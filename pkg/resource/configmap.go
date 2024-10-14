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

type configMapManagerInterface interface {
	GetConfigMap(string) (*apiv1.ConfigMap, error)
	UpdateConfigMap(string, func(*apiv1.ConfigMap) error, metav1.UpdateOptions) error
	ListConfigMaps(labels.Selector) ([]*apiv1.ConfigMap, error)
	DeleteConfigMap(string) error
	CreateConfigMap(*apiv1.ConfigMap, metav1.CreateOptions) error
}

type ConfigMapManager struct {
	configmaps corev1listers.ConfigMapLister
}

func newConfigMapManager(configmaps corev1listers.ConfigMapLister) *ConfigMapManager {
	return &ConfigMapManager{
		configmaps: configmaps,
	}
}

func (cmm *ConfigMapManager) GetConfigMap(name string) (*apiv1.ConfigMap, error) {
	cm, err := cmm.configmaps.ConfigMaps(config.GetNamespace()).Get(name)

	if err != nil {
		if !apierr.IsNotFound(err) {
			return nil, err
		}
		log.Infof("configmap %s not found, will create", name)
		cm = &apiv1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:   name,
				Labels: ManageByGatepointLabelMap(),
			},
			Data: make(map[string]string),
		}
		err = cmm.CreateConfigMap(cm, metav1.CreateOptions{})
		if err != nil {
			return nil, err
		}
		log.Infof("configmap %s created success", name)
	}
	if cm.Data == nil {
		cm.Data = map[string]string{}
	}
	return cm, nil
}

func (cmm *ConfigMapManager) UpdateConfigMap(name string, callback func(*apiv1.ConfigMap) error, updateOptions metav1.UpdateOptions) error {
	cm, err := cmm.GetConfigMap(name)
	if err != nil {
		return err
	}

	beforeUpdate := cm.DeepCopy()
	err = callback(cm)
	if err != nil {
		return err
	}
	if reflect.DeepEqual(beforeUpdate.Data, cm.Data) {
		return nil
	}

	_, err = Manager.clientset.CoreV1().ConfigMaps(config.GetNamespace()).Update(context.Background(), cm, updateOptions)
	return err
}

func (cmm *ConfigMapManager) ListConfigMaps(selector labels.Selector) ([]*apiv1.ConfigMap, error) {
	return cmm.configmaps.ConfigMaps(config.GetNamespace()).List(selector)
}

func (cmm *ConfigMapManager) DeleteConfigMap(name string) error {
	return Manager.clientset.CoreV1().ConfigMaps(config.GetNamespace()).Delete(context.Background(), name, metav1.DeleteOptions{})
}

func (cmm *ConfigMapManager) CreateConfigMap(cm *apiv1.ConfigMap, createOptions metav1.CreateOptions) error {
	_, err := Manager.clientset.CoreV1().ConfigMaps(config.GetNamespace()).Create(context.Background(), cm, createOptions)
	return err
}
