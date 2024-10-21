package OneApiSdk

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestCreateOpenAIUser(t *testing.T) {
	clientOptions := &OneApiOptions{
		ServerUrl:   os.Getenv("ONE_API_URL"),
		HttpTimeout: time.Second * 60,
		SystemToken: os.Getenv("ONE_API_SYSTEM_TOKEN"),
	}

	client := NewOneApiClient(clientOptions)
	data, err := client.CreateOpenAIUser(context.TODO(), "test100")
	if err != nil {
		t.Errorf("CreateOpenAIUser err: %v", err)
	}

	t.Log("success: ", data.Id)
}

func TestAddUserQuota(t *testing.T) {
	clientOptions := &OneApiOptions{
		ServerUrl:   os.Getenv("ONE_API_URL"),
		HttpTimeout: time.Second * 60,
		SystemToken: os.Getenv("ONE_API_SYSTEM_TOKEN"),
	}
	client := NewOneApiClient(clientOptions)

	if err := client.AddUserQuota(context.TODO(), 23, 100000); err != nil {
		t.Errorf("AddUserQuota err: %v", err)
	}

	t.Log("success: ")
}
