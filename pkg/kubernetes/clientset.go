package kubernetes

import (
	"flag"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func GetClientset(kubeconfig string) (kubernetes.Clientset, error) {
	var config = rest.Config{}

	if kubeconfig == "" {
		config = GetInClusterConfig()
	} else {
		config = GetOutClusterConfig(kubeconfig)
	}

	// creates the clientset
	clientset, err := kubernetes.NewForConfig(&config)
	if err != nil {
		panic(err.Error())
	}
	return *clientset, nil
}

func GetInClusterConfig() rest.Config {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	return *config
}

func GetOutClusterConfig(kubeconfig string) rest.Config {
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = *flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = *flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()
	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	return *config
}
