package OneAPISDK

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type OneApiClient struct {
	serverUrl   string
	httpClient  *http.Client
	systemToken string
	seededRand  *rand.Rand
}

type OneApiOptions struct {
	ServerUrl       string
	HttpTimeout     time.Duration
	SystemToken     string
	TransportConfig http.RoundTripper
}

type tokenTransport struct {
	transport http.RoundTripper
	token     string
}

func NewOneApiClient(options *OneApiOptions) *OneApiClient {
	client := &OneApiClient{
		serverUrl: options.ServerUrl,
		httpClient: &http.Client{
			Transport: options.TransportConfig,
			Timeout:   options.HttpTimeout,
		},
		systemToken: options.SystemToken,
		seededRand:  rand.New(rand.NewSource(time.Now().UnixNano())),
	}
	return client
}

func (c *OneApiClient) sendReq(endpoint string, method string, token string, sessionID string, reqData interface{}, additionHeader map[string]string) (*CommonAPIRes, error) {
	var jsonData []byte
	callUrl := fmt.Sprintf("%s%s", c.serverUrl, endpoint)
	u, err := url.Parse(callUrl)
	if err != nil {
		return nil, fmt.Errorf("error Parse url: %v", err)
	}

	if method == http.MethodGet && reqData != nil {
		queryParams := reqData.(map[string]string)
		q := u.Query()
		for key, value := range queryParams {
			q.Add(key, value)
		}
		u.RawQuery = q.Encode()
	} else {
		jsonData, err = json.Marshal(reqData)
		if err != nil {
			return nil, fmt.Errorf("error Marshal data: %v", err)
		}
	}

	req, err := http.NewRequest(method, u.String(), bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	if token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	}

	if sessionID != "" {
		req.AddCookie(&http.Cookie{Name: "session", Value: sessionID})
	}

	if additionHeader != nil {
		for k, v := range additionHeader {
			req.Header.Set(k, v)
		}
	}

	if method != http.MethodGet {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error call api: %v", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error read resp: %v", err)
	}

	responseData := &CommonAPIRes{
		res: &oneAPIRes{},
	}

	for _, cookie := range resp.Cookies() {
		if cookie.Name == "session" {
			responseData.sessionID = cookie.Value
			break
		}
	}

	err = json.Unmarshal(body, responseData.res)
	if err != nil {
		return responseData, fmt.Errorf("error Unmarshal resp: %v", err)
	}

	if !responseData.res.Success {
		return responseData, fmt.Errorf("%s", responseData.res.Message)
	}

	return responseData, nil
}

func (c *OneApiClient) randomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[c.seededRand.Intn(len(charset))]
	}
	return string(b)
}
