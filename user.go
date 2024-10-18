package one_api_sdk

import (
	"fmt"
	"net/http"
	"strconv"
)

type OneApiUserData struct {
	Id               int    `json:"id"`
	Username         string `json:"username"`
	Password         string `json:"password"`
	DisplayName      string `json:"display_name"`
	Role             int    `json:"role"`   // admin, util
	Status           int    `json:"status"` // enabled, disabled
	Email            string `json:"email"`
	GitHubId         string `json:"github_id"`
	WeChatId         string `json:"wechat_id"`
	LarkId           string `json:"lark_id"`
	OidcId           string `json:"oidc_id"`
	VerificationCode string `json:"verification_code"`
	AccessToken      string `json:"access_token"`
	Quota            int64  `json:"quota"`
	UsedQuota        int64  `json:"used_quota"`
	RequestCount     int    `json:"request_count"`
	Group            string `json:"group"`
	AffCode          string `json:"aff_code"`
	InviterId        int    `json:"inviter_id"`
}

type OpenAIUserData struct {
	BaseInfo    *BaseUserinfo
	AccessToken string
	Quota       int64
	UsedQuota   int64
	KeyData     []*OpenAIKeyData
}

type CommonUserReq struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	DisplayName string `json:"display_name"`
}

type BaseUserinfo struct {
	Id          int    `json:"id"`
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	Role        int    `json:"role"`
	Status      int    `json:"status"`
}

func (c *OneApiClient) AddUser(req *CommonUserReq) error {
	_, err := c.sendReq(userEndpoint, http.MethodPost, c.systemToken, "", req, nil)
	if err != nil {
		return fmt.Errorf("add user call error %v", err)
	}
	return nil
}

func (c *OneApiClient) CreateOpenAIUser(openAIModels string, initQuota int) (*OpenAIUserData, error) {
	req := &CommonUserReq{
		Username: c.randomString(12),
		Password: c.randomString(20),
	}
	if err := c.AddUser(req); err != nil {
		return nil, fmt.Errorf("add user call error %v", err)
	}

	loginData, err := c.UserLogin(req)
	if err != nil {
		return nil, fmt.Errorf("login with add user call error %v", err)
	}

	returnData := loginData.res.Data.(map[string]interface{})
	userData := &BaseUserinfo{
		Id:          int(returnData["id"].(float64)),
		Username:    returnData["username"].(string),
		DisplayName: returnData["display_name"].(string),
		Role:        int(returnData["role"].(float64)),
		Status:      int(returnData["status"].(float64)),
	}

	token, err := c.GetToken(strconv.Itoa(userData.Id), loginData.sessionID)
	if err != nil {
		return nil, fmt.Errorf("get token with sessionID error %v", err)
	}

	openAIKeyReq := &AddTokenReq{
		Name:               c.randomString(12),
		ExpiredTime:        0,
		ModelLimitsEnabled: len(openAIModels) > 0,
		ModelLimits:        openAIModels,
	}

	if err := c.GenerateOpenAPIKey(token, openAIKeyReq); err != nil {
		return nil, fmt.Errorf("create openai key with token error %v", err)
	}

	fmt.Println("GenerateOpenAPIKey success ")

	openAIKeys, err := c.GetOpenAPIKey(token)
	if err != nil {
		return nil, fmt.Errorf("get openai key with token error %v", err)
	}

	keyData := &OpenAIKeyData{
		Key:    openAIKeys[0].Key,
		Models: openAIKeys[0].Models,
	}

	userBasicData := &OpenAIUserData{
		BaseInfo:    userData,
		Quota:       int64(initQuota),
		UsedQuota:   0,
		KeyData:     []*OpenAIKeyData{keyData},
		AccessToken: token,
	}
	return userBasicData, nil
}

func (c *OneApiClient) GetUser(userID int) (*OneApiUserData, error) {
	data, err := c.sendReq(fmt.Sprintf("%s/%d", userEndpoint, userID), http.MethodGet, c.systemToken, "", nil, nil)
	if err != nil {
		return nil, fmt.Errorf("get user call error %v", err)
	}

	returnData := data.res.Data.(map[string]interface{})
	userData := &OneApiUserData{
		Id:           int(returnData["id"].(float64)),
		Username:     returnData["username"].(string),
		Password:     returnData["password"].(string),
		DisplayName:  returnData["display_name"].(string),
		Role:         int(returnData["role"].(float64)),
		Status:       int(returnData["status"].(float64)),
		Email:        returnData["email"].(string),
		GitHubId:     returnData["github_id"].(string),
		WeChatId:     returnData["wechat_id"].(string),
		AccessToken:  returnData["access_token"].(string),
		Quota:        int64(returnData["quota"].(float64)),
		UsedQuota:    int64(returnData["used_quota"].(float64)),
		RequestCount: int(returnData["request_count"].(float64)),
		Group:        returnData["group"].(string),
		AffCode:      returnData["aff_code"].(string),
		InviterId:    int(returnData["inviter_id"].(float64)),
	}

	return userData, nil
}

func (c *OneApiClient) UserLogin(req *CommonUserReq) (*CommonAPIRes, error) {
	data, err := c.sendReq(loginEndpoint, http.MethodPost, "", "", req, nil)
	if err != nil {
		return nil, fmt.Errorf("add user call error %v", err)
	}
	return data, nil
}

func (c *OneApiClient) GetToken(userID string, sessionID string) (string, error) {
	newApiUserMap := map[string]string{"new-api-user": userID}
	data, err := c.sendReq(tokenEndpoint, http.MethodGet, "", sessionID, nil, newApiUserMap)
	if err != nil {
		return "", fmt.Errorf("get token call error %v", err)
	}
	token := data.res.Data.(string)
	return token, nil
}

func (c *OneApiClient) UpdateUser(req *OneApiUserData) error {
	_, err := c.sendReq(userEndpoint, http.MethodPut, c.systemToken, "", req, nil)
	if err != nil {
		return fmt.Errorf("update user info %v", err)
	}
	return nil
}

func (c *OneApiClient) AddUserQuota(userID int, quota int64) error {
	userData, err := c.GetUser(userID)
	if err != nil {
		return fmt.Errorf("get user call error %v", err)
	}
	userData.Quota += quota
	if err := c.UpdateUser(userData); err != nil {
		return fmt.Errorf("update user quota %v", err)
	}
	return nil
}
