package kube

import (
	"log"

	"github.com/gatepoint/gatepoint/pkg/config"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	k8sclicfg "sigs.k8s.io/controller-runtime/pkg/client/config"
)

var kubeClient kubernetes.Interface

func GetKubeClient() kubernetes.Interface {
	if kubeClient != nil {
		return kubeClient
	}
	cli, err := InitClientset(config.GetKubeConfig())
	if err != nil {
		panic(err)
	}
	kubeClient = cli
	return cli
}

func InitClientset(kubeconfig string, opts ...func(config *rest.Config)) (*kubernetes.Clientset, error) {
	var config *rest.Config

	if kubeconfig == "" {
		config = getInClusterConfig()
	} else {
		config = getClusterConfig(kubeconfig)
	}

	for _, opt := range opts {
		opt(config)
	}

	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("init kubernetes clientset error: %s", err)
	}
	kubeClient = clientset
	return clientset, nil
}

func getClusterConfig(kubeconfig string) *rest.Config {
	log.Println("init kubernetes config with out-cluster")
	restConfig, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatalf("init kubernetes config error: %s", err)
	}
	return restConfig
}

func getInClusterConfig() *rest.Config {
	log.Println("init kubernetes config with in-cluster")
	return k8sclicfg.GetConfigOrDie()
}
