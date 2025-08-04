package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"example.com/m/v2/auth"
)

type SpotRequest struct {
	Method      string
	Url         string
	Body        io.Reader
	rawJson     any
	hasJson     bool
	headers     map[string]string
	queryParams map[string]string
}

func NewSpotRequest(method string, url string) *SpotRequest {
	return &SpotRequest{
		Method:      method,
		Url:         url,
		Body:        http.NoBody,
		hasJson:     false,
		headers:     map[string]string{},
		queryParams: map[string]string{},
	}
}

func (sr *SpotRequest) WithJson(body any) *SpotRequest {
	sr.rawJson = body
	sr.hasJson = true
	return sr
}

func (sr *SpotRequest) WithHeader(key, value string) *SpotRequest {
	sr.headers[key] = value
	return sr
}

func (sr *SpotRequest) WithQueryParam(key, value string) *SpotRequest {
	sr.queryParams[key] = value
	return sr
}

func (sr *SpotRequest) WithAuth(token *auth.FreshToken) *SpotRequest {
	sr.headers["Authorization"] = "Bearer " + token.AccessToken
	return sr
}

func (sr *SpotRequest) Do() (*http.Response, error) {
	var err error
	if sr.hasJson {
		var data []byte
		data, err = json.Marshal(sr.rawJson)
		if err != nil {
			return nil, err
		}
		sr.Body = bytes.NewReader(data)
		sr.headers["Content-Type"] = "application/json"
	}

	if len(sr.queryParams) > 0 {
		var builder strings.Builder
		builder.WriteString(sr.Url)
		builder.WriteRune('?')
		for key, value := range sr.queryParams {
			builder.WriteString(key)
			builder.WriteRune('=')
			builder.WriteString(value)
			builder.WriteRune('&')
		}
		realUrl := builder.String()
		sr.Url = realUrl[:len(realUrl)-1]
	}

	req, err := http.NewRequest(sr.Method, sr.Url, sr.Body)
	if err != nil {
		return nil, err
	}

	for key, val := range sr.headers {
		req.Header.Set(key, val)
	}

	return http.DefaultClient.Do(req)
}

func ValidateResponse(resp *http.Response, err error) error {
	if err != nil {
		return err
	}

	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf(resp.Status + ":" + string(body))
	}

	return nil
}
