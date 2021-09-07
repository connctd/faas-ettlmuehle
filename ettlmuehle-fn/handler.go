package function

import (
	"io/ioutil"
	"net/http"
	"os"

	handler "github.com/openfaas/templates-sdk/go-http"
	log "github.com/sirupsen/logrus"
)

// Handle is called
func Handle(req handler.Request) (handler.Response, error) {
	clientSecret, err := getSecret("ettlmuehle-client-secret")
	clientID := os.Getenv("ettlmuehle-client-id")
	baseURL := os.Getenv("base-url")
	gatewayURL := os.Getenv("gateway-url")

	// faas not configured properly
	if err != nil {
		log.WithError(err).Errorln("Missing ettlmuehle-client-secret env var")
		panic("FAAS execution failed due to missing ettlmuehle-client-secret")
	}

	// restore or retrieve access token
	client, err := NewAPIClient(string(baseURL), string(gatewayURL), string(clientID), string(clientSecret))
	if err != nil {
		log.WithError(err).Errorln("Failed to create api client")
		panic("Failed to create api client")
	}

	simulateSensors(client)

	return handler.Response{
		StatusCode: http.StatusOK,
	}, err
}

func getSecret(secretName string) (secretBytes []byte, err error) {
	// read from the openfaas secrets folder
	secretBytes, err = ioutil.ReadFile("/var/openfaas/secrets/" + secretName)
	if err != nil {
		// read from the original location for backwards compatibility with openfaas <= 0.8.2
		secretBytes, err = ioutil.ReadFile("/run/secrets/" + secretName)
	}

	return secretBytes, err
}
