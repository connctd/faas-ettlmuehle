package function

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"

	log "github.com/sirupsen/logrus"

	"net/http"
	"strconv"
)

var (
	tmpAccessTokenFile = "/tmp/token"
)

// APIClient allows sending requests against connctd
type APIClient interface {
	UpdateProperty(thingID, componentID, propertyID, newValue string)
	GetToken() Token
}

// Client defines an api client
type Client struct {
	clientID     string
	clientSecret string
	token        Token
	httpClient   *http.Client
	baseURL      string
	gatewayURL   string
}

// NewAPIClient creates new client. If no token could be retrieved or restored error is returned
func NewAPIClient(baseURL string, gatewayURL string, clientID string, clientSecret string) (APIClient, error) {
	client := Client{baseURL: baseURL, gatewayURL: gatewayURL, clientID: clientID, clientSecret: clientSecret, httpClient: http.DefaultClient}

	token, err := client.retrieveToken()
	if err != nil {
		return nil, err
	}

	client.token = token

	return &client, nil
}

// UpdateProperty implements inferface definition
func (c *Client) UpdateProperty(thingID, componentID, propertyID, newValue string) {
	header := make(map[string]string)

	header["Authorization"] = "Bearer " + c.token.AccessToken
	header["Content-Type"] = "application/json"
	header["X-External-Subject-ID"] = "default"

	payload := PropertyUpdate{Value: newValue}
	if err := c.doRequest(context.Background(), http.MethodPut, c.baseURL, "/api/v1/things/"+thingID+"/components/"+componentID+"/properties/"+propertyID, http.StatusNoContent, payload, nil, &header); err != nil {
		log.WithError(err).Warningln("Failed to update thing property")
	}
}

// GetToken implements inferface definition
func (c *Client) GetToken() Token {
	return c.token
}

func (c *Client) retrieveToken() (Token, error) {
	var token Token
	payload := TokenRequest{ClientID: c.clientID, ClientSecret: c.clientSecret}
	if err := c.doRequest(context.Background(), http.MethodPost, c.gatewayURL, "/function/authcc-fn", http.StatusOK, payload, &token, nil); err != nil {
		log.WithError(err).Warningln("Failed to retrieve token")
		return token, err
	}

	return token, nil
}

func (c *Client) doRequest(ctx context.Context, method string, baseURL string, endpoint string, expectedStatusCode int, payload interface{}, response interface{}, includedHeaders *map[string]string) error {
	var req *http.Request
	var err error

	if payload != nil {
		payloadBytes, perr := json.Marshal(payload)

		if perr != nil {
			return perr
		}

		req, err = http.NewRequest(method, baseURL+endpoint, bytes.NewBuffer(payloadBytes))
	} else {
		req, err = http.NewRequest(method, baseURL+endpoint, nil)
	}

	// propagate context so that premature cancelation can be done or timeouts realized
	req = req.WithContext(ctx)

	if err != nil {
		return err
	}

	if payload != nil {
		defer req.Body.Close()
	}

	// add headers
	if includedHeaders != nil {
		for key, value := range *includedHeaders {
			req.Header.Add(key, value)
		}
	}

	resp, err := c.httpClient.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != expectedStatusCode {
		return errors.New("Invalid status code: " + strconv.Itoa(resp.StatusCode))
	}

	if response != nil {
		err = json.NewDecoder(resp.Body).Decode(response)

		if err != nil {
			return err
		}
	}

	return nil
}

// PropertyUpdate describes a property update
type PropertyUpdate struct {
	Value string `json:"value"`
}

// TokenRequest describes a token request
type TokenRequest struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}
