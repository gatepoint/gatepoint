package resource

import (
	"context"
	"reflect"

	"github.com/gatepoint/gatepoint/pkg/config"
	"github.com/gatepoint/gatepoint/pkg/log"
	apiv1 "k8s.io/api/core/v1"
	apierr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	corev1listers "k8s.io/client-go/listers/core/v1"
)

type secretManagerInterface interface {
	GetSecret(string) (*apiv1.Secret, error)
	UpdateSecret(string, func(*apiv1.Secret) error, metav1.UpdateOptions) error
	ListSecrets(labels.Selector) ([]*apiv1.Secret, error)
	DeleteSecret(string) error
	CreateSecret(*apiv1.Secret, metav1.CreateOptions) error
}

type SecretManager struct {
	secrets corev1listers.SecretLister
}

func newSecretManager(secrets corev1listers.SecretLister) *SecretManager {
	return &SecretManager{
		secrets: secrets,
	}
}

func (sm *SecretManager) GetSecret(name string) (*apiv1.Secret, error) {
	secret, err := sm.secrets.Secrets(config.GetNamespace()).Get(name)
	if err != nil {
		if !apierr.IsNotFound(err) {
			return nil, err
		}
		log.Infof("secret %s not found,will create", name)
		secret = &apiv1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:   name,
				Labels: ManageByGatepointLabelMap(),
			},
			Data: make(map[string][]byte),
			Type: apiv1.SecretTypeOpaque,
		}

		err = sm.CreateSecret(secret, metav1.CreateOptions{})
		if err != nil {
			return nil, err
		}
		log.Infof("secret %s created success", name)
	}
	if secret.Data == nil {
		secret.Data = make(map[string][]byte)
	}
	return secret, nil
}

func (sm *SecretManager) UpdateSecret(name string, callback func(*apiv1.Secret) error, updateOptions metav1.UpdateOptions) error {
	secret, err := sm.GetSecret(name)
	if err != nil {
		return err
	}

	if secret.Data == nil {
		secret.Data = make(map[string][]byte)
	}

	beforeUpdate := secret.DeepCopy()
	err = callback(secret)
	if err != nil {
		return err
	}
	if reflect.DeepEqual(beforeUpdate.Data, secret.Data) {
		return nil
	}

	_, err = Manager.clientset.CoreV1().Secrets(config.GetNamespace()).Update(context.Background(), secret, updateOptions)
	return err
}

func (sm *SecretManager) ListSecrets(selector labels.Selector) ([]*apiv1.Secret, error) {
	return sm.secrets.Secrets(config.GetNamespace()).List(selector)
}

func (sm *SecretManager) DeleteSecret(name string) error {
	return Manager.clientset.CoreV1().Secrets(config.GetNamespace()).Delete(context.Background(), name, metav1.DeleteOptions{})
}

func (sm *SecretManager) CreateSecret(secret *apiv1.Secret, createOptions metav1.CreateOptions) error {
	_, err := Manager.clientset.CoreV1().Secrets(config.GetNamespace()).Create(context.Background(), secret, createOptions)
	return err
}
