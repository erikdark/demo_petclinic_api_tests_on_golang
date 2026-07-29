package tests

import (
 "net/http"
 "testing"

 "github.com/stretchr/testify/suite"

 "github.com/erikdark/demo_petclinic_api_tests_on_golang/internal/checkers"
 "github.com/erikdark/demo_petclinic_api_tests_on_golang/internal/fixtures"
)

type OwnersCRUDSuite struct {
 BaseSuite
}

func TestOwnersCRUD(t *testing.T) {
 suite.Run(t, new(OwnersCRUDSuite))
}

func (s *OwnersCRUDSuite) TestOwnerLifecycle() {
 t := s.T()
 in := fixtures.Owner()

 created, resp, err := s.cli.CreateOwner(s.ctx, in)
 checkers.StatusCode(t, resp, err, http.StatusCreated)
 s.Require().Greater(created.ID, int32(0), "id не проставился")
 s.cleanup.AddOwner(created.ID)

 checkers.OwnerFields(t, in, created)
 s.Assert().Empty(created.Pets, "у нового владельца питомцев быть не должно")

 got, resp, err := s.cli.GetOwner(s.ctx, created.ID)
 checkers.StatusCode(t, resp, err, http.StatusOK)
 checkers.SameID(t, created.ID, got.ID, "owner")
 checkers.OwnerFields(t, in, got)

 updated := in
 updated.City = "Newcity"
 updated.Telephone = fixtures.Telephone()

 resp, err = s.cli.UpdateOwner(s.ctx, created.ID, updated)
 checkers.NoContent(t, resp, err)

 got, resp, err = s.cli.GetOwner(s.ctx, created.ID)
 checkers.StatusCode(t, resp, err, http.StatusOK)
 checkers.OwnerFields(t, updated, got)

 resp, err = s.cli.DeleteOwner(s.ctx, created.ID)
 checkers.NoContent(t, resp, err)

 _, resp, err = s.cli.GetOwner(s.ctx, created.ID)
 checkers.NotFoundEmptyBody(t, resp, err)
}
