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
	URL        string
	Log        logger.Logger
	ClientHTTP *http.Client
}

// CheckUser checks if the user exists in the collector service
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

// RegisterUser reg the user in collection service
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

// GetCollections gets list of collections for user
// Need JWT token for this opperation
// Authorization: Bearer TOKEN
func (c *HTTPCollectorClient) GetUserWithCollections(reqData *collector.GetCollectionsRequest) (*collector.GetCollectionsResponse, error) {
	c.Log.Info("Get user's list of collections", logger.String("token", reqData.Token))

	request, err := http.NewRequest(http.MethodGet, c.URL+"collections", nil)
	if err != nil {
		c.Log.Error("Failed to create request", logger.Error(err))
		return nil, err
	}

	request.Header.Set("Authorization", "Bearer "+reqData.Token)

	resp, err := c.ClientHTTP.Do(request)
	if err != nil {
		c.Log.Error("Failed to do a request", logger.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	var respData collector.GetCollectionsResponse
	err = json.NewDecoder(resp.Body).Decode(&reqData)
	if err != nil {
		c.Log.Error("Failed to decode a body request", logger.Error(err))
		return nil, err
	}

	return &respData, nil
}

// Need JWT token for this opperation
func (c *HTTPCollectorClient) CreateCollection(reqData *collector.CreateCollectionRequest) (*collector.CreateCollectionResponse, error) {
	// doesn't look good to logging this information!!!
	c.Log.Info("Delete collection", logger.String("token_auth", reqData.Token))

	var collectionName = struct {
		Name string `json:"name"`
	}{
		Name: reqData.CollectionName,
	}
	body, err := json.Marshal(&collectionName)
	if err != nil {
		c.Log.Error("Failed to marshal data", logger.Error(err))
		return nil, err
	}

	request, err := http.NewRequest(http.MethodPost, c.URL+"/collections", bytes.NewBuffer(body))
	if err != nil {
		c.Log.Error("Failed to create request", logger.Error(err))
		return nil, err
	}

	request.Header.Set("Authorization", "Bearer "+reqData.Token)

	resp, err := c.ClientHTTP.Do(request)
	if err != nil {
		c.Log.Error("Failed to send request to collector service", logger.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		c.Log.Error("Failed to create collection's for user", logger.Int("status_code", resp.StatusCode))
		return nil, fmt.Errorf("failed to create collection's, status code: %d", resp.StatusCode)
	}

	var respData collector.CreateCollectionResponse
	err = json.NewDecoder(resp.Body).Decode(&respData)
	if err != nil {
		c.Log.Error("Failed to decode a responser body", logger.Error(err))
		return nil, err
	}

	return &respData, nil
}

// Need JWT token for this opperation
func (c *HTTPCollectorClient) RenameCollection(reqData *collector.RenameCollectionRequest) error {
	c.Log.Info("Rename collection", logger.String("collection_id", reqData.CollectionID))

	request, err := http.NewRequest(http.MethodPatch, c.URL+"/collections/"+reqData.CollectionID, nil)
	if err != nil {
		c.Log.Error("Failed to create request", logger.Error(err))
		return err
	}

	request.Header.Set("Authorization", "Bearer "+reqData.Token)

	resp, err := c.ClientHTTP.Do(request)
	if err != nil {
		c.Log.Error("Failed to do a request", logger.Error(err))
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		c.Log.Error("Failed to rename collection", logger.Int("status_code", resp.StatusCode))
		return fmt.Errorf("failed to rename the user's collection, status code: %d", resp.StatusCode)
	}

	return nil
}

// Need JWT token for this opperation
func (c *HTTPCollectorClient) DeleteCollection(reqData *collector.DeleteCollectionRequest) error {
	// doesn't look good to logging this information!!!
	c.Log.Info("Delete collection", logger.String("token_auth", reqData.Token))

	request, err := http.NewRequest(http.MethodDelete, c.URL+"/collections/"+reqData.CollectionID, nil)
	if err != nil {
		c.Log.Error("Failed to prepare a request", logger.Error(err))
		return err
	}

	request.Header.Set("Authorization", "Bearer "+reqData.Token)

	resp, err := c.ClientHTTP.Do(request)
	if err != nil {
		c.Log.Error("Failed to do a request", logger.Error(err))
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		c.Log.Error("Failed to delete the user's collection", logger.Int("status_code", resp.StatusCode))
		return fmt.Errorf("failed to delete the user's collection, status code: %d", resp.StatusCode)
	}

	return nil
}
