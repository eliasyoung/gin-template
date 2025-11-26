package client

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"time"
)

type HttpClient struct {
	httpClient *http.Client
}

func NewHttpClient(insecureSkipVerify bool) *HttpClient {
	// 测试客户端
	tlsConfig := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: insecureSkipVerify},
	}

	httpClient := &http.Client{
		Timeout:   30 * time.Second,
		Transport: tlsConfig, //测试时使用
	}

	return &HttpClient{
		httpClient: httpClient,
	}
}

func (hc *HttpClient) DoGetRequest(url string, headerMap map[string]string) ([]byte, int, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, 0, fmt.Errorf("create GET request: %w", err)
	}

	if len(headerMap) > 0 {
		for headerKey, headerValue := range headerMap {
			req.Header.Set(headerKey, headerValue)
		}
	}

	resp, err := hc.httpClient.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("send GET request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, resp.StatusCode, fmt.Errorf("GET request failed with http status code : %d", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("read response: %w", err)
	}

	return bodyBytes, resp.StatusCode, nil
}

func (hc *HttpClient) DoGetRequestWithCtx(ctx context.Context, url string, headerMap map[string]string) ([]byte, int, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, 0, fmt.Errorf("create GET request: %w", err)
	}

	if len(headerMap) > 0 {
		for headerKey, headerValue := range headerMap {
			req.Header.Set(headerKey, headerValue)
		}
	}

	resp, err := hc.httpClient.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("send GET request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, resp.StatusCode, fmt.Errorf("GET request failed with http status code: %d", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("read response: %w", err)
	}

	return bodyBytes, resp.StatusCode, nil
}

func (hc *HttpClient) DoPostRequest(url string, body []byte, headerMap map[string]string) ([]byte, int, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, 0, fmt.Errorf("create POST request: %w", err)
	}

	if len(headerMap) > 0 {
		for headerKey, headerValue := range headerMap {
			req.Header.Set(headerKey, headerValue)
		}
	}

	resp, err := hc.httpClient.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("send POST request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, resp.StatusCode, fmt.Errorf("POST request failed with http status code: %d", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("read response: %w", err)
	}

	return bodyBytes, resp.StatusCode, nil
}

func (hc *HttpClient) DoPostRequestWithCtx(ctx context.Context, url string, body []byte, headerMap map[string]string) ([]byte, int, error) {
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, 0, fmt.Errorf("create POST request: %w", err)
	}

	if len(headerMap) > 0 {
		for headerKey, headerValue := range headerMap {
			req.Header.Set(headerKey, headerValue)
		}
	}

	resp, err := hc.httpClient.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("send POST request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, resp.StatusCode, fmt.Errorf("POST request failed with http status code: %d", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("read response: %w", err)
	}

	return bodyBytes, resp.StatusCode, nil
}
