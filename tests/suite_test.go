package tests

import (
 "context"

 "github.com/stretchr/testify/suite"

 "github.com/erikdark/demo_petclinic_api_tests_on_golang/config"
 "github.com/erikdark/demo_petclinic_api_tests_on_golang/internal/client"
)

// BaseSuite — общая обвязка: клиент и контекст для всех сценариев.
type BaseSuite struct {
 suite.Suite

 cli *client.Client
 ctx context.Context
}

func (s *BaseSuite) SetupSuite() {
 cfg := config.Load()

 s.cli = client.New(cfg.BaseURL, cfg.Timeout)
 s.ctx = context.Background()
}
