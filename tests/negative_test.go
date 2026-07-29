package tests

import (
 "net/http"
 "testing"

 "github.com/stretchr/testify/suite"

 "github.com/erikdark/demo_petclinic_api_tests_on_golang/internal/checkers"
 "github.com/erikdark/demo_petclinic_api_tests_on_golang/internal/client"
 "github.com/erikdark/demo_petclinic_api_tests_on_golang/internal/fixtures"
 "github.com/erikdark/demo_petclinic_api_tests_on_golang/internal/models"
)

const missingID = int32(999999)

type NegativeSuite struct {
 BaseSuite
}

func TestNegative(t *testing.T) {
 suite.Run(t, new(NegativeSuite))
}

func (s *NegativeSuite) TestMissingEntities() {
 cases := map[string]func() (*client.Response, error){
  "получить владельца": func() (*client.Response, error) {
   _, resp, err := s.cli.GetOwner(s.ctx, missingID)
   return resp, err
  },
  "обновить владельца": func() (*client.Response, error) {
   return s.cli.UpdateOwner(s.ctx, missingID, fixtures.Owner())
  },
  "удалить владельца": func() (*client.Response, error) {
   return s.cli.DeleteOwner(s.ctx, missingID)
  },
  "получить питомца": func() (*client.Response, error) {
   _, resp, err := s.cli.GetPet(s.ctx, missingID)
   return resp, err
  },
  "добавить питомца несуществующему владельцу": func() (*client.Response, error) {
   _, resp, err := s.cli.CreatePet(s.ctx, missingID, fixtures.Pet(s.anyPetType()))
   return resp, err
  },
 }

 for name, call := range cases {
  s.Run(name, func() {
   resp, err := call()
   checkers.NotFoundEmptyBody(s.T(), resp, err)
  })
 }
}

func (s *NegativeSuite) TestOwnerValidation() {
 badPhone := fixtures.Owner()
 badPhone.Telephone = "not-a-number"

 noFirstName := fixtures.Owner()
 noFirstName.FirstName = ""

 digitsInLastName := fixtures.Owner()
 digitsInLastName.LastName = "Auto12345"

 cases := map[string]struct {
  payload    models.OwnerFields
  wantFields []string
 }{
  "телефон не из цифр": {payload: badPhone, wantFields: []string{"telephone"}},
  "пустое имя":         {payload: noFirstName, wantFields: []string{"firstName"}},
  "цифры в фамилии":    {payload: digitsInLastName, wantFields: []string{"lastName"}},
  "все поля пустые": {
   payload:    models.OwnerFields{},
   wantFields: []string{"firstName", "lastName", "address", "city", "telephone"},
  },
 }

 for name, tc := range cases {
  s.Run(name, func() {
   t := s.T()

   _, resp, err := s.cli.CreateOwner(s.ctx, tc.payload)
   checkers.StatusCode(t, resp, err, http.StatusBadRequest)
   checkers.ValidationProblem(t, resp, tc.wantFields...)
  })
 }
}

func (s *NegativeSuite) TestPetWithoutType() {
 t := s.T()
 owner := s.createOwner()

 resp, err := s.cli.CreatePetRaw(s.ctx, owner.ID, []byte(`{"name":"NoType","birthDate":"2020-01-01","typeId":1}`))
 checkers.StatusCode(t, resp, err, http.StatusBadRequest)
 checkers.ValidationProblem(t, resp, "type")
}

func (s *NegativeSuite) TestUnknownPetType() {
 t := s.T()
 owner := s.createOwner()

 _, resp, err := s.cli.CreatePet(s.ctx, owner.ID, fixtures.Pet(models.PetType{ID: missingID, Name: "unknown"}))
 checkers.StatusCode(t, resp, err, http.StatusNotFound)

 problem, err := client.Problem(resp)
 s.Require().NoError(err)
 s.Assert().Equalf("DataIntegrityViolationException", problem.Title, "ждал ошибку по несуществующему типу, а пришло: %s", resp.Body)
}

func (s *NegativeSuite) TestMalformedBody() {
 t := s.T()

 resp, err := s.cli.CreateOwnerRaw(s.ctx, []byte(`{"firstName":`))
 checkers.StatusCode(t, resp, err, http.StatusInternalServerError)

 problem, err := client.Problem(resp)
 s.Require().NoError(err)
 s.Assert().Equalf("HttpMessageNotReadableException", problem.Title, "ждал ошибку разбора json, а пришло: %s", resp.Body)
}

