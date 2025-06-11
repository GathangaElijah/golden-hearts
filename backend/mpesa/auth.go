package mpesa

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"time"
)

var (
	consumerKey    = os.Getenv("MPESA_CONSUMER_KEY")
	consumerSecret = os.Getenv("MPESA_CONSUMER_SECRET")
	baseURL        = "https://sandbox.safaricom.co.ke"
	tokenCache     string
	tokenExpiry    time.Time
)

type accessTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   string `json:"expires_in"`
}

// GetAccessToken fetches or returns cached token
func GetAccessToken() (string, error) {
	if time.Now().Before(tokenExpiry) && tokenCache != "" {
		return tokenCache, nil
	}

	url := baseURL + "/oauth/v1/generate?grant_type=client_credentials"
	req, _ := http.NewRequest("GET", url, nil)

	credentials := base64.StdEncoding.EncodeToString([]byte(consumerKey + ":" + consumerSecret))
	req.Header.Add("Authorization", "Basic "+credentials)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", errors.New("failed to get access token")
	}

	var tokenResp accessTokenResponse
	err = json.NewDecoder(resp.Body).Decode(&tokenResp)
	if err != nil {
		return "", err
	}

	tokenCache = tokenResp.AccessToken
	tokenExpiry = time.Now().Add(50 * time.Minute)

	return tokenCache, nil
}
