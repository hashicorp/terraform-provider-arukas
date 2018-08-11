package arukas

import (
	"fmt"
	"time"

	"github.com/yamamoto-febc/go-arukas"
)

const (
	JSONTokenParamName   = "ARUKAS_JSON_API_TOKEN"
	JSONSecretParamName  = "ARUKAS_JSON_API_SECRET"
	JSONUrlParamName     = "ARUKAS_JSON_API_URL"
	JSONDebugParamName   = "ARUKAS_DEBUG"
	JSONTimeoutParamName = "ARUKAS_TIMEOUT"
	userAgentFormat      = "terraform-provider-arukas(go-arukas: v%s)"
)

type Config struct {
	Token   string
	Secret  string
	URL     string
	Trace   string
	Timeout int
}

func (c *Config) NewClient() (*arukasClient, error) {

	timeout := time.Duration(0)
	if c.Timeout > 0 {
		timeout = time.Duration(c.Timeout) * time.Second
	}

	client, err := arukas.NewClient(&arukas.ClientParam{
		Token:      c.Token,
		Secret:     c.Secret,
		APIBaseURL: c.URL,
		Trace:      c.Trace != "",
		UserAgent:  fmt.Sprintf(userAgentFormat, arukas.Version),
	})
	if err != nil {
		return nil, err
	}

	return &arukasClient{
		Client:  client,
		timeout: timeout,
	}, nil
}

type arukasClient struct {
	arukas.Client
	timeout time.Duration
}
