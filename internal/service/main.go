package service

type Service struct {
	//clientset *kubernetes.Clientset
}

func (s *Service) DemoService() DemoService {
	return *NewDemoService()
}
