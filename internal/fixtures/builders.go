package fixtures

import (
 "fmt"
 "math/rand/v2"

 "github.com/erikdark/demo_petclinic_api_tests_on_golang/internal/models"
)

const (
 lastNamePrefix = "Autoqa"
 suffixLen      = 8
)

func Owner() models.OwnerFields {
 return models.OwnerFields{
  FirstName: "Auto",
  LastName:  UniqueLastName(),
  Address:   "1 Test St.",
  City:      "Testcity",
  Telephone: Telephone(),
 }
}

func UniqueLastName() string {
 suffix := make([]byte, suffixLen)
 for i := range suffix {
  suffix[i] = byte(97 + rand.IntN(26))
 }

 return lastNamePrefix + string(suffix)
}

func Telephone() string {
 return fmt.Sprintf("%010d", rand.IntN(1000000000))
}

func Pet(petType models.PetType) models.PetFields {
 return models.PetFields{
  Name:      fmt.Sprintf("Pet-%d", rand.IntN(1000000)),
  BirthDate: "2020-01-01",
  Type:      petType,
 }
}
