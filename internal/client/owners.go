package client

import (
 "context"
 "encoding/json"
 "fmt"
 "net/http"
 "net/url"

 "github.com/erikdark/demo_petclinic_api_tests_on_golang/internal/models"
)

func (c *Client) CreateOwner(ctx context.Context, in models.OwnerFields) (models.Owner, *Response, error) {
 var owner models.Owner

 resp, err := c.do(ctx, http.MethodPost, "/owners", nil, in)
 if err != nil {
  return owner, nil, err
 }

 return owner, resp, decode(resp, &owner)
}

func (c *Client) CreateOwnerRaw(ctx context.Context, payload []byte) (*Response, error) {
 return c.doRaw(ctx, http.MethodPost, "/owners", nil, payload, true)
}

func (c *Client) GetOwner(ctx context.Context, id int32) (models.Owner, *Response, error) {
 var owner models.Owner

 resp, err := c.do(ctx, http.MethodGet, ownerPath(id), nil, nil)
 if err != nil {
  return owner, nil, err
 }

 return owner, resp, decode(resp, &owner)
}


func (c *Client) ListOwners(ctx context.Context, lastName string) ([]models.Owner, *Response, error) {
 var (
  owners []models.Owner
  query  url.Values
 )

 if lastName != "" {
  query = url.Values{"lastName": []string{lastName}}
 }

 resp, err := c.do(ctx, http.MethodGet, "/owners", query, nil)
 if err != nil {
  return nil, nil, err
 }

 return owners, resp, decode(resp, &owners)
}


func (c *Client) UpdateOwner(ctx context.Context, id int32, in models.OwnerFields) (*Response, error) {
 return c.do(ctx, http.MethodPut, ownerPath(id), nil, in)
}

func (c *Client) DeleteOwner(ctx context.Context, id int32) (*Response, error) {
 return c.do(ctx, http.MethodDelete, ownerPath(id), nil, nil)
}


func Problem(resp *Response) (models.ProblemDetail, error) {
 var problem models.ProblemDetail
 if len(resp.Body) == 0 {
  return problem, nil
 }

 return problem, json.Unmarshal(resp.Body, &problem)
}

func ownerPath(id int32) string {
 return fmt.Sprintf("/owners/%d", id)
}
