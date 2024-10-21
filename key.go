package OneApiSdk

import (
	"context"
	"fmt"
	"net/http"
)

// AddTokenReq 添加调用OpenAPI的token的请求数据
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

// TokenInfo 查询Token的返回数据
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

// GenerateOpenAPIKey 根据用户的accessToken生成此用户的OpenAI key
func (c *OneApiClient) GenerateOpenAPIKey(ctx context.Context, accessToken string, req *AddTokenReq) error {
	_, err := c.sendReq(ctx, keyEndpoint, http.MethodPost, accessToken, "", req, nil)
	if err != nil {
		return fmt.Errorf("generate ppen API key fail: %v", err)
	}
	return nil
}

// GetOpenAPIKey 根据用户的accessToken生成获取此用户的key信息
func (c *OneApiClient) GetOpenAPIKey(ctx context.Context, accessToken string) ([]*TokenInfo, error) {
	data, err := c.sendReq(ctx, keyEndpoint, http.MethodGet, accessToken, "", nil, nil)
	if err != nil {
		return nil, fmt.Errorf("get open API key fail: %v", err)
	}

	info := formTokenData(data)
	return info, nil
}

// GenerateSpecificNameKey GenerateSpecificNameKey 便利接口
// 这里需要传入根据你的用户维度，唯一命名的Key
// 如果你自己管理key，有信心保证用户维度的Key唯一，可以传递check为false，提升效率
func (c *OneApiClient) GenerateSpecificNameKey(ctx context.Context, accessToken string, keyName string, check bool) (*TokenInfo, error) {
	if check {
		info, err := c.GetSpecificNameKey(ctx, accessToken, keyName)
		if err != nil {
			return nil, fmt.Errorf("check openAPI key fail: %v", err)
		}
		if info != nil && len(info) > 0 {
			return info[0], nil
		}
	}

	req := &AddTokenReq{
		Name: keyName,
	}

	if err := c.GenerateOpenAPIKey(ctx, accessToken, req); err != nil {
		return nil, fmt.Errorf("generate specific API key fail in create: %v", err)
	}

	searchData := map[string]string{"keyword": req.Name}
	tokenData, err := c.sendReq(ctx, searchTokenEndpoint, http.MethodGet, accessToken, "", searchData, nil)
	if err != nil {
		return nil, fmt.Errorf("ggenerate specific API key fail in get: %v", err)
	}

	info := formTokenData(tokenData)
	if info == nil {
		return nil, fmt.Errorf("ggenerate specific API key fail odd error: %v", tokenData)
	}

	return info[0], nil
}

// GetSpecificNameKey 根据OpenKey 名称查询Key
func (c *OneApiClient) GetSpecificNameKey(ctx context.Context, accessToken string, keyName string) ([]*TokenInfo, error) {
	searchData := map[string]string{"keyword": keyName}
	tokenData, err := c.sendReq(ctx, searchTokenEndpoint, http.MethodGet, accessToken, "", searchData, nil)
	if err != nil {
		return nil, fmt.Errorf("ggenerate specific API key fail in get: %v", err)
	}

	info := formTokenData(tokenData)
	if info == nil {
		return nil, fmt.Errorf("ggenerate specific API key fail odd error: %v", tokenData)
	}

	return info, nil
}

func formTokenData(data *CommonAPIRes) []*TokenInfo {
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

	return info
}
