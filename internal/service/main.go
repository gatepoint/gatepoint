package service

type Service struct {
}

func (s *Service) DemoService() DemoService {
	return *NewDemoService()
}
