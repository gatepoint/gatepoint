package account

import v1 "k8s.io/api/core/v1"

const (
	accountsKeyPrefix          = "accounts"
	accountPasswordSuffix      = "password"
	accountPasswordMtimeSuffix = "passwordMtime"
	accountEnabledSuffix       = "enabled"
	accountTokensSuffix        = "tokens"

	// Admin superuser password storage
	// settingAdminPasswordHashKey designates the key for a root password hash inside a Kubernetes secret.
	settingAdminPasswordHashKey = "admin.password"
	// settingAdminPasswordMtimeKey designates the key for a root password mtime inside a Kubernetes secret.
	settingAdminPasswordMtimeKey = "admin.passwordMtime"
	settingAdminEnabledKey       = "admin.enabled"
	settingAdminTokensKey        = "admin.tokens"
)

const (
	passwordLetters            = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"
	initialPasswordSecretName  = "gatepoint-initial-admin-secret"
	initialPasswordSecretField = "password"

	initialPasswordLength = 16
)

func updateAccountSecret(secret *v1.Secret, key string, val string, defVal string) {
	existingVal := string(secret.Data[key])
	if existingVal != val {
		if val == "" || val == defVal {
			delete(secret.Data, key)
		} else {
			secret.Data[key] = []byte(val)
		}
	}
}

func updateAccountMap(cm *v1.ConfigMap, key string, val string, defVal string) {
	existingVal := cm.Data[key]
	if existingVal != val {
		if val == "" || val == defVal {
			delete(cm.Data, key)
		} else {
			cm.Data[key] = val
		}
	}
}
