package auth

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"Aether/pkg/fs"
)

// AccountType defines the type of account (e.g., offline, microsoft)
type AccountType string

const (
	TypeOffline   AccountType = "offline"
	TypeMicrosoft AccountType = "microsoft"
)

// Account represents a user profile
type Account struct {
	ID       string      `json:"id"`
	Type     AccountType `json:"type"`
	Username string      `json:"username"`
}

// AccountStore represents the accounts.json file structure
type AccountStore struct {
	ActiveAccountID string    `json:"activeAccountId"`
	Accounts        []Account `json:"accounts"`
}

var store AccountStore

func getStorePath() string {
	return filepath.Join(fs.GetDataDir(), "accounts.json")
}

// LoadAccounts reads accounts.json into memory
func LoadAccounts() error {
	path := getStorePath()
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			// Initialize empty store
			store = AccountStore{Accounts: []Account{}}
			return nil
		}
		return err
	}

	return json.Unmarshal(data, &store)
}

// SaveAccounts writes the in-memory store to accounts.json
func SaveAccounts() error {
	path := getStorePath()
	
	// Ensure directory exists (though fs.GetDataDir() should handle base dir)
	os.MkdirAll(filepath.Dir(path), 0755)

	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// GenerateOfflineUUID generates a deterministic UUID for offline/cracked mode based on username
func GenerateOfflineUUID(username string) string {
	hash := md5.Sum([]byte("OfflinePlayer:" + username))
	// Set version 3 and variant bits
	hash[6] = (hash[6] & 0x0f) | 0x30
	hash[8] = (hash[8] & 0x3f) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", hash[0:4], hash[4:6], hash[6:8], hash[8:10], hash[10:16])
}

// AddOfflineAccount creates a new offline account and sets it as active
func AddOfflineAccount(username string) (Account, error) {
	if err := LoadAccounts(); err != nil {
		return Account{}, err
	}

	id := GenerateOfflineUUID(username)
	
	// Check if account already exists
	for _, acc := range store.Accounts {
		if acc.ID == id {
			store.ActiveAccountID = id
			SaveAccounts()
			return acc, nil
		}
	}

	newAcc := Account{
		ID:       id,
		Type:     TypeOffline,
		Username: username,
	}

	store.Accounts = append(store.Accounts, newAcc)
	store.ActiveAccountID = id
	err := SaveAccounts()
	return newAcc, err
}

// GetActiveAccount returns the currently active account, or nil if none
func GetActiveAccount() *Account {
	if err := LoadAccounts(); err != nil {
		return nil
	}

	for _, acc := range store.Accounts {
		if acc.ID == store.ActiveAccountID {
			return &acc
		}
	}
	return nil
}

// GetAccounts returns all saved accounts
func GetAccounts() []Account {
	if err := LoadAccounts(); err != nil {
		return []Account{}
	}
	return store.Accounts
}

// SetActiveAccount sets the active account by ID
func SetActiveAccount(id string) error {
	if err := LoadAccounts(); err != nil {
		return err
	}
	
	for _, acc := range store.Accounts {
		if acc.ID == id {
			store.ActiveAccountID = id
			return SaveAccounts()
		}
	}
	return fmt.Errorf("account not found")
}
