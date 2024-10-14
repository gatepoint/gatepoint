package account

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/gatepoint/gatepoint/pkg/common"
	"github.com/gatepoint/gatepoint/pkg/log"
	"github.com/gatepoint/gatepoint/pkg/resource"
	"github.com/gatepoint/gatepoint/pkg/utils"
	"github.com/gatepoint/gatepoint/pkg/utils/password"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/retry"
)

type AccountInterface interface {
	UpdateAccount(name string, callback func(account *Account) error) error
	CreateAccount(account Account) error
}

type AccountManager struct {
	resourceManager *resource.ResourceManager
}

func NewAccountManager() *AccountManager {
	accountManager := &AccountManager{
		resourceManager: resource.Manager,
	}
	accountManager.Init()
	return accountManager
}

type Token struct {
	ID        string
	IssuedAt  int64
	ExpiresAt int64
}

type AccountCapability string

const (
	// AccountCapabilityLogin represents capability to create UI session tokens.
	AccountCapabilityLogin AccountCapability = "login"
	// AccountCapabilityApiKey represents capability to generate API auth tokens.
	AccountCapabilityApiKey AccountCapability = "apiKey"
)

type Account struct {
	Username      string
	PasswordHash  string
	PasswordMtime *time.Time
	Enabled       bool
	Capabilities  []AccountCapability
	Tokens        []Token
}

func (am *AccountManager) UpdateAccount(name string, callback func(account *Account) error) error {

	return retry.RetryOnConflict(retry.DefaultBackoff, func() error {
		account, err := am.GetAccount(name)
		if err != nil {
			return err
		}
		err = callback(account)
		if err != nil {
			return err
		}
		return am.saveAccount(account)
	})
}

func (am *AccountManager) CreateAccount(account *Account) error {
	exist, err := am.accountExist(account.Username)
	if err != nil {
		return err
	}

	if exist {
		return fmt.Errorf("account %s already exists", account.Username)
	}
	return am.saveAccount(account)
}

func (am *AccountManager) GetAccount(name string) (*Account, error) {
	accounts, err := am.ListAccounts()
	if err != nil {
		return nil, err
	}

	if account, ok := accounts[name]; ok {
		return &account, nil
	}
	//return nil, errors.ErrMap[commonv1.ErrType_ERR_TYPE_NOT_FOUND].WithError(fmt.Errorf("account %s not found", name))
	return nil, fmt.Errorf("account %s not found", name)
}

func (am *AccountManager) accountExist(name string) (bool, error) {
	accounts, err := am.ListAccounts()
	if err != nil {
		return false, err
	}

	if _, ok := accounts[name]; ok {
		return true, nil
	}
	return false, nil
}

func (am *AccountManager) saveAccount(account *Account) error {
	userName := account.Username
	return am.resourceManager.UpdateSecret(userName, func(secret *v1.Secret) error {
		return am.resourceManager.UpdateConfigMap(userName, func(cm *v1.ConfigMap) error {
			return saveAccount(secret, cm, account)
		}, metav1.UpdateOptions{})
	}, metav1.UpdateOptions{})
}

func saveAccount(secret *v1.Secret, cm *v1.ConfigMap, account *Account) error {
	name := account.Username
	tokens, err := json.Marshal(account.Tokens)
	if err != nil {
		return err
	}
	if name == common.GatepointAdminUsername {
		updateAccountSecret(secret, settingAdminPasswordHashKey, account.PasswordHash, "")
		updateAccountSecret(secret, settingAdminPasswordMtimeKey, account.FormatPasswordMtime(), "")
		updateAccountSecret(secret, settingAdminTokensKey, string(tokens), "[]")
		updateAccountMap(cm, settingAdminEnabledKey, strconv.FormatBool(account.Enabled), "true")
	} else {
		updateAccountSecret(secret, fmt.Sprintf("%s.%s.%s", accountsKeyPrefix, name, accountPasswordSuffix), account.PasswordHash, "")
		updateAccountSecret(secret, fmt.Sprintf("%s.%s.%s", accountsKeyPrefix, name, accountPasswordMtimeSuffix), account.FormatPasswordMtime(), "")
		updateAccountSecret(secret, fmt.Sprintf("%s.%s.%s", accountsKeyPrefix, name, accountTokensSuffix), string(tokens), "[]")
		updateAccountMap(cm, fmt.Sprintf("%s.%s.%s", accountsKeyPrefix, name, accountEnabledSuffix), strconv.FormatBool(account.Enabled), "true")
		updateAccountMap(cm, fmt.Sprintf("%s.%s", accountsKeyPrefix, name), account.FormatCapabilities(), "")
	}
	return nil
}

// FormatCapabilities returns comma separate list of user capabilities.
func (a *Account) FormatCapabilities() string {
	var items []string
	for i := range a.Capabilities {
		items = append(items, string(a.Capabilities[i]))
	}
	return strings.Join(items, ",")
}

// FormatPasswordMtime return the formatted password modify time or empty string of password modify time is nil.
func (a *Account) FormatPasswordMtime() string {
	if a.PasswordMtime == nil {
		return ""
	}
	return a.PasswordMtime.Format(time.RFC3339)
}

func (am *AccountManager) DeleteAccount(name string) error {
	return nil
}

func (am *AccountManager) ListAccounts() (map[string]Account, error) {

	secret, err := am.resourceManager.GetSecret(common.GatepointSecretName)

	if err != nil {
		return nil, err
	}

	cm, err := am.resourceManager.GetConfigMap(common.GatepointConfigMapName)

	if err != nil {
		return nil, err
	}

	return parseAccounts(secret, cm)
}

func (am *AccountManager) Init() {
	am.UpdateAccount(common.GatepointAdminUsername, func(adminAccount *Account) error {
		if adminAccount.Enabled {
			now := utils.NowUTC()
			if adminAccount.PasswordHash == "" {
				randBytes := make([]byte, initialPasswordLength)
				for i := 0; i < initialPasswordLength; i++ {
					num, err := rand.Int(rand.Reader, big.NewInt(int64(len(passwordLetters))))
					if err != nil {
						return err
					}
					randBytes[i] = passwordLetters[num.Int64()]
				}
				initialPassword := string(randBytes)

				hashedPassword, err := password.HashPassword(initialPassword)
				if err != nil {
					return err
				}

				err = am.resourceManager.UpdateSecret(initialPasswordSecretName, func(secret *v1.Secret) error {
					secret.Data[initialPasswordSecretField] = []byte(initialPassword)
					return nil
				}, metav1.UpdateOptions{})

				if err != nil {
					return err
				}
				adminAccount.PasswordHash = hashedPassword
				adminAccount.PasswordMtime = now
				log.Info("Initialized admin password")
			}
			if adminAccount.PasswordMtime == nil || adminAccount.PasswordMtime.IsZero() {
				adminAccount.PasswordMtime = now
				log.Info("Initialized admin mtime")
			}
		} else {
			log.Info("admin disabled")
		}
		return nil
	})
}

func parseAdminAccount(secret *v1.Secret, cm *v1.ConfigMap) (*Account, error) {
	adminAccount := &Account{Enabled: true, Capabilities: []AccountCapability{AccountCapabilityLogin}, Username: common.GatepointAdminUsername}
	if adminPasswordHash, ok := secret.Data[settingAdminPasswordHashKey]; ok {
		adminAccount.PasswordHash = string(adminPasswordHash)
	}
	if adminPasswordMtimeBytes, ok := secret.Data[settingAdminPasswordMtimeKey]; ok {
		if mTime, err := time.Parse(time.RFC3339, string(adminPasswordMtimeBytes)); err == nil {
			adminAccount.PasswordMtime = &mTime
		}
	}

	adminAccount.Tokens = make([]Token, 0)
	if tokensStr, ok := secret.Data[settingAdminTokensKey]; ok && string(tokensStr) != "" {
		if err := json.Unmarshal(tokensStr, &adminAccount.Tokens); err != nil {
			return nil, err
		}
	}

	if enabledStr, ok := cm.Data[settingAdminEnabledKey]; ok {
		if enabled, err := strconv.ParseBool(enabledStr); err == nil {
			adminAccount.Enabled = enabled
		} else {
			log.Warnf("ConfigMap has invalid key %s: %v", settingAdminTokensKey, err)
		}
	}

	return adminAccount, nil
}

func parseAccounts(secret *v1.Secret, cm *v1.ConfigMap) (map[string]Account, error) {
	adminAccount, err := parseAdminAccount(secret, cm)
	if err != nil {
		return nil, err
	}
	accounts := map[string]Account{
		common.GatepointAdminUsername: *adminAccount,
	}

	for key, v := range cm.Data {
		if !strings.HasPrefix(key, fmt.Sprintf("%s.", accountsKeyPrefix)) {
			continue
		}

		val := v
		var accountName, suffix string

		parts := strings.Split(key, ".")
		switch len(parts) {
		case 2:
			accountName = parts[1]
		case 3:
			accountName = parts[1]
			suffix = parts[2]
		default:
			log.Warnf("Unexpected key %s in ConfigMap '%s'", key, cm.Name)
			continue
		}

		account, ok := accounts[accountName]
		if !ok {
			account = Account{Enabled: true}
			accounts[accountName] = account
		}
		switch suffix {
		case "":
			for _, capability := range strings.Split(val, ",") {
				capability = strings.TrimSpace(capability)
				if capability == "" {
					continue
				}

				switch capability {
				case string(AccountCapabilityLogin):
					account.Capabilities = append(account.Capabilities, AccountCapabilityLogin)
				case string(AccountCapabilityApiKey):
					account.Capabilities = append(account.Capabilities, AccountCapabilityApiKey)
				default:
					log.Warnf("not supported account capability '%s' in config map key '%s'", capability, key)
				}
			}
		case accountEnabledSuffix:
			account.Enabled, err = strconv.ParseBool(val)
			if err != nil {
				return nil, err
			}
		}
		accounts[accountName] = account
	}

	for name, account := range accounts {
		if name == common.GatepointAdminUsername {
			continue
		}

		if passwordHash, ok := secret.Data[fmt.Sprintf("%s.%s.%s", accountsKeyPrefix, name, accountPasswordSuffix)]; ok {
			account.PasswordHash = string(passwordHash)
		}
		if passwordMtime, ok := secret.Data[fmt.Sprintf("%s.%s.%s", accountsKeyPrefix, name, accountPasswordMtimeSuffix)]; ok {
			if mTime, err := time.Parse(time.RFC3339, string(passwordMtime)); err != nil {
				return nil, err
			} else {
				account.PasswordMtime = &mTime
			}
		}
		if tokensStr, ok := secret.Data[fmt.Sprintf("%s.%s.%s", accountsKeyPrefix, name, accountTokensSuffix)]; ok {
			account.Tokens = make([]Token, 0)
			if string(tokensStr) != "" {
				if err := json.Unmarshal(tokensStr, &account.Tokens); err != nil {
					log.Errorf("Account '%s' has invalid token in secret '%s'", name, secret.Name)
				}
			}
		}
		accounts[name] = account
	}

	return accounts, nil
}
