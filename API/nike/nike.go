package nike

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"os"
	"runsync/API"
	"sort"
	"time"
)

const (
	baseURL = "https://api.nike.com/"

	httpTimeout = 30 * time.Second
)

func GetRunsFromNRC() ([]activity, error) {
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
	activityIds, err := GetActivities(ctx, *accessToken)
	if err != nil {
		return nil, errors.WithMessagef(err, "Fail to get activities from Nike Run Club")
	}

	// Filter to keep runs activities only
	runs := []activity{}
	for _, activity := range activityIds {
		if activity.Type == "run" {
			log.Infof("[nike] Retrieve run details for [%v]", activity.ID)
			run, err := GetActivity(ctx, *accessToken, activity.ID)
			if err != nil {
				return nil, errors.WithMessagef(err, "Fail to get run from Nike Run Club for [%v]", activity.ID)
			}
			runs = append(runs, *run)
		} else {
			log.Infof("[nike] activity [%v] skipped cause it has type [%v]", activity.ID, activity.Type)
		}
	}

	// TODO Do we have to filter on longitude et latitude?

	return runs, nil
}

func BuildGpxFromActivity(activity activity) *API.GPX {
	unixStartTime := time.Unix(activity.StartEpoch/1000, activity.StartEpoch%1000).UTC()
	startTimeString := unixStartTime.Format(time.RFC3339Nano)

	var trackpoints []API.TrackPoint

	latitudes := findMetric(activity.Metrics, "latitude")
	longitudes := findMetric(activity.Metrics, "longitude")
	elevations := findMetric(activity.Metrics, "elevation")
	heartRates := findMetric(activity.Metrics, "heart_rate")

	if latitudes != nil && longitudes != nil {
		for i := 0; i < len(latitudes.Values); i++ {
			tp := API.TrackPoint{
				Latitude:  fmt.Sprintf("%v", latitudes.Values[i].Value),
				Longitude: fmt.Sprintf("%v", longitudes.Values[i].Value),
				Start:     latitudes.Values[i].Start,
				Time:      time.Unix(latitudes.Values[i].Start/1000, latitudes.Values[i].Start%1000).UTC().Format(time.RFC3339Nano),
			}
			trackpoints = append(trackpoints, tp)
		}
	}

	if elevations != nil {
		var index = 0
		for i := 0; i < len(trackpoints); i++ {
			point := trackpoints[i]
			if elevations.Values[index].Start < point.Start && index < (len(elevations.Values)-1) {
				index++
			}
			trackpoints[i].Elevation = fmt.Sprintf("%v", elevations.Values[index].Value)
		}
	}

	if heartRates != nil {
		var index = 0
		for i := 0; i < len(trackpoints); i++ {
			point := trackpoints[i]
			if heartRates.Values[index].Start < point.Start && index < (len(heartRates.Values)-1) {
				index++
			}
			trackpoints[i].Extensions = []API.Extensions{
				{
					TrackPointExtensions: []API.TrackPointExtension{
						{
							HeartRate: int(heartRates.Values[index].Value),
						},
					},
				},
			}
		}
	}

	return &API.GPX{
		Creator:        "StravaGPX",
		XmlnsXsi:       "http://www.w3.org/2001/XMLSchema-instance",
		SchemaLocation: "http://www.topografix.com/GPX/1/1 http://www.topografix.com/GPX/1/1/gpx.xsd http://www.garmin.com/xmlschemas/GpxExtensions/v3 http://www.garmin.com/xmlschemas/GpxExtensionsv3.xsd http://www.garmin.com/xmlschemas/TrackPointExtension/v1 http://www.garmin.com/xmlschemas/TrackPointExtensionv1.xsd",
		Version:        "1.1",
		Xmlns:          "http://www.topografix.com/GPX/1/1",
		XmlnsGpxtpx:    "http://www.garmin.com/xmlschemas/TrackPointExtension/v1",
		XmlnsGpxx:      "http://www.garmin.com/xmlschemas/GpxExtensions/v3",

		Metadata: API.Metadata{
			Time: startTimeString,
		},

		Track: API.Track{
			Name: fmt.Sprintf("%v run - NRC", unixStartTime.Weekday()),
			Type: 9,
			TrackSegment: API.TrackSegment{
				TrackPoints: trackpoints,
			},
		},
	}
}

func BuildTcxFromActivity(activity activity) *API.TrainingCenterDatabase {
	unixStartTime := time.Unix(activity.StartEpoch/1000, activity.StartEpoch%1000).UTC().Format(time.RFC3339)

	distance := findSummary(activity.Summaries, "distance")
	calories := findSummary(activity.Summaries, "calories")
	heartRate := findSummary(activity.Summaries, "heart_rate")
	speedMean := findSummary(activity.Summaries, "speed")

	speedValues := findMetric(activity.Metrics, "speed").Values
	sort.Slice(speedValues, func(i, j int) bool {
		return speedValues[i].Value > speedValues[j].Value
	})

	heartRatesValues := findMetric(activity.Metrics, "heart_rate").Values
	sort.Slice(heartRatesValues, func(i, j int) bool {
		return heartRatesValues[i].Value > heartRatesValues[j].Value
	})

	trackpoints := []API.TcxTrackpoint{}
	speeds := findMetric(activity.Metrics, "speed").Values
	distances := findMetric(activity.Metrics, "distance").Values
	heartRates := findMetric(activity.Metrics, "heart_rate").Values

	sort.Slice(speeds, func(i, j int) bool {
		return speeds[i].Start < speeds[j].Start
	})
	for _, speed := range speeds {
		tp := API.TcxTrackpoint{
			Time:           time.Unix(speed.Start/1000, speed.Start%1000).UTC().Format(time.RFC3339),
			DistanceMeters: 0,
			//HeartRateBpm:   &API.Value{},
			Extensions: API.TrackExtension{
				TPX: API.TPX{
					Xmlns: "http://www.garmin.com/xmlschemas/ActivityExtension/v2",
					Speed: float32(speed.Value) * 0.277778,
				},
			},
		}
		trackpoints = append(trackpoints, tp)
	}

	for i := range trackpoints {
		if i >= len(distances) {
			break
		}
		d := float32(distances[i].Value) * 1000
		if i > 0 {
			d = d + trackpoints[i-1].DistanceMeters
		}
		trackpoints[i].DistanceMeters = d
	}

	for i := 0; i < len(heartRates)-1; i++ {
		for j := range trackpoints {
			hrt, err := time.Parse(time.RFC3339, trackpoints[j].Time)
			if err != nil {
				fmt.Println(err)
			}
			ts := hrt.Unix() * 1000

			if ts >= heartRates[i].Start {
				trackpoints[j].HeartRateBpm = &API.Value{
					Value: int32(heartRates[i].Value),
				}
				break
			}

		}
	}

	tcxActivities := []API.Activity{}
	tcxActivity := API.Activity{
		Sport: "Running",
		ID:    unixStartTime,
		Lap: API.Lap{
			StartTime:        unixStartTime,
			TotalTimeSeconds: float32(activity.ActivityDuration) / 1000.0,
			DistanceMeters:   distance.Value * 1000.0,
			MaximumSpeed:     float32(speedValues[0].Value) * 0.277778,
			Calories:         int32(calories.Value),
			AverageHeartRateBpm: API.Value{
				Value: int32(heartRate.Value),
			},
			MaximumHeartRateBpm: API.Value{
				Value: int32(heartRatesValues[0].Value),
			},
			Intensity:     "Active",
			TriggerMethod: "Manual",
			Track: API.TcxTrack{
				Trackpoint: trackpoints,
			},
			Extensions: API.LapExtensions{
				LX: API.LX{
					AvgSpeed: speedMean.Value * 0.277778,
				},
			},
		},
	}
	tcxActivities = append(tcxActivities, tcxActivity)

	return &API.TrainingCenterDatabase{
		SchemaLocation: "http://www.garmin.com/xmlschemas/TrainingCenterDatabase/v2 http://www.garmin.com/xmlschemas/TrainingCenterDatabasev2.xsd",
		Xmlns:          "http://www.garmin.com/xmlschemas/TrainingCenterDatabase/v2",
		XmlnsXsi:       "http://www.w3.org/2001/XMLSchema-instance",
		XmlnsNs2:       "http://www.garmin.com/xmlschemas/UserProfile/v2",
		XmlnsNs3:       "http://www.garmin.com/xmlschemas/ActivityExtension/v2",
		XmlnsNs4:       "http://www.garmin.com/xmlschemas/ProfileExtension/v1",
		XmlnsNs5:       "http://www.garmin.com/xmlschemas/ActivityGoals/v1",

		Activities: API.Activities{
			Activities: tcxActivities,
		},

		Author: API.Author{
			Type: "Application_t",
			Name: "Aurélien Allienne",
			Build: API.Build{
				Version: API.Version{
					VersionMajor: 1,
					VersionMinor: 0,
					BuildMajor:   1,
					BuildMinor:   0,
				},
			},
			LangID: "en",
		},
	}
}

func findMetric(metrics []metric, t string) *metric {
	for _, n := range metrics {
		if t == n.Type {
			return &n
		}
	}
	return nil
}
func findSummary(metrics []summary, t string) *summary {
	for _, n := range metrics {
		if t == n.Metric {
			return &n
		}
	}
	return nil
}
