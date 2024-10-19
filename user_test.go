package OneAPISDK

import (
	"testing"
	"time"
)

func TestCreateOpenAIUser(t *testing.T) {
	clientOptions := &OneApiOptions{
		ServerUrl:   "",
		HttpTimeout: time.Second * 60,
		SystemToken: "",
	}

	client := NewOneApiClient(clientOptions)
	data, err := client.CreateOpenAIUser("", 0)
	if err != nil {
		t.Errorf("CreateOpenAIUser err: %v", err)
	}

	t.Log("success: ", data.BaseInfo.Id, data.KeyData[0].Key)
}

func TestAddUserQuota(t *testing.T) {
	clientOptions := &OneApiOptions{
		ServerUrl:   "",
		HttpTimeout: time.Second * 60,
		SystemToken: "",
	}
	client := NewOneApiClient(clientOptions)

	if err := client.AddUserQuota(23, 100000); err != nil {
		t.Errorf("AddUserQuota err: %v", err)
	}

	t.Log("success: ")
}
