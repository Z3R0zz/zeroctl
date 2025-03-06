package utils

import (
	"fmt"
	"time"

	"github.com/valyala/fasthttp"
)

var client *fasthttp.Client

func init() {
	readTimeout, _ := time.ParseDuration("5000ms")
	writeTimeout, _ := time.ParseDuration("500ms")
	maxIdleConnDuration, _ := time.ParseDuration("1h")

	client = &fasthttp.Client{
		ReadTimeout:                   readTimeout,
		WriteTimeout:                  writeTimeout,
		MaxIdleConnDuration:           maxIdleConnDuration,
		NoDefaultUserAgentHeader:      true,
		DisableHeaderNamesNormalizing: true,
		DisablePathNormalizing:        true,
		Dial: (&fasthttp.TCPDialer{
			Concurrency:      4096,
			DNSCacheDuration: time.Hour,
		}).Dial,
	}
}

func Get(url string) (*fasthttp.Response, error) {
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(url)
	req.Header.SetMethod(fasthttp.MethodGet)
	resp := fasthttp.AcquireResponse()

	err := client.Do(req, resp)
	fasthttp.ReleaseRequest(req)

	if err != nil {
		fasthttp.ReleaseResponse(resp)
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	copiedResp := &fasthttp.Response{}
	resp.CopyTo(copiedResp)
	fasthttp.ReleaseResponse(resp)

	return copiedResp, nil
}
