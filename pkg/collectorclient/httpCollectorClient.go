package collectorclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ShenokZlob/collector-ouphe/pkg/authctx"
	"github.com/ShenokZlob/collector-ouphe/pkg/contracts/auth"
	"github.com/ShenokZlob/collector-ouphe/pkg/contracts/collections"
	"github.com/ShenokZlob/collector-ouphe/pkg/logger"
)

type HTTPCollectorClient struct {
	URL        string
	Log        logger.Logger
	ClientHTTP *http.Client
}

func NewHTTPCollectorClient(url string, log logger.Logger) *HTTPCollectorClient {
	return &HTTPCollectorClient{
		URL:        url,
		Log:        log,
		ClientHTTP: http.DefaultClient,
	}
}

// CheckUser checks if the user exists in the collector service
func (c *HTTPCollectorClient) CheckUser(reqData *auth.CheckUserRequest) (*auth.CheckUserResponse, error) {
	c.Log.Info("Checking user in collector service", logger.String("method", "HTTPCollectorClient.CheckUser"), logger.Int("telegram_id", int(reqData.TelegramID)))

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

	var respData auth.CheckUserResponse
	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		c.Log.Error("Failed to decode response data", logger.Error(err))
		return nil, err
	}

	return &respData, nil
}

// RegisterUser reg the user in collection service
func (c *HTTPCollectorClient) RegisterUser(reqdata *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	c.Log.Info("Registering user in collector service", logger.String("method", "HTTPCollectorClient.RegisterUser"), logger.Int("telegram_id", int(reqdata.TelegramID)))

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

	if resp.StatusCode != http.StatusCreated {
		c.Log.Error("Failed to register user in collector service", logger.Int("status_code", resp.StatusCode))
		return nil, fmt.Errorf("failed to register user in collector service, status code: %d", resp.StatusCode)
	}

	var respData auth.RegisterResponse
	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		c.Log.Error("Failed to decode response data", logger.Error(err))
		return nil, err
	}

	return &respData, nil
}

// GetCollections gets list of collections for user
// Need JWT token for this opperation
// Authorization: Bearer TOKEN
func (c *HTTPCollectorClient) GetUserCollections(ctx context.Context) ([]collections.Collection, error) {
	token, ok := authctx.GetJWT(ctx)
	if !ok || token == "" {
		c.Log.Error("Authorization token is missing")
		return nil, fmt.Errorf("authorization token is missing")
	}

	c.Log.Info("Get user's list of collections", logger.String("method", "HTTPCollectorClient.GetUserCollections"), logger.String("token_auth", token))

	request, err := http.NewRequest(http.MethodGet, c.URL+"/collections", nil)
	if err != nil {
		c.Log.Error("Failed to create request", logger.Error(err))
		return nil, err
	}

	request.Header.Set("Authorization", "Bearer "+token)
	c.Log.Debug("Req token", logger.String("header", request.Header.Get("Authorization")))

	resp, err := c.ClientHTTP.Do(request)
	if err != nil {
		c.Log.Error("Failed to do a request", logger.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errorRepsonse collections.ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errorRepsonse); err != nil {
			c.Log.Error("Failed to decode error response", logger.Error(err))
			return nil, fmt.Errorf("failed to decode error response, status code: %d", resp.StatusCode)
		}
		c.Log.Error("Failed to get user's collections", logger.String("message", errorRepsonse.Message))
		return nil, fmt.Errorf("failed to get user's collections, status code: %d", resp.StatusCode)
	}

	var collections []collections.Collection
	err = json.NewDecoder(resp.Body).Decode(&collections)
	if err != nil {
		c.Log.Error("Failed to decode a body request", logger.Error(err))
		return nil, err
	}

	return collections, nil
}

// Need JWT token for this opperation
func (c *HTTPCollectorClient) CreateCollection(ctx context.Context, req *collections.CreateCollectionRequest) (*collections.Collection, error) {
	token, ok := authctx.GetJWT(ctx)
	if !ok || token == "" {
		c.Log.Error("Authorization token is missing")
		return nil, fmt.Errorf("authorization token is missing")
	}

	// doesn't look good to logging this information!!!
	c.Log.Info("Create collection", logger.String("metod", "HTTPCollectorClient.CreateCollection"), logger.String("token_auth", token))

	body, err := json.Marshal(&req)
	if err != nil {
		c.Log.Error("Failed to marshal data", logger.Error(err))
		return nil, err
	}

	request, err := http.NewRequest(http.MethodPost, c.URL+"/collections", bytes.NewBuffer(body))
	if err != nil {
		c.Log.Error("Failed to create request", logger.Error(err))
		return nil, err
	}

	request.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.ClientHTTP.Do(request)
	if err != nil {
		c.Log.Error("Failed to send request to collector service", logger.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		var errorResponse collections.ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			c.Log.Error("Failed to decode error response", logger.Error(err))
			return nil, fmt.Errorf("failed to decode error response, status code: %d", resp.StatusCode)
		}
		c.Log.Error("Failed to create collection's for user", logger.String("message", errorResponse.Message))
		return nil, fmt.Errorf("failed to create collection's, status code: %d", resp.StatusCode)
	}

	var collection collections.Collection
	err = json.NewDecoder(resp.Body).Decode(&collection)
	if err != nil {
		c.Log.Error("Failed to decode a responser body", logger.Error(err))
		return nil, err
	}

	return &collection, nil
}

// Need JWT token for this opperation
func (c *HTTPCollectorClient) RenameCollection(ctx context.Context, collectionID string, req *collections.RenameCollectionRequest) error {
	token, ok := authctx.GetJWT(ctx)
	if !ok || token == "" {
		c.Log.Error("Authorization token is missing")
		return fmt.Errorf("authorization token is missing")
	}

	c.Log.Info("Rename collection", logger.String("collection_id", collectionID), logger.String("method", "HTTPCollectorClient.RenameCollection"), logger.String("token_auth", token))

	body, err := json.Marshal(req)
	if err != nil {
		c.Log.Error("Failed to marshal request data")
		return err
	}

	request, err := http.NewRequest(http.MethodPatch, c.URL+"/collections/"+collectionID, bytes.NewBuffer(body))
	if err != nil {
		c.Log.Error("Failed to create request", logger.Error(err))
		return err
	}

	request.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.ClientHTTP.Do(request)
	if err != nil {
		c.Log.Error("Failed to do a request", logger.Error(err))
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		var errorResponse collections.ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			c.Log.Error("Failed to decode error response", logger.Error(err))
			return fmt.Errorf("failed to decode error response, status code: %d", resp.StatusCode)
		}
		c.Log.Error("Failed to rename collection", logger.String("message", errorResponse.Message))
		return fmt.Errorf("failed to rename the user's collection, status code: %d", resp.StatusCode)
	}

	return nil
}

// Need JWT token for this opperation
func (c *HTTPCollectorClient) DeleteCollection(ctx context.Context, collectionID string) error {
	token, ok := authctx.GetJWT(ctx)
	if !ok || token == "" {
		c.Log.Error("Authorization token is missing")
		return fmt.Errorf("authorization token is missing")
	}

	// doesn't look good to logging this information!!!
	c.Log.Info("Delete collection", logger.String("method", "HTTPCollectorClient.DeleteCollection"), logger.String("token_auth", token))

	request, err := http.NewRequest(http.MethodDelete, c.URL+"/collections/"+collectionID, nil)
	if err != nil {
		c.Log.Error("Failed to prepare a request", logger.Error(err))
		return err
	}

	request.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.ClientHTTP.Do(request)
	if err != nil {
		c.Log.Error("Failed to do a request", logger.Error(err))
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		var errorResponse collections.ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			c.Log.Error("Failed to decode error response", logger.Error(err))
			return fmt.Errorf("failed to decode error response, status code: %d", resp.StatusCode)
		}
		c.Log.Error("Failed to delete the user's collection", logger.String("message", errorResponse.Message))
		return fmt.Errorf("failed to delete the user's collection, status code: %d", resp.StatusCode)
	}

	return nil
}

func (c *HTTPCollectorClient) GetUsersCollectionByName(ctx context.Context, collectionName string) (*collections.Collection, error) {
	token, ok := authctx.GetJWT(ctx)
	if !ok || token == "" {
		c.Log.Error("Authorization token is missing")
		return nil, fmt.Errorf("authorization token is missing")
	}

	c.Log.Info("Get user's collection by name", logger.String("method", "HTTPCollectorClient.GetUsersCollectionByName"), logger.String("token_auth", token), logger.String("collection_name", collectionName))

	request, err := http.NewRequest(http.MethodGet, c.URL+"/collections/name/"+collectionName, nil)
	if err != nil {
		c.Log.Error("Failed to create request", logger.Error(err))
		return nil, err
	}

	request.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.ClientHTTP.Do(request)
	if err != nil {
		c.Log.Error("Failed to do a request", logger.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errorResponse collections.ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			c.Log.Error("Failed to decode error response", logger.Error(err))
			return nil, fmt.Errorf("failed to decode error response, status code: %d", resp.StatusCode)
		}
		c.Log.Error("Failed to get user's collection by name", logger.String("message", errorResponse.Message))
		return nil, fmt.Errorf("failed to get user's collection by name, status code: %d", resp.StatusCode)
	}

	var collection collections.Collection
	err = json.NewDecoder(resp.Body).Decode(&collection)
	if err != nil {
		c.Log.Error("Failed to decode a body request", logger.Error(err))
		return nil, err
	}

	return &collection, nil
}
