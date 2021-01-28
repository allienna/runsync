package API

import (
	"encoding/xml"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
)

type GPX struct {
	Creator        string   `xml:"creator,attr"`
	XmlnsXsi       string   `xml:"xmlns:xsi,attr"`
	SchemaLocation string   `xml:"xsi:schemaLocation,attr"`
	Version        string   `xml:"version,attr"`
	Xmlns          string   `xml:"xmlns,attr"`
	XmlnsGpxtpx    string   `xml:"xmlns:gpxtpx,attr"`
	XmlnsGpxx      string   `xml:"xmlns:gpxx,attr"`
	Metadata       Metadata `xml:"metadata"`
	Track          Track    `xml:"trk"`
}
type Metadata struct {
	Time string `xml:"time"`
}

type Track struct {
	Name         string       `xml:"name"`
	Type         int          `xml:"type"`
	TrackSegment TrackSegment `xml:"trkseg"`
}

type TrackSegment struct {
	TrackPoints []TrackPoint `xml:"trkpt"`
}

type TrackPoint struct {
	Latitude   string       `xml:"lat,attr"`
	Longitude  string       `xml:"lon,attr"`
	Time       string       `xml:"time"`
	Elevation  string       `xml:"ele"`
	Extensions []Extensions `xml:"extensions"`
	Start      int64
}
type Extensions struct {
	TrackPointExtensions []TrackPointExtension `xml:"gpxtpx:TrackPointExtension"`
}

type TrackPointExtension struct {
	HeartRate int `xml:"gpxtpx:hr"`
}

func WriteGpxToFile(activityID string, gpx *GPX) string {
	file, err := xml.MarshalIndent(&gpx, "", " ")
	if err != nil {
		log.Errorf("Fail to transform struct to XML for id [%v] : %v", activityID, err)
	}

	os.Mkdir("activities", os.ModePerm)
	path := fmt.Sprintf("./activities/activity_%v.gpx", activityID)
	err = ioutil.WriteFile(path, file, 0644)

	if err != nil {
		log.Errorf("Fail to write file for id [%v], %v", activityID, err)
	}
	return path
}

func WriteTcxToFile(activityID string, tcx *TrainingCenterDatabase) string {
	file, err := xml.MarshalIndent(&tcx, "", " ")
	if err != nil {
		log.Errorf("Fail to transform struct to XML for id [%v] : %v", activityID, err)
	}

	file = []byte(xml.Header + string(file))

	os.Mkdir("activities", os.ModePerm)
	path := fmt.Sprintf("./activities/activity_%v.tcx", activityID)
	err = ioutil.WriteFile(path, file, 0644)

	if err != nil {
		log.Errorf("Fail to write file for id [%v], %v", activityID, err)
	}
	return path
}
