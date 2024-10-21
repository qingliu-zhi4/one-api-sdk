package OneApiSdk

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
)

// OneApiUserData getUser时的返回
type OneApiUserData struct {
	Id               int    `json:"id"`
	Username         string `json:"username"`
	Password         string `json:"password"`
	DisplayName      string `json:"display_name"`
	Role             int    `json:"role"`
	Status           int    `json:"status"`
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

// CommonUserReq 创建用户的入参
type CommonUserReq struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	DisplayName string `json:"display_name"`
}

// BaseUserinfo 便利的创建用户方法，返回的基础用户数据
type BaseUserinfo struct {
	Id          int    `json:"id"`
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	Role        int    `json:"role"`
	Status      int    `json:"status"`
	AccessToken string `json:"access_token"`
}

// AddUser 创建普通用户
func (c *OneApiClient) AddUser(ctx context.Context, req *CommonUserReq) error {
	_, err := c.sendReq(ctx, userEndpoint, http.MethodPost, c.systemToken, "", req, nil)
	if err != nil {
		return fmt.Errorf("add user call error %v", err)
	}
	return nil
}

// CreateOpenAIUser 便利方法，创建用户，同时创建这个用户的AccessToken
// 返回自动化运营常用的基础数据
// 请保证UserName的唯一性，这样便利性方法才保证正确, UserName长度最大为12
// 如果你没有信心保证UserName唯一，可以传空，后期利用UserID管理用户
func (c *OneApiClient) CreateOpenAIUser(ctx context.Context, userName string) (*BaseUserinfo, error) {
	uniqueName := c.randomString(12)
	if len(userName) > 0 {
		uniqueName = userName
	}

	req := &CommonUserReq{
		Username: uniqueName,
		Password: c.randomString(20),
	}
	if err := c.AddUser(ctx, req); err != nil {
		return nil, fmt.Errorf("add user call error %v", err)
	}

	loginData, err := c.UserLogin(ctx, req)
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

	token, err := c.GetToken(ctx, strconv.Itoa(userData.Id), loginData.sessionID)
	if err != nil {
		return nil, fmt.Errorf("get token with sessionID error %v", err)
	}

	userData.AccessToken = token

	return userData, nil
}

// GetUser 获取用户信息
func (c *OneApiClient) GetUser(ctx context.Context, userID int) (*OneApiUserData, error) {
	data, err := c.sendReq(ctx, fmt.Sprintf("%s/%d", userEndpoint, userID), http.MethodGet, c.systemToken, "", nil, nil)
	if err != nil {
		return nil, fmt.Errorf("get user call error %v", err)
	}

	userData := formUserData(data)
	return userData, nil
}

// UserLogin 用户登录，这个方法本质是为了获取用户的sessionID
func (c *OneApiClient) UserLogin(ctx context.Context, req *CommonUserReq) (*CommonAPIRes, error) {
	data, err := c.sendReq(ctx, loginEndpoint, http.MethodPost, "", "", req, nil)
	if err != nil {
		return nil, fmt.Errorf("add user call error %v", err)
	}
	return data, nil
}

// GetToken 获取用户的AccessToken，这个方法需要用户的sessionID
func (c *OneApiClient) GetToken(ctx context.Context, userID string, sessionID string) (string, error) {
	newApiUserMap := map[string]string{"new-api-user": userID}
	data, err := c.sendReq(ctx, tokenEndpoint, http.MethodGet, "", sessionID, nil, newApiUserMap)
	if err != nil {
		return "", fmt.Errorf("get token call error %v", err)
	}
	token := data.res.Data.(string)
	return token, nil
}

// UpdateUser 更新用户数据，注意不是所有的数据都可以通过这个接口更新
// 目前用户名，密码，额度，状态，角色可以通过此接口更新
func (c *OneApiClient) UpdateUser(ctx context.Context, req *OneApiUserData) error {
	_, err := c.sendReq(ctx, userEndpoint, http.MethodPut, c.systemToken, "", req, nil)
	if err != nil {
		return fmt.Errorf("update user info %v", err)
	}
	return nil
}

// AddUserQuota 给指定用户添加限额的便利方法
func (c *OneApiClient) AddUserQuota(ctx context.Context, userID int, quota int64) error {
	userData, err := c.GetUser(ctx, userID)
	if err != nil {
		return fmt.Errorf("get user call error %v", err)
	}
	userData.Quota += quota
	if err := c.UpdateUser(ctx, userData); err != nil {
		return fmt.Errorf("update user quota %v", err)
	}
	return nil
}

func formUserData(data *CommonAPIRes) *OneApiUserData {
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

	return userData
}
