package client

import (
 "context"
 "fmt"
 "net/http"

 "github.com/erikdark/demo_petclinic_api_tests_on_golang/internal/models"
)

func (c *Client) CreatePet(ctx context.Context, ownerID int32, in models.PetFields) (models.Pet, *Response, error) {
 var pet models.Pet

 resp, err := c.do(ctx, http.MethodPost, ownerPetsPath(ownerID), nil, in)
 if err != nil {
  return pet, nil, err
 }

 return pet, resp, decode(resp, &pet)
}

func (c *Client) CreatePetRaw(ctx context.Context, ownerID int32, payload []byte) (*Response, error) {
 return c.doRaw(ctx, http.MethodPost, ownerPetsPath(ownerID), nil, payload, true)
}

func (c *Client) GetPet(ctx context.Context, petID int32) (models.Pet, *Response, error) {
 var pet models.Pet

 resp, err := c.do(ctx, http.MethodGet, petPath(petID), nil, nil)
 if err != nil {
  return pet, nil, err
 }

 return pet, resp, decode(resp, &pet)
}

func (c *Client) DeletePet(ctx context.Context, petID int32) (*Response, error) {
 return c.do(ctx, http.MethodDelete, petPath(petID), nil, nil)
}

func (c *Client) ListPetTypes(ctx context.Context) ([]models.PetType, *Response, error) {
 var types []models.PetType

 resp, err := c.do(ctx, http.MethodGet, "/pettypes", nil, nil)
 if err != nil {
  return nil, nil, err
 }

 return types, resp, decode(resp, &types)
}

func petPath(id int32) string {
 return fmt.Sprintf("/pets/%d", id)
}

func ownerPetsPath(ownerID int32) string {
 return fmt.Sprintf("/owners/%d/pets", ownerID)
}
