package OneApiSdk

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestGenerateSpecificNameKey(t *testing.T) {
	clientOptions := &OneApiOptions{
		ServerUrl:   os.Getenv("ONE_API_URL"),
		HttpTimeout: time.Second * 60,
		SystemToken: os.Getenv("ONE_API_SYSTEM_TOKEN"),
	}
	client := NewOneApiClient(clientOptions)

	userData, err := client.GetUser(context.TODO(), 23)
	if err != nil {
		t.Errorf("GetUser err: %v", err)
	}

	keyData, err := client.GenerateSpecificNameKey(context.TODO(), userData.AccessToken, "testName", true)
	if err != nil {
		t.Errorf("GenerateSpecificNameKey err: %v", err)
	}

	t.Log("success: ", keyData.Key)
}
