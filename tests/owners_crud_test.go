package tests

import (
 "net/http"
 "testing"

 "github.com/stretchr/testify/suite"

 "github.com/erikdark/demo_petclinic_api_tests_on_golang/internal/models"
)

type OwnersCRUDSuite struct {
 BaseSuite
}

func TestOwnersCRUD(t *testing.T) {
 suite.Run(t, new(OwnersCRUDSuite))
}

// TestOwnerLifecycle — полный цикл: создать, прочитать, обновить, удалить
// и убедиться, что запись пропала.
func (s *OwnersCRUDSuite) TestOwnerLifecycle() {
 in := models.OwnerFields{
  FirstName: "Auto",
  LastName:  "Autoqacrud",
  Address:   "1 Test St.",
  City:      "Testcity",
  Telephone: "1234567890",
 }

 created, resp, err := s.cli.CreateOwner(s.ctx, in)
 s.Require().NoError(err)
 s.Require().Equal(http.StatusCreated, resp.StatusCode, "Создание владельца должно возвращать 201")
 s.Require().Greater(created.ID, int32(0), "Сервер должен присвоить положительный id")
 s.Assert().Equal(in.LastName, created.LastName, "lastName не совпадает с отправленным")
 s.Assert().Empty(created.Pets, "У нового владельца не должно быть питомцев")

 got, resp, err := s.cli.GetOwner(s.ctx, created.ID)
 s.Require().NoError(err)
 s.Require().Equal(http.StatusOK, resp.StatusCode)
 s.Assert().Equal(created.ID, got.ID)
 s.Assert().Equal(in.City, got.City, "city не совпадает с отправленным")

 updated := in
 updated.City = "Newcity"
 updated.Telephone = "9999999999"

 // PUT отвечает 204 без тела, поэтому результат сверяем следующим GET.
 resp, err = s.cli.UpdateOwner(s.ctx, created.ID, updated)
 s.Require().NoError(err)
 s.Require().Equal(http.StatusNoContent, resp.StatusCode, "Обновление должно возвращать 204")
 s.Assert().Empty(resp.Body, "На 204 тела быть не должно")

 got, resp, err = s.cli.GetOwner(s.ctx, created.ID)
 s.Require().NoError(err)
 s.Require().Equal(http.StatusOK, resp.StatusCode)
 s.Assert().Equal("Newcity", got.City, "Изменения не применились")
 s.Assert().Equal("9999999999", got.Telephone, "Изменения не применились")

 resp, err = s.cli.DeleteOwner(s.ctx, created.ID)
 s.Require().NoError(err)
 s.Require().Equal(http.StatusNoContent, resp.StatusCode, "Удаление должно возвращать 204")

 _, resp, err = s.cli.GetOwner(s.ctx, created.ID)
 s.Require().NoError(err)
 s.Require().Equal(http.StatusNotFound, resp.StatusCode, "После удаления владелец должен отдаваться как 404")
 s.Assert().Empty(resp.Body, "На 404 сервис отдаёт пустое тело")
}
