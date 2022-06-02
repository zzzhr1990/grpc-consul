package server

import (
	"time"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

type SimpleConsulServer struct {
	// grpcServer  *grpc.Server
	healthcheck *health.Server
	res         *ConsulResult
}

func (s *SimpleConsulServer) BaseInit(cfg *ConsulRegisterConfig, grpcServer *grpc.Server) error {
	err := s.initHealthCheck(cfg, grpcServer)
	if err != nil {
		return errors.Wrap(err, "init health check failed")
	}
	err = s.initConsul(cfg)
	if err != nil {
		return errors.Wrap(err, "init consul failed")
	}
	return nil
}

func (s *SimpleConsulServer) BaseShutdown() error {
	if s.healthcheck != nil {
		s.healthcheck.Shutdown()
	}
	time.Sleep(time.Second * 5)
	if s.res != nil {
		err := s.res.ShutdownAgent()
		if err != nil {
			return errors.Wrap(err, "shutdown consul failed")
		}
	}
	time.Sleep(time.Second * 3)
	return nil
}

func (s *SimpleConsulServer) initHealthCheck(cfg *ConsulRegisterConfig, grpcServer *grpc.Server) error {
	s.healthcheck = health.NewServer()
	healthpb.RegisterHealthServer(grpcServer, s.healthcheck)
	s.healthcheck.SetServingStatus(cfg.Name, healthpb.HealthCheckResponse_SERVING)
	return nil
}

func (s *SimpleConsulServer) initConsul(cfg *ConsulRegisterConfig) error {
	// consulServer.dd
	//hname, _ := os.Hostname()
	var err error = nil
	s.res, err = RegisterToConsul(cfg)
	return err
}
