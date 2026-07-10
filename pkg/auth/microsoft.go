package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// Microsoft OAuth2 client ID shared by third-party Minecraft launchers (Xbox app client)
const msClientID = "00000000402b5328"

// GetMicrosoftAuthURL returns the URL to open in the browser to begin Microsoft login.
func GetMicrosoftAuthURL(redirectPort int) string {
	redirectURI := fmt.Sprintf("http://localhost:%d/callback", redirectPort)
	params := url.Values{}
	params.Set("client_id", msClientID)
	params.Set("response_type", "code")
	params.Set("redirect_uri", redirectURI)
	params.Set("scope", "XboxLive.signin offline_access")
	params.Set("prompt", "select_account")
	return "https://login.live.com/oauth20_authorize.srf?" + params.Encode()
}

// LoginWithMicrosoft runs the full 6-step token chain and returns a fully-populated Account.
func LoginWithMicrosoft(authCode string, redirectPort int) (Account, error) {
	redirectURI := fmt.Sprintf("http://localhost:%d/callback", redirectPort)

	// Step 2: Exchange auth code for Microsoft access token
	msToken, msRefresh, err := getMicrosoftToken(authCode, redirectURI)
	if err != nil {
		return Account{}, fmt.Errorf("MS token exchange failed: %w", err)
	}

	// Step 3: Authenticate with Xbox Live
	xblToken, userHash, err := getXboxLiveToken(msToken)
	if err != nil {
		return Account{}, fmt.Errorf("XBL auth failed: %w", err)
	}

	// Step 4: Get XSTS token
	xstsToken, err := getXSTSToken(xblToken)
	if err != nil {
		return Account{}, fmt.Errorf("XSTS auth failed: %w", err)
	}

	// Step 5: Get Minecraft access token
	mcToken, err := getMinecraftToken(xstsToken, userHash)
	if err != nil {
		return Account{}, fmt.Errorf("Minecraft token failed: %w", err)
	}

	// Step 6: Fetch Minecraft profile
	uuid, username, err := getMinecraftProfile(mcToken)
	if err != nil {
		return Account{}, fmt.Errorf("profile fetch failed: %w", err)
	}

	return Account{
		ID:           uuid,
		Type:         TypeMicrosoft,
		Username:     username,
		AccessToken:  mcToken,
		RefreshToken: msRefresh,
	}, nil
}

// RefreshMicrosoftToken silently refreshes the Microsoft access token using the stored refresh token.
// Returns the updated Account if successful.
func RefreshMicrosoftToken(acc *Account) (*Account, error) {
	if acc.RefreshToken == "" {
		return nil, fmt.Errorf("no refresh token stored for account %s", acc.Username)
	}

	body := url.Values{}
	body.Set("client_id", msClientID)
	body.Set("refresh_token", acc.RefreshToken)
	body.Set("grant_type", "refresh_token")
	body.Set("scope", "XboxLive.signin offline_access")

	resp, err := http.PostForm("https://login.live.com/oauth20_token.srf", body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		Error        string `json:"error"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	if result.Error != "" {
		return nil, fmt.Errorf("refresh failed: %s", result.Error)
	}

	// Re-run the XBL → XSTS → MC chain with the new MS token
	xblToken, userHash, err := getXboxLiveToken(result.AccessToken)
	if err != nil {
		return nil, err
	}
	xstsToken, err := getXSTSToken(xblToken)
	if err != nil {
		return nil, err
	}
	mcToken, err := getMinecraftToken(xstsToken, userHash)
	if err != nil {
		return nil, err
	}

	acc.AccessToken = mcToken
	acc.RefreshToken = result.RefreshToken
	return acc, nil
}

// getMicrosoftToken exchanges an auth code for a Microsoft access + refresh token pair.
func getMicrosoftToken(code, redirectURI string) (accessToken, refreshToken string, err error) {
	body := url.Values{}
	body.Set("client_id", msClientID)
	body.Set("code", code)
	body.Set("grant_type", "authorization_code")
	body.Set("redirect_uri", redirectURI)

	resp, err := http.PostForm("https://login.live.com/oauth20_token.srf", body)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	var result struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		Error        string `json:"error"`
		ErrorDesc    string `json:"error_description"`
	}
	rawBody, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(rawBody, &result); err != nil {
		return "", "", fmt.Errorf("failed to parse MS token response: %w", err)
	}
	if result.Error != "" {
		return "", "", fmt.Errorf("%s: %s", result.Error, result.ErrorDesc)
	}
	return result.AccessToken, result.RefreshToken, nil
}

// getXboxLiveToken exchanges a Microsoft access token for an Xbox Live token.
func getXboxLiveToken(msAccessToken string) (xblToken, userHash string, err error) {
	payload := map[string]interface{}{
		"Properties": map[string]interface{}{
			"AuthMethod": "RPS",
			"SiteName":   "user.auth.xboxlive.com",
			"RpsTicket":  msAccessToken,
		},
		"RelyingParty": "http://auth.xboxlive.com",
		"TokenType":    "JWT",
	}

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "https://user.auth.xboxlive.com/user/authenticate", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	var result struct {
		Token         string `json:"Token"`
		DisplayClaims struct {
			Xui []struct {
				Uhs string `json:"uhs"`
			} `json:"xui"`
		} `json:"DisplayClaims"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", "", err
	}
	if result.Token == "" {
		return "", "", fmt.Errorf("empty XBL token in response")
	}
	if len(result.DisplayClaims.Xui) == 0 {
		return "", "", fmt.Errorf("no user hash in XBL response")
	}
	return result.Token, result.DisplayClaims.Xui[0].Uhs, nil
}

// getXSTSToken exchanges an XBL token for an XSTS token required by Mojang.
func getXSTSToken(xblToken string) (string, error) {
	payload := map[string]interface{}{
		"Properties": map[string]interface{}{
			"SandboxId":  "RETAIL",
			"UserTokens": []string{xblToken},
		},
		"RelyingParty": "rp://api.minecraftservices.com/",
		"TokenType":    "JWT",
	}

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "https://xsts.auth.xboxlive.com/xsts/authorize", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// XSTS can return an error body even on non-200 status
	rawBody, _ := io.ReadAll(resp.Body)

	var result struct {
		Token string `json:"Token"`
		XErr  int64  `json:"XErr"`
	}
	if err := json.Unmarshal(rawBody, &result); err != nil {
		return "", err
	}

	// Translate known XErr codes
	switch result.XErr {
	case 2148916233:
		return "", fmt.Errorf("this Microsoft account has no Xbox account. Please create one at xbox.com")
	case 2148916235:
		return "", fmt.Errorf("Xbox Live is not available in your region")
	case 2148916238:
		return "", fmt.Errorf("this account is a child account — an adult must add it to a Family")
	}
	if result.Token == "" {
		return "", fmt.Errorf("empty XSTS token (XErr: %d)", result.XErr)
	}
	return result.Token, nil
}

// getMinecraftToken exchanges XSTS credentials for a Minecraft access token.
func getMinecraftToken(xstsToken, userHash string) (string, error) {
	payload := map[string]string{
		"identityToken": fmt.Sprintf("XBL3.0 x=%s;%s", userHash, xstsToken),
	}

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "https://api.minecraftservices.com/authentication/login_with_xbox", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	if result.AccessToken == "" {
		return "", fmt.Errorf("empty Minecraft access token in response")
	}
	return result.AccessToken, nil
}

// getMinecraftProfile fetches the player's UUID and username using their Minecraft access token.
func getMinecraftProfile(mcToken string) (uuid, username string, err error) {
	req, _ := http.NewRequest("GET", "https://api.minecraftservices.com/minecraft/profile", nil)
	req.Header.Set("Authorization", "Bearer "+mcToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	var profile struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
		return "", "", err
	}

	// Format UUID with dashes: 8-4-4-4-12
	raw := profile.ID
	if len(raw) == 32 {
		uuid = strings.Join([]string{raw[0:8], raw[8:12], raw[12:16], raw[16:20], raw[20:32]}, "-")
	} else {
		uuid = raw
	}

	if profile.Name == "" {
		return "", "", fmt.Errorf("Minecraft profile returned empty name — account may not own Minecraft")
	}

	return uuid, profile.Name, nil
}
