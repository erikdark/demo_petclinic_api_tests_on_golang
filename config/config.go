package config

import (
 "os"
 "time"
)

const (
 baseURLEnv = "PETCLINIC_BASE_URL"
 timeoutEnv = "PETCLINIC_TIMEOUT"

 defaultBaseURL = "http://localhost:9966/petclinic/api"
 defaultTimeout = 10 * time.Second
)


type Config struct {
 BaseURL string
 Timeout time.Duration
}


func Load() Config {
 cfg := Config{
  BaseURL: defaultBaseURL,
  Timeout: defaultTimeout,
 }

 if v := os.Getenv(baseURLEnv); v != "" {
  cfg.BaseURL = v
 }

 if v := os.Getenv(timeoutEnv); v != "" {
  if d, err := time.ParseDuration(v); err == nil {
   cfg.Timeout = d
  }
 }

 return cfg
}