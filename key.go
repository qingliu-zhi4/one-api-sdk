package OneAPISDK

import (
	"fmt"
	"net/http"
)

type AddTokenReq struct {
	Name               string `json:"name"`
	RemainQuota        int    `json:"remain_quota"`
	ExpiredTime        int    `json:"expired_time"`
	UnlimitedQuota     bool   `json:"unlimited_quota"`
	ModelLimitsEnabled bool   `json:"model_limits_enabled"`
	ModelLimits        string `json:"model_limits"`
	AllowIps           string `json:"allow_ips"`
	Group              string `json:"group"`
}

type TokenInfo struct {
	UserId         int    `json:"user_id"`
	Key            string `json:"key"`
	Status         int    `json:"status"`
	Name           string `json:"name"`
	CreatedTime    int64  `json:"created_time"`
	AccessedTime   int64  `json:"accessed_time"`
	ExpiredTime    int64  `json:"expired_time"`
	RemainQuota    int64  `json:"remain_quota"`
	UnlimitedQuota bool   `json:"unlimited_quota"`
	UsedQuota      int64  `json:"used_quota"`
	Models         string `json:"models"`
	Subnet         string `json:"subnet"`
}

type OpenAIKeyData struct {
	Key    string
	Models string
}

func (c *OneApiClient) GenerateOpenAPIKey(accessToken string, req *AddTokenReq) error {
	_, err := c.sendReq(keyEndpoint, http.MethodPost, accessToken, "", req, nil)
	if err != nil {
		return fmt.Errorf("generate ppen API key fail: %v", err)
	}
	return nil
}

func (c *OneApiClient) GetOpenAPIKey(accessToken string) ([]*TokenInfo, error) {
	data, err := c.sendReq(keyEndpoint, http.MethodGet, accessToken, "", nil, nil)
	if err != nil {
		return nil, fmt.Errorf("get open API key fail: %v", err)
	}
	returnData := data.res.Data.([]interface{})
	info := make([]*TokenInfo, len(returnData))

	for index, _ := range returnData {
		item := returnData[index].(map[string]interface{})
		info[index] = &TokenInfo{
			UserId:         int(item["user_id"].(float64)),
			Key:            item["key"].(string),
			Status:         int(item["status"].(float64)),
			Name:           item["name"].(string),
			CreatedTime:    int64(item["created_time"].(float64)),
			AccessedTime:   int64(item["accessed_time"].(float64)),
			ExpiredTime:    int64(item["expired_time"].(float64)),
			RemainQuota:    int64(item["remain_quota"].(float64)),
			UnlimitedQuota: item["unlimited_quota"].(bool),
			UsedQuota:      int64(item["used_quota"].(float64)),
		}

		if models, ok := item["models"]; ok {
			info[index].Models = models.(string)
		}

		if subnet, ok := item["subnet"]; ok {
			info[index].Subnet = subnet.(string)
		}
	}

	return info, nil
}
