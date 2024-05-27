package service

import (
	"context"
	generalv1 "github.com/gatepoint/gatepoint/api/general/v1"
	projectv1 "github.com/gatepoint/gatepoint/api/gatepoint/v1"
)

type DemoService struct {
	// This is generated by protoc
	projectv1.UnimplementedDemoServer
}

func NewDemoService() *DemoService {
	return &DemoService{}
}

func (s *DemoService) Demo(ctx context.Context, req *generalv1.DemoRequest) (*generalv1.DemoResponse, error) {
	return &generalv1.DemoResponse{
		Demo: &generalv1.Demo{
			Demo: req.Demo,
		},
	}, nil
}
