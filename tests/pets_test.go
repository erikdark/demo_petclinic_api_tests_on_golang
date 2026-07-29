package tests

import (
 "net/http"
 "testing"

 "github.com/stretchr/testify/suite"

 "github.com/erikdark/demo_petclinic_api_tests_on_golang/internal/checkers"
 "github.com/erikdark/demo_petclinic_api_tests_on_golang/internal/fixtures"
)

type PetsSuite struct {
 BaseSuite
}

func TestPets(t *testing.T) {
 suite.Run(t, new(PetsSuite))
}

func (s *PetsSuite) TestAddPetToOwner() {
 t := s.T()

 owner := s.createOwner()
 in := fixtures.Pet(s.anyPetType())

 pet, resp, err := s.cli.CreatePet(s.ctx, owner.ID, in)
 checkers.StatusCode(t, resp, err, http.StatusCreated)
 s.Require().Greater(pet.ID, int32(0), "id не проставился")
 s.cleanup.AddPet(pet.ID)

 checkers.PetFields(t, in, pet)
 checkers.SameID(t, owner.ID, pet.OwnerID, "pet.ownerId")
 s.Assert().Empty(pet.Visits, "у нового питомца визитов быть не должно")

 got, resp, err := s.cli.GetOwner(s.ctx, owner.ID)
 checkers.StatusCode(t, resp, err, http.StatusOK)
 s.Require().Lenf(got.Pets, 1, "у владельца должен быть один питомец, а их %d", len(got.Pets))
 checkers.SameID(t, pet.ID, got.Pets[0].ID, "owner.pets[0]")
 checkers.PetFields(t, in, got.Pets[0])
}
