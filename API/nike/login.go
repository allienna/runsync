package nike

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"net/http"
	"runsync/API"
)

const(
	getTokenEndpoint = "idn/shim/oauth/2.0/token"
)

type loginRequest struct {
	ClientID     string `json:"client_id"`
	GrantType    string `json:"grant_type"`
	UxID         string `json:"ux_id"`
	RefreshToken string `json:"refresh_token"`
}

type loginResponse struct {
	AccessToken string `json:"access_token"`
}

func GetBearer(ctx context.Context, clientID, refreshToken string) (*string, error) {
	ctx, cancel := context.WithTimeout(ctx, httpTimeout)
	defer cancel()

	b, err := json.Marshal(loginRequest{
		ClientID:     clientID,
		GrantType:    "refresh_token",
		UxID:         "com.nike.sport.running.ios.6.5.1",
		RefreshToken: refreshToken,
	})
	if err != nil {
		return nil, err
	}
	body := bytes.NewReader(b)

	request, err := http.NewRequest(http.MethodPost, baseURL+getTokenEndpoint, body)
	if err != nil {
		return nil, err
	}

	header := request.Header
	header.Set("Content-Type", "application/json")

	response, err := API.GetClient().Do(request)
	if err != nil {
		return nil, errors.WithMessage(err, "Failed to connect to Nike API")
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

	//glog.Infof("[nike] Bearer: %v ", data.AccessToken)
	return &data.AccessToken, nil
}