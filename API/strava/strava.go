package strava

import (
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"runsync/API"
	"strings"
	"time"
)

const (
	httpTimeout = 30 * time.Second
)

func ImportDataFromFiles(path string) {
	log.Infof("[strava] Import file %v", path)
	stravaClientId := os.Getenv("STRAVA_CLIENT_ID")
	stravaClientSecret := os.Getenv("STRAVA_CLIENT_SECRET")
	stravaRefreshToken := os.Getenv("STRAVA_REFRESH_TOKEN")

	if len(stravaClientId) == 0 || len(stravaClientSecret) == 0 || len(stravaRefreshToken) == 0 {
		log.Error("[strava] Please set your Strava application parameters in .env")
	}

	ctx := context.Background()
	accessToken, err := GetBearer(ctx, stravaClientId, stravaClientSecret, stravaRefreshToken)
	if err != nil {
		log.WithError(err).Error("[strava] Fail to get Bearer")
	}

	err = upload(*accessToken, path)
	if err != nil {
		log.WithError(err).Error("[strava] Upload failed")
	}
}

func upload(accessToken, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", path)

	gzBuffer := &bytes.Buffer{}
	gzWriter := gzip.NewWriter(gzBuffer)

	_, err = io.Copy(gzWriter, file)
	gzWriter.Close()
	if err != nil {
		return err
	}

	io.Copy(part, gzBuffer)

	writer.WriteField("description", "Uploaded from NRC")
	if strings.HasSuffix(path, ".gpx") {
		writer.WriteField("data_type", "gpx.gz")
	} else if strings.HasSuffix(path, ".tcx") {
		writer.WriteField("data_type", "tcx.gz")
	} else {
		log.Fatal("[strava] Unrecognized file type [%v]", path)
	}

	writer.Close()

	req, err := http.NewRequest("POST", "https://www.strava.com/api/v3/uploads", body)
	req.Header.Add("Content-Type", "multipart/form-data; boundary="+writer.Boundary())
	req.Header.Add("Authorization", "Bearer "+accessToken)


	resp, err := API.GetClient().Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusCreated {
		fmt.Println("Error")
	}

	_, err = ioutil.ReadAll(resp.Body)

	var bodyContent []byte
	resp.Body.Read(bodyContent)

	defer resp.Body.Close()

	log.Infof("[strava] Import done")
	return nil
}

type UnexpectedError struct {
	Errors  []Error `json:"errors"`
	Message string  `json:"message"`
}

type Error struct {
	Code     string `json:"code"`
	Field    string `json:"field"`
	Resource string `json:"resource"`
}
