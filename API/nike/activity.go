package nike

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	"runsync/API"
	"strconv"
	"time"
)

const (
	getActivitiesByTimeEndpoint = "sport/v3/me/activities/after_time/"
	getActivitiesByIdEndpoint   = "sport/v3/me/activity/%s?metrics=ALL"
)

type activities struct {
	Activities []activity `json:"activities"`
	Paging     paging     `json:"paging"`
}

type activity struct {
	ID          string   `json:"id"`
	Type        string   `json:"type"`
	StartEpoch  int64    `json:"start_epoch_ms"`
	MetricTypes []string `json:"metric_types"`
	Metrics     []metric `json:"metrics"`
}

type metric struct {
	Type   string        `json:"type"`
	Unit   string        `json:"unit"`
	Values []metricValue `json:"values"`
}

type metricValue struct {
	Start int64   `json:"start_epoch_ms"`
	End   int64   `json:"end_epoch_ms"`
	Value float64 `json:"value"`
}

type paging struct {
	AfterTime int64  `json:"after_time"`
	AfterID   string `json:"after_id"`
}

// Get activities for the last 14 days
func GetActivityIds(ctx context.Context, accessToken string) ([]activity, error) {
	tm := time.Now().AddDate(0, 0, -14).Unix()

	ctx, cancel := context.WithTimeout(ctx, httpTimeout)
	defer cancel()

	var activityIds []activity = make([]activity, 0)
	for ; tm > 0; {
		request, err := http.NewRequest(
			http.MethodGet,
			baseURL+getActivitiesByTimeEndpoint+strconv.FormatInt(tm, 10),
			nil)
		if err != nil {
			return nil, err
		}

		header := request.Header
		header.Set("Authorization", "Bearer "+accessToken)

		response, err := API.GetClient().Do(request)
		if err != nil {
			return nil, errors.WithMessage(err, "Fail to connect to Nike API")
		}

		defer response.Body.Close()
		if response.StatusCode != http.StatusOK {
			return nil, errors.WithMessage(errors.New(response.Status), "Fail to get activities")
		}

		var data activities
		decoder := json.NewDecoder(response.Body)
		if err = decoder.Decode(&data); err != nil {
			return nil, errors.WithMessage(err, API.ErrInvalidLoginResponse.Error())
		}

		activityIds = append(activityIds, data.Activities...)

		tm = data.Paging.AfterTime
	}
	return activityIds, nil
}

func GetActivity(ctx context.Context, accessToken string, activityId string) (*activity, error) {
	ctx, cancel := context.WithTimeout(ctx, httpTimeout)
	defer cancel()

	request, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf(baseURL+getActivitiesByIdEndpoint, activityId),
		nil)

	if err != nil {
		return nil, errors.WithMessage(err, "Fail to create a new request")
	}

	header := request.Header
	header.Set("Authorization", "Bearer "+accessToken)

	response, err := API.GetClient().Do(request)
	if err != nil {
		return nil, errors.WithMessage(err, "Fail to connect to Nike API")
	}

	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, errors.WithMessagef(errors.New(response.Status), "Fail to get activity with id %v", activityId)
	}

	var data activity
	decoder := json.NewDecoder(response.Body)
	if err = decoder.Decode(&data); err != nil {
		return nil, errors.WithMessage(err, API.ErrInvalidLoginResponse.Error())
	}

	return &data, nil
}
