package OneApiSdk

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const defaultHttpTimeout = time.Second * 60

// OneApiClient OneAPISDK客户端实例
type OneApiClient struct {
	serverUrl   string
	httpClient  *http.Client
	systemToken string
	seededRand  *rand.Rand
}

// OneApiOptions 创建OneAPISDK客户端的选项
type OneApiOptions struct {
	//OneAPI server的地址
	ServerUrl string
	//root用户的AccessToken
	SystemToken     string
	HttpTimeout     time.Duration
	TransportConfig http.RoundTripper
}

// NewOneApiClient 创建新的OneApiClient，ServerUrl和SystemToken是必选的，
// 这两项配置也可以通过环境变量ONE_API_HOST和ONE_API_SYSTEM_TOKEN传入
func NewOneApiClient(options ...*OneApiOptions) *OneApiClient {
	oneApiHost := os.Getenv("ONE_API_HOST")
	systemToken := os.Getenv("ONE_API_SYSTEM_TOKEN")
	randTripper := http.DefaultTransport
	httpTimeout := defaultHttpTimeout

	if len(options) == 1 {
		oneApiHost = options[0].ServerUrl
		systemToken = options[0].SystemToken
		if options[0].TransportConfig == nil {
			randTripper = options[0].TransportConfig
		}
		if options[0].HttpTimeout != 0 {
			httpTimeout = options[0].HttpTimeout
		}
	}

	if len(systemToken) == 0 || len(oneApiHost) == 0 {
		return nil
	}

	client := &OneApiClient{
		serverUrl: oneApiHost,
		httpClient: &http.Client{
			Transport: randTripper,
			Timeout:   httpTimeout,
		},
		systemToken: systemToken,
		seededRand:  rand.New(rand.NewSource(time.Now().UnixNano())),
	}
	return client
}

func (c *OneApiClient) sendReq(ctx context.Context, endpoint string, method string, token string, sessionID string, reqData interface{}, additionHeader map[string]string) (*CommonAPIRes, error) {
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

	req, err := http.NewRequestWithContext(ctx, method, u.String(), bytes.NewBuffer(jsonData))
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
		req.Header.Set("Content-Type", contentTypeJson)
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
