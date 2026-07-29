package client

import (
 "bytes"
 "context"
 "encoding/json"
 "io"
 "net/http"
 "net/url"
 "time"
)


type Response struct {
 StatusCode int
 Body       []byte
 Header     http.Header
}

type Client struct {
 http    *http.Client
 baseURL string
}

func New(baseURL string, timeout time.Duration) *Client {
 return &Client{
  http:    &http.Client{Timeout: timeout},
  baseURL: baseURL,
 }
}

func (c *Client) do(ctx context.Context, method, path string, query url.Values, body interface{}) (*Response, error) {
 var payload []byte
 if body != nil {
  raw, err := json.Marshal(body)
  if err != nil {
   return nil, err
  }

  payload = raw
 }

 return c.doRaw(ctx, method, path, query, payload, payload != nil)
}


func (c *Client) doRaw(ctx context.Context, method, path string, query url.Values, payload []byte, withContentType bool) (*Response, error) {
 var body io.Reader
 if payload != nil {
  body = bytes.NewReader(payload)
 }

 u := c.baseURL + path
 if len(query) > 0 {
  u += "?" + query.Encode()
 }

 req, err := http.NewRequestWithContext(ctx, method, u, body)
 if err != nil {
  return nil, err
 }

 req.Header.Set("Accept", "application/json")
 if withContentType {
  req.Header.Set("Content-Type", "application/json")
 }

 resp, err := c.http.Do(req)
 if err != nil {
  return nil, err
 }
 defer resp.Body.Close()

 raw, err := io.ReadAll(resp.Body)
 if err != nil {
  return nil, err
 }

 return &Response{StatusCode: resp.StatusCode, Body: raw, Header: resp.Header}, nil
}


func decode(resp *Response, target interface{}) error {
 if resp.StatusCode >= http.StatusMultipleChoices || len(resp.Body) == 0 {
  return nil
 }

 return json.Unmarshal(resp.Body, target)
}
