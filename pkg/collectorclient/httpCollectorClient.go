package collectorclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ShenokZlob/collector-ouphe/pkg/contracts/collector"
	"github.com/ShenokZlob/collector-ouphe/pkg/logger"
)

type HTTPCollectorClient struct {
	URL string
	Log logger.Logger
}

func (c *HTTPCollectorClient) CheckUser(reqData *collector.CheckUserRequest) (*collector.CheckUserResponse, error) {
	c.Log.Info("Checking user in collector service", logger.Int("telegram_id", int(reqData.TelegramID)))

	body, err := json.Marshal(reqData)
	if err != nil {
		c.Log.Error("Failed to marshal request data", logger.Error(err))
		return nil, err
	}

	resp, err := http.Post(c.URL+"/login", "application/json", bytes.NewBuffer(body))
	if err != nil {
		c.Log.Error("Failed to send request to collector service", logger.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.Log.Error("Failed to check user in collector service", logger.Int("status_code", resp.StatusCode))
		return nil, fmt.Errorf("failed to check user in collector service, status code: %d", resp.StatusCode)
	}

	var respData collector.CheckUserResponse
	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		c.Log.Error("Failed to decode response data", logger.Error(err))
		return nil, err
	}

	return &respData, nil
}

func (c *HTTPCollectorClient) RegisterUser(reqdata *collector.RegisterRequest) (*collector.RegisterResponse, error) {
	c.Log.Info("Registering user in collector service", logger.Int("telegram_id", int(reqdata.TelegramID)))

	body, err := json.Marshal(reqdata)
	if err != nil {
		c.Log.Error("Failed to marshal request data", logger.Error(err))
		return nil, err
	}

	resp, err := http.Post(c.URL+"/register", "application/json", bytes.NewBuffer(body))
	if err != nil {
		c.Log.Error("Failed to send request to collector service", logger.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.Log.Error("Failed to register user in collector service", logger.Int("status_code", resp.StatusCode))
		return nil, fmt.Errorf("failed to register user in collector service, status code: %d", resp.StatusCode)
	}

	var respData collector.RegisterResponse
	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		c.Log.Error("Failed to decode response data", logger.Error(err))
		return nil, err
	}

	return &respData, nil
}
