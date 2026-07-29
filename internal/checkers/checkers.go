package checkers

import (
 "fmt"
 "net/http"
 "testing"

 "github.com/stretchr/testify/assert"
 "github.com/stretchr/testify/require"

 "github.com/erikdark/demo_petclinic_api_tests_on_golang/internal/client"
 "github.com/erikdark/demo_petclinic_api_tests_on_golang/internal/models"
)

func StatusCode(t *testing.T, resp *client.Response, err error, want int) {
 t.Helper()

 require.NoError(t, err, "запрос не ушёл")
 require.Equalf(t, want, resp.StatusCode, "ждал %d, пришло %d, тело: %s", want, resp.StatusCode, resp.Body)
}

func NotFoundEmptyBody(t *testing.T, resp *client.Response, err error) {
 t.Helper()

 StatusCode(t, resp, err, http.StatusNotFound)
 assert.Emptyf(t, resp.Body, "на 404 тело должно быть пустое, а пришло: %s", resp.Body)
}

func NoContent(t *testing.T, resp *client.Response, err error) {
 t.Helper()

 StatusCode(t, resp, err, http.StatusNoContent)
 assert.Emptyf(t, resp.Body, "на 204 тела быть не должно, а пришло: %s", resp.Body)
}

func OwnerFields(t *testing.T, want models.OwnerFields, got models.Owner) {
 t.Helper()

 assert.Equal(t, want.FirstName, got.FirstName, "firstName не тот, что отправляли")
 assert.Equal(t, want.LastName, got.LastName, "lastName не тот, что отправляли")
 assert.Equal(t, want.Address, got.Address, "address не тот, что отправляли")
 assert.Equal(t, want.City, got.City, "city не тот, что отправляли")
 assert.Equal(t, want.Telephone, got.Telephone, "telephone не тот, что отправляли")
}

func PetFields(t *testing.T, want models.PetFields, got models.Pet) {
 t.Helper()

 assert.Equal(t, want.Name, got.Name, "name не тот, что отправляли")
 assert.Equal(t, want.BirthDate, got.BirthDate, "birthDate не тот, что отправляли")
 assert.Equal(t, want.Type.ID, got.Type.ID, "type.id не тот, что отправляли")
 assert.Equal(t, want.Type.Name, got.Type.Name, "type.name не тот, что отправляли")
}

func OwnerInList(t *testing.T, owners []models.Owner, id int32) {
 t.Helper()

 for _, o := range owners {
  if o.ID == id {
   return
  }
 }

 assert.Failf(t, "нашего владельца нет в выдаче", "id %d не найден среди %d записей", id, len(owners))
}

func ValidationProblem(t *testing.T, resp *client.Response, wantFields ...string) models.ProblemDetail {
 t.Helper()

 problem, err := client.Problem(resp)
 require.NoErrorf(t, err, "тело ошибки не парсится, пришло: %s", resp.Body)

 assert.Containsf(t, resp.Header.Get("Content-Type"), "problem+json", "ошибка должна приходить как problem+json, а пришёл %s", resp.Header.Get("Content-Type"))
 assert.Equal(t, resp.StatusCode, problem.Status, "status в теле не совпал с http-статусом")
 require.NotEmpty(t, problem.SchemaValidationErrors, "сервис не вернул ни одной ошибки валидации")

 for _, field := range wantFields {
  assert.Containsf(t, problem.FieldsInError(), field, "сервис должен ругнуться на %q, а ругнулся на %v", field, problem.FieldsInError())
 }

 return problem
}

func SameID(t *testing.T, want, got int32, what string) {
 t.Helper()

 assert.Equal(t, want, got, fmt.Sprintf("%s: id другой", what))
}
