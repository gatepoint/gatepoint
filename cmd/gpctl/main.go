package main

import (
	_ "context"
	"fmt"
	"time"

	k8scorev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	_ "k8s.io/apimachinery/pkg/util/runtime"
	corev1informers "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	// 加载 kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", "/Users/yangyang/Desktop/环境/demo-dev-new.yaml")
	if err != nil {
		panic(err)
	}

	// 创建 Kubernetes 客户端
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	// 创建 ConfigMap informer
	informer := corev1informers.NewFilteredConfigMapInformer(
		clientset,
		"",
		30*time.Second, // resync period
		cache.Indexers{
			cache.NamespaceIndex: cache.MetaNamespaceIndexFunc,
		},
		func(options *metav1.ListOptions) {
			//options.LabelSelector = "app=myapp" // 根据标签过滤
		},
	)

	// 处理 ConfigMap 事件
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			cm := obj.(*k8scorev1.ConfigMap)
			fmt.Println("ConfigMap added:", cm.Name)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			oldCM := oldObj.(*k8scorev1.ConfigMap)
			newCM := newObj.(*k8scorev1.ConfigMap)
			fmt.Println("ConfigMap updated:", oldCM.Name, "->", newCM.Name)
		},
		DeleteFunc: func(obj interface{}) {
			cm := obj.(*k8scorev1.ConfigMap)
			fmt.Println("ConfigMap deleted:", cm.Name)
		},
	})

	// 启动 informer
	stopCh := make(chan struct{})
	defer close(stopCh)
	go informer.Run(stopCh)

	for {
		getConfigMapsWithIndex(informer)

		time.Sleep(10 * time.Second)
	}

	//
	//
	//// 阻塞主 goroutine
	//<-stopCh
}

func getConfigMapsWithIndex(informer cache.SharedIndexInformer) {
	index := informer.GetIndexer()
	//a := index.List()
	//fmt.Println(a)
	objs, err := index.ByIndex(cache.NamespaceIndex, "skoala-system")
	if err != nil {
		panic(err)
	}
	for _, obj := range objs {
		//cm := obj.(*corev1.ConfigMap)
		fmt.Println("Found ConfigMap:", obj)
	}
}
