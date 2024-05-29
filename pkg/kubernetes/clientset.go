package kubernetes

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"reflect"

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
		log.Fatalf("init kubernetes clientset error: %s", err)
	}
	return *clientset, nil
}

func GetInClusterConfig() rest.Config {
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatalf("init kubernetes config error: %s", err)
	}
	return *config
}

func GetOutClusterConfigTest1(kubeconfig string) rest.Config {
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = *flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = *flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()
	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatalf("init kubernetes config error: %s", err)
	}
	return *config
}

func GetOutClusterConfigWithContext(kubeconfig, context string) rest.Config {
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = *flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = *flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfig},
		&clientcmd.ConfigOverrides{
			CurrentContext: context,
		}).ClientConfig()
	if err != nil {
		log.Fatalf("init kubernetes config error: %s", err)
	}
	return *config
}

func GetOutClusterConfigTest(kubeconfig string) rest.Config {
	//config, err := clientcmd.RESTConfigFromKubeConfig([]byte(kubeconfig))
	//clientcmd.BuildConfigFromKubeconfigGetter("127.0.0.1", clientcmd.KubeconfigGetter())
	//if err != nil {
	//	log.Fatalf("init kubernetes config error: %s", err)
	//}
	tmpfile, err := os.CreateTemp("", "kubeconfig")
	if err != nil {
		log.Fatalf("init kubernetes config error1: %s", err)
	}
	defer os.Remove(tmpfile.Name())
	if err := os.WriteFile(tmpfile.Name(), []byte(kubeconfig), 0666); err != nil {
		log.Fatalf("init kubernetes config error2: %s", err)
	}
	config, err := clientcmd.BuildConfigFromFlags("https://172.30.120.220:6443", tmpfile.Name())
	if err != nil {
		log.Fatalf("init kubernetes config error: %s", err)
	}

	return *config
}

func GetOutClusterConfigTest2(kubeconfig string) rest.Config {
	config, err := clientcmd.RESTConfigFromKubeConfig([]byte(kubeconfig))
	if err != nil {
		log.Fatalf("init kubernetes config error: %s", err)
	}

	return *config
}

func GetOutClusterConfig(kubeconfig string) rest.Config {
	tmpfile, err := os.CreateTemp("./", "kubeconfig")
	if err != nil {
		log.Fatalf("convert kubeconfig to file error: %s", err)
	}
	defer os.Remove(tmpfile.Name())
	if err := os.WriteFile(tmpfile.Name(), []byte(kubeconfig), 0666); err != nil {
		log.Fatalf("write temp kubeconfig error: %s", err)
	}

	log.Println("temp kubeconfig at: ", tmpfile.Name())

	currentContext := ""

	if kc := clientcmd.GetConfigFromFileOrDie(tmpfile.Name()); kc != nil {
		if kc.CurrentContext != "" {
			currentContext = kc.CurrentContext
		} else {
			contexts := reflect.ValueOf(kc.Contexts).MapKeys()
			currentContext = contexts[0].String()
		}
	} else {
		log.Fatalf("init kubernetes config error: kubeconfig error")
	}

	config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: tmpfile.Name()},
		&clientcmd.ConfigOverrides{
			CurrentContext: currentContext,
		}).ClientConfig()

	if err != nil {
		log.Fatalf("init kubernetes config error: %s", err)
	}
	return *config
}
