package tests

import (
 "net/http"
 "strings"
 "testing"

 "github.com/stretchr/testify/suite"

 "github.com/erikdark/demo_petclinic_api_tests_on_golang/internal/checkers"
 "github.com/erikdark/demo_petclinic_api_tests_on_golang/internal/fixtures"
)

type OwnersListSuite struct {
 BaseSuite
}

func TestOwnersList(t *testing.T) {
 suite.Run(t, new(OwnersListSuite))
}

func (s *OwnersListSuite) TestSearchByLastName() {
 t := s.T()
 owner := s.createOwner()

 found, resp, err := s.cli.ListOwners(s.ctx, owner.LastName)
 checkers.StatusCode(t, resp, err, http.StatusOK)

 s.Require().Lenf(found, 1, "в поиске должен быть один владелец, а не %d", len(found))
 checkers.SameID(t, owner.ID, found[0].ID, "owner")
 checkers.OwnerFields(t, owner.OwnerFields, found[0])
}

func (s *OwnersListSuite) TestListContainsCreated() {
 t := s.T()
 owner := s.createOwner()

 all, resp, err := s.cli.ListOwners(s.ctx, "")
 checkers.StatusCode(t, resp, err, http.StatusOK)

 s.Require().NotEmpty(all, "список владельцев пустой")
 checkers.OwnerInList(t, all, owner.ID)
}

func (s *OwnersListSuite) TestSearchSemantics() {
 owner := s.createOwner()

 cases := map[string]struct {
  lastName   string
  wantStatus int
  wantFound  bool
 }{
  "полное совпадение":      {lastName: owner.LastName, wantStatus: http.StatusOK, wantFound: true},
  "префикс фамилии":        {lastName: owner.LastName[:len(owner.LastName)-3], wantStatus: http.StatusOK, wantFound: true},
  "другой регистр":         {lastName: strings.ToLower(owner.LastName), wantStatus: http.StatusNotFound},
  "подстрока в середине":   {lastName: owner.LastName[3:], wantStatus: http.StatusNotFound},
  "несуществующая фамилия": {lastName: fixtures.UniqueLastName(), wantStatus: http.StatusNotFound},
 }

 for name, tc := range cases {
  s.Run(name, func() {
   t := s.T()

   found, resp, err := s.cli.ListOwners(s.ctx, tc.lastName)
   checkers.StatusCode(t, resp, err, tc.wantStatus)

   if tc.wantFound {
    checkers.OwnerInList(t, found, owner.ID)
    return
   }

   s.Assert().Emptyf(resp.Body, "на пустую выдачу приходит 404 без тела, а пришло: %s", resp.Body)
  })
 }
}
