package fixtures

import (
 "context"
 "net/http"
 "sync"
 "testing"

 "github.com/erikdark/demo_petclinic_api_tests_on_golang/internal/client"
)

type Cleanup struct {
 cli *client.Client

 mu     sync.Mutex
 pets   []int32
 owners []int32
}

func NewCleanup(cli *client.Client) *Cleanup {
 return &Cleanup{cli: cli}
}

func (c *Cleanup) AddOwner(id int32) {
 c.mu.Lock()
 defer c.mu.Unlock()

 c.owners = append(c.owners, id)
}

func (c *Cleanup) AddPet(id int32) {
 c.mu.Lock()
 defer c.mu.Unlock()

 c.pets = append(c.pets, id)
}

func (c *Cleanup) Flush(ctx context.Context, t *testing.T) {
 t.Helper()

 c.mu.Lock()
 pets, owners := c.pets, c.owners
 c.pets, c.owners = nil, nil
 c.mu.Unlock()

 for _, id := range pets {
  resp, err := c.cli.DeletePet(ctx, id)
  check(t, "pet", id, resp, err)
 }

 for _, id := range owners {
  resp, err := c.cli.DeleteOwner(ctx, id)
  check(t, "owner", id, resp, err)
 }
}

func check(t *testing.T, kind string, id int32, resp *client.Response, err error) {
 t.Helper()

 if err != nil {
  t.Errorf("не смог удалить %s %d: %v", kind, id, err)
  return
 }

 if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusNotFound {
  t.Errorf("удаление %s %d вернуло %d, ждал 204 или 404", kind, id, resp.StatusCode)
 }
}
