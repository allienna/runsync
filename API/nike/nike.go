package nike

import (
	"context"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

const (
	baseURL = "https://api.nike.com/"

	httpTimeout = 30 * time.Second
)

func GetActivitiesFromNRC() ([]activity, error) {
	clientId := os.Getenv("NIKE_CLIENT_ID")
	refreshToken := os.Getenv("NIKE_REFRESH_TOKEN")

	if len(clientId) == 0 || len(refreshToken) == 0 {
		return nil, errors.New("Please set your Nike Run Club application parameters in .env")
	}

	ctx := context.Background()
	accessToken, err := GetBearer(ctx, clientId, refreshToken)
	if err != nil {
		return nil, errors.WithMessage(err, "Fail to get bearer from Nike Run Club")
	}


	log.Infof("[nike] Bearer retrieved with success")

	// Fetch from Nike API
	activityIds, err := GetActivityIds(ctx, *accessToken)
	if err != nil {
		return nil, errors.WithMessagef(err, "Fail to get activities from Nike Run Club")
	}

	return activityIds, nil
}
