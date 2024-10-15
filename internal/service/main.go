package service

type Service struct {
	//clientset *kubernetes.Clientset
}

func (s *Service) DemoService() DemoService {
	return *NewDemoService()
}

func (s *Service) GlobalService() GlobalService {
	return *NewGlobalService()
}
