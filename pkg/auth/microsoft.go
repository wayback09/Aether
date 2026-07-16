package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Microsoft OAuth2 client ID shared by third-party Minecraft launchers (Xbox app client)
const msClientID = "00000000402b5328"

type DeviceCodeResponse struct {
	UserCode        string `json:"user_code"`
	DeviceCode      string `json:"device_code"`
	VerificationURI string `json:"verification_uri"`
	ExpiresIn       int    `json:"expires_in"`
	Interval        int    `json:"interval"`
	Message         string `json:"message"`
}

// GetDeviceCode initiates the Microsoft device code flow.
func GetDeviceCode() (*DeviceCodeResponse, error) {
	body := url.Values{}
	body.Set("client_id", msClientID)
	body.Set("scope", "XboxLive.signin offline_access")

	resp, err := http.PostForm("https://login.microsoftonline.com/consumers/oauth2/v2.0/devicecode", body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		raw, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to request device code: %s", string(raw))
	}

	var result DeviceCodeResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

// PollForToken repeatedly checks if the user has completed the device code login.
func PollForToken(deviceCode string, intervalSeconds int) (accessToken, refreshToken string, err error) {
	ticker := time.NewTicker(time.Duration(intervalSeconds) * time.Second)
	defer ticker.Stop()
	
	// Timeout after 15 minutes (max for device code flow)
	timeout := time.After(15 * time.Minute)

	for {
		select {
		case <-timeout:
			return "", "", fmt.Errorf("authentication timed out")
		case <-ticker.C:
			body := url.Values{}
			body.Set("client_id", msClientID)
			body.Set("grant_type", "urn:ietf:params:oauth:grant-type:device_code")
			body.Set("device_code", deviceCode)

			resp, err := http.PostForm("https://login.microsoftonline.com/consumers/oauth2/v2.0/token", body)
			if err != nil {
				return "", "", err
			}

			var result struct {
				Error        string `json:"error"`
				ErrorDesc    string `json:"error_description"`
				AccessToken  string `json:"access_token"`
				RefreshToken string `json:"refresh_token"`
			}

			rawBody, _ := io.ReadAll(resp.Body)
			resp.Body.Close()

			if err := json.Unmarshal(rawBody, &result); err != nil {
				return "", "", fmt.Errorf("failed to parse MS token response: %w", err)
			}

			if result.Error == "authorization_pending" {
				continue // Still waiting
			}

			if result.Error != "" {
				return "", "", fmt.Errorf("%s: %s", result.Error, result.ErrorDesc)
			}

			return result.AccessToken, result.RefreshToken, nil
		}
	}
}

// LoginWithMicrosoft runs the full 5-step token chain using a fresh Microsoft access token.
func LoginWithMicrosoft(msToken, msRefresh string) (Account, error) {
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

	// Re-run the XBL -> XSTS -> MC chain with the new MS token
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

// getXboxLiveToken exchanges a Microsoft access token for an Xbox Live token.
func getXboxLiveToken(msAccessToken string) (xblToken, userHash string, err error) {
	payload := map[string]interface{}{
		"Properties": map[string]interface{}{
			"AuthMethod": "RPS",
			"SiteName":   "user.auth.xboxlive.com",
			"RpsTicket":  "d=" + msAccessToken, // Important prefix for modern MS tokens in Minecraft API
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
