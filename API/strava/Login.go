package strava

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"runsync/API"
	"strings"
)

type loginResponse struct {
	AccessToken string `json:"access_token"`
}

func GetBearer(ctx context.Context, clientID, clientSecret, refreshToken string) (*string, error) {
	ctx, cancel := context.WithTimeout(ctx, httpTimeout)
	defer cancel()

	body := url.Values{}
	body.Set("client_id", clientID)
	body.Set("client_secret", clientSecret)
	body.Set("grant_type", "refresh_token")
	body.Set("refresh_token", refreshToken)

	request, err := http.NewRequest(
		http.MethodPost,
		"https://www.strava.com/api/v3/oauth/token",
		strings.NewReader(body.Encode()))
	if err != nil {
		return nil, err
	}

	header := request.Header
	header.Set("Content-Type", "application/x-www-form-urlencoded")

	response, err := API.GetClient().Do(request)
	if err != nil {
		return nil, errors.WithMessage(err, "Failed to connect to Strava API")
	}

	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, errors.WithMessage(errors.New(response.Status), "Failed to login")
	}

	var data loginResponse
	decoder := json.NewDecoder(response.Body)
	if err = decoder.Decode(&data); err != nil {
		return nil, errors.WithMessage(err, API.ErrInvalidLoginResponse.Error())
	}

	log.Debugf("Oauth token: %v", data)
	return &data.AccessToken, nil
}
