package kube

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"reflect"

	"github.com/gatepoint/gatepoint/pkg/config"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
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
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			log.Fatalf("remove temp kubeconfig file error: %s", err)
		}
	}(tmpfile.Name())
	if err := os.WriteFile(tmpfile.Name(), []byte(kubeconfig), 0666); err != nil {
		log.Fatalf("init kubernetes config error: %s", err)
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
	tmpfile, err := os.CreateTemp("", "kubeconfig")
	if err != nil {
		log.Fatalf("convert kubeconfig to file error: %s", err)
	}

	defer func(f *os.File) {
		name := f.Name()
		err := f.Close()
		if err != nil {
			log.Fatalf("close temp kubeconfig error: %s", err)
		}
		if err := os.Remove(name); err != nil {
			log.Fatalf("remove temp kubeconfig error: %s", err)
		}
	}(tmpfile)

	if err := os.WriteFile(tmpfile.Name(), []byte(kubeconfig), 0777); err != nil {
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
