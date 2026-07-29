package tests

import (
 "context"
 "net/http"

 "github.com/stretchr/testify/suite"

 "github.com/erikdark/demo_petclinic_api_tests_on_golang/config"
 "github.com/erikdark/demo_petclinic_api_tests_on_golang/internal/client"
 "github.com/erikdark/demo_petclinic_api_tests_on_golang/internal/fixtures"
 "github.com/erikdark/demo_petclinic_api_tests_on_golang/internal/models"
)

type BaseSuite struct {
 suite.Suite

 cli     *client.Client
 ctx     context.Context
 cleanup *fixtures.Cleanup
}

func (s *BaseSuite) SetupSuite() {
 cfg := config.Load()

 s.cli = client.New(cfg.BaseURL, cfg.Timeout)
 s.ctx = context.Background()
 s.cleanup = fixtures.NewCleanup(s.cli)
}

func (s *BaseSuite) TearDownSuite() {
 s.cleanup.Flush(s.ctx, s.T())
}

func (s *BaseSuite) createOwner() models.Owner {
 s.T().Helper()

 owner, resp, err := s.cli.CreateOwner(s.ctx, fixtures.Owner())
 s.Require().NoError(err)
 s.Require().Equalf(http.StatusCreated, resp.StatusCode, "владелец для теста не создался, тело: %s", resp.Body)

 s.cleanup.AddOwner(owner.ID)

 return owner
}

func (s *BaseSuite) anyPetType() models.PetType {
 s.T().Helper()

 types, resp, err := s.cli.ListPetTypes(s.ctx)
 s.Require().NoError(err)
 s.Require().Equal(http.StatusOK, resp.StatusCode)
 s.Require().NotEmpty(types, "справочник типов пустой, не из чего выбрать")

 return types[0]
}
