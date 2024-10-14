package resource

import (
	"context"
	"errors"
	"reflect"
	"sync"
	"time"

	"github.com/gatepoint/gatepoint/pkg/common"
	"github.com/gatepoint/gatepoint/pkg/config"
	"github.com/gatepoint/gatepoint/pkg/log"
	apiv1 "k8s.io/api/core/v1"
	apierr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	v1 "k8s.io/client-go/informers/core/v1"
	corev1listers "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"

	"k8s.io/client-go/kubernetes"
)

const (
	settingsResyncDuration = 3 * time.Minute
)

var Manager *ResourceManager

type ResourceManager struct {
	ctx       context.Context
	clientset kubernetes.Interface

	ConfigMapManager
	SecretManager

	//secrets          corev1listers.SecretLister
	//secretsInformers cache.SharedIndexInformer
	//
	//configmaps corev1listers.ConfigMapLister

	initContextCancel func()

	//subscribers []chan<- struct{}

	mu *sync.Mutex
}

type ReourceManagerOpts func(opts *ResourceManager)

func NewResourceManager(ctx context.Context, clientset kubernetes.Interface, opts ...ReourceManagerOpts) error {
	Manager = &ResourceManager{
		ctx:       ctx,
		clientset: clientset,
		mu:        &sync.Mutex{},
	}

	for i := range opts {
		opts[i](Manager)
	}
	return Manager.ensureSynced(true)

	//return mgr
}
func (mgr *ResourceManager) ensureSynced(forceResync bool) error {
	mgr.mu.Lock()
	defer mgr.mu.Unlock()
	if !forceResync && mgr.secrets != nil && mgr.configmaps != nil {
		return nil
	}

	if mgr.initContextCancel != nil {
		mgr.initContextCancel()
	}
	ctx, cancel := context.WithCancel(mgr.ctx)
	mgr.initContextCancel = cancel
	return mgr.initialize(ctx)
}

func (mgr *ResourceManager) initialize(ctx context.Context) error {

	eventHandler := cache.ResourceEventHandlerFuncs{
		UpdateFunc: func(oldObj, newObj interface{}) {
			log.Infof("Resource updated from %v to %v", oldObj, newObj)
		},
		AddFunc: func(obj interface{}) {
			log.Infof("Resource added %v", obj)

		},
		DeleteFunc: func(obj interface{}) {
			log.Infof("Resource deleted %v", obj)
		},
	}
	cmInformer := v1.NewFilteredConfigMapInformer(mgr.clientset, config.GetNamespace(), settingsResyncDuration, cache.Indexers{}, func(options *metav1.ListOptions) {
		options.LabelSelector = fields.ParseSelectorOrDie(ManageByGatepointLabel()).String()
	})

	secretsInformer := v1.NewSecretInformer(mgr.clientset, config.GetNamespace(), 3*time.Minute, cache.Indexers{})

	if _, err := secretsInformer.AddEventHandler(eventHandler); err != nil {
		log.Errorf("Failed to add event handler to Secret informer: %v", err)
	}

	if _, err := cmInformer.AddEventHandler(eventHandler); err != nil {
		log.Errorf("Failed to add event handler to ConfigMap informer: %v", err)
	}

	log.Info("Starting ConfigMap informer")
	go func() {
		cmInformer.Run(ctx.Done())
		log.Info("ConfigMap informer stopped")
	}()

	log.Info("Starting Secret informer")
	go func() {
		secretsInformer.Run(ctx.Done())
		log.Info("Secret informer stopped")
	}()

	if !cache.WaitForCacheSync(ctx.Done(), cmInformer.HasSynced, secretsInformer.HasSynced) {
		log.Errorf("Timed out waiting for resource cache to sync")
		return errors.New("failed to sync caches")
	}

	log.Info("Resource cache synced")

	//todo add more eventHandlers

	mgr.secrets = corev1listers.NewSecretLister(secretsInformer.GetIndexer())
	mgr.configmaps = corev1listers.NewConfigMapLister(cmInformer.GetIndexer())
	//mgr.secretsInformers = secretsInformer
	return nil

}

func (mgr *ResourceManager) invalidateCache() {
	mgr.mu.Lock()
	defer mgr.mu.Unlock()
	//mgr.secrets = nil
	//mgr.configmaps = nil
}

func (mgr *ResourceManager) ResyncInformers() error {
	return mgr.ensureSynced(true)
}

func (mgr *ResourceManager) updateSecret(callback func(*apiv1.Secret) error) error {
	err := mgr.ensureSynced(false)
	if err != nil {
		return err
	}
	argoCDSecret, err := mgr.secrets.Secrets(config.GetNamespace()).Get(common.GatepointSecretName)
	createSecret := false
	if err != nil {
		if !apierr.IsNotFound(err) {
			return err
		}
		argoCDSecret = &apiv1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name: common.GatepointSecretName,
			},
			Data: make(map[string][]byte),
		}
		createSecret = true
	}
	if argoCDSecret.Data == nil {
		argoCDSecret.Data = make(map[string][]byte)
	}

	updatedSecret := argoCDSecret.DeepCopy()
	err = callback(updatedSecret)
	if err != nil {
		return err
	}

	if !createSecret && reflect.DeepEqual(argoCDSecret.Data, updatedSecret.Data) {
		return nil
	}

	if createSecret {
		_, err = mgr.clientset.CoreV1().Secrets(config.GetNamespace()).Create(context.Background(), updatedSecret, metav1.CreateOptions{})
	} else {
		_, err = mgr.clientset.CoreV1().Secrets(config.GetNamespace()).Update(context.Background(), updatedSecret, metav1.UpdateOptions{})
	}
	if err != nil {
		return err
	}

	return mgr.ResyncInformers()
}
