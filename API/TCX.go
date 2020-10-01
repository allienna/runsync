package API

type TrainingCenterDatabase struct {
	SchemaLocation string `xml:"xsi:schemaLocation,attr"`
	Xmlns          string `xml:"xmlns,attr"`
	XmlnsNs5       string `xml:"xmlns:n5,attr"`
	XmlnsNs4       string `xml:"xmlns:n4,attr"`
	XmlnsNs3       string `xml:"xmlns:n3,attr"`
	XmlnsNs2       string `xml:"xmlns:n2,attr"`
	XmlnsXsi       string `xml:"xmlns:xsi,attr"`

	Activities Activities `xml:"Activities"`
	Author     Author     `xml:"Author"`
}
type Activities struct {
	Activities []Activity `xml:"Activity"`
}

type Activity struct {
	Sport string `xml:"Sport,attr"`

	ID  string `xml:"Id"`
	Lap Lap    `xml:"Lap"`
}

type Lap struct {
	StartTime string `xml:"StartTime,attr"`

	TotalTimeSeconds    float32       `xml:"TotalTimeSeconds"`
	DistanceMeters      float32       `xml:"DistanceMeters"`
	MaximumSpeed        float32       `xml:"MaximumSpeed"`
	Calories            int32         `xml:"Calories"`
	AverageHeartRateBpm Value         `xml:"AverageHeartRateBpm"`
	MaximumHeartRateBpm Value         `xml:"MaximumHeartRateBpm"`
	Intensity           string        `xml:"Intensity"`
	TriggerMethod       string        `xml:"TriggerMethod"`
	Track               TcxTrack      `xml:"Track"`
	Extensions          LapExtensions `xml:"Extensions"`
}

type Value struct {
	Value int32 `xml:"Value,omitempty"`
}

type TcxTrack struct {
	Trackpoint []TcxTrackpoint `xml:"Trackpoint"`
}

type TcxTrackpoint struct {
	Time           string         `xml:"Time"`
	DistanceMeters float32        `xml:"DistanceMeters"`
	HeartRateBpm   *Value         `xml:"HeartRateBpm,omitempty"`
	Extensions     TrackExtension `xml:"Extensions"`
}

type TrackExtension struct {
	TPX TPX `xml:"TPX"`
}

type TPX struct {
	Xmlns string  `xml:"xmlns,attr"`
	Speed float32 `xml:"Speed"`
	//RunCadence int32   `xml:"RunCadence"`
}
type LapExtensions struct {
	LX LX `xml:"LX"`
}

type LX struct {
	AvgSpeed      float32 `xml:"AvgSpeed,omitempty"`
	AvgRunCadence int32   `xml:"AvgRunCadence,omitempty"`
	MaxRunCadence int32   `xml:"MaxRunCadence,omitempty"`
}

type Author struct {
	Type   string `xml:"xsi:type,attr"`
	Name   string `xml:"Name"`
	Build  Build  `xml:"Build"`
	LangID string `xml:"LangID"`
}

type Build struct {
	Version Version `xml:"Version"`
}

type Version struct {
	VersionMajor int `xml:"VersionMajor"`
	VersionMinor int `xml:"VersionMinor"`
	BuildMajor   int `xml:"BuildMajor"`
	BuildMinor   int `xml:"BuildMinor"`
}
