package service

import "k8s.io/client-go/kubernetes"

type Service struct {
	clientset *kubernetes.Clientset
}

func (s *Service) DemoService() DemoService {
	return *NewDemoService()
}
