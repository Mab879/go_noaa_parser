package main

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/twpayne/go-geos"
	pgxgeos "github.com/twpayne/pgx-geos"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"
)

type MsgType int64

const (
	Alert  MsgType = 0
	Update MsgType = 1
	Cancel MsgType = 2
	Ack    MsgType = 3
	Error  MsgType = 4
)

func StringToMsgType(s string) (MsgType, error) {
	switch s {
	case "Alert":
		return Alert, nil
	case "Update":
		return Update, nil
	case "Cancel":
		return Cancel, nil
	case "Ack":
		return Ack, nil
	case "Error":
		return Error, nil
	default:
		return 0, fmt.Errorf("invalid MsgType: %s", s)
	}
}

type Category int64

const (
	Geo Category = iota
	Met
	Safety
	Security
	Rescue
	Fire
	Health
	Env
	Transport
	Infra
	CBRNE
	Other
)

func StringToCategory(s string) (Category, error) {
	switch s {
	case "Geo":
		return Geo, nil
	case "Met":
		return Met, nil
	case "Safety":
		return Safety, nil
	case "Security":
		return Security, nil
	case "Rescue":
		return Rescue, nil
	case "Fire":
		return Fire, nil
	case "Health":
		return Health, nil
	case "Env":
		return Env, nil
	case "Transport":
		return Transport, nil
	case "Infra":
		return Infra, nil
	case "CBRNE":
		return CBRNE, nil
	case "Other":
		return Other, nil
	default:
		return 0, fmt.Errorf("invalid Category: %s", s)
	}

}

type Urgency int64

const (
	Immediate Urgency = iota
	Expected
	Future
	Past
	UnknownUrgency
)

func StringToUrgency(s string) (Urgency, error) {
	switch s {
	case "Immediate":
		return Immediate, nil
	case "Expected":
		return Expected, nil
	case "Future":
		return Future, nil
	case "Past":
		return Past, nil
	case "Unknown":
		return UnknownUrgency, nil
	default:
		return 0, fmt.Errorf("invalid Urgency: %s", s)
	}
}

type Severity int64

const (
	Extreme Severity = iota
	Severe
	Moderate
	Minor
	UnknownSeverity
)

func StringToSeverity(s string) (Severity, error) {
	switch s {
	case "Extreme":
		return Extreme, nil
	case "Severe":
		return Severe, nil
	case "Moderate":
		return Moderate, nil
	case "Minor":
		return Minor, nil
	case "Unknown":
		return UnknownSeverity, nil
	default:
		return 0, fmt.Errorf("invalid Severity: %s", s)
	}
}

type Certainty int64

const (
	Observed Certainty = iota
	Likely
	Possible
	Unlikely
	UnknownCertainty
)

func StringToCertainty(s string) (Certainty, error) {
	switch s {
	case "Observed":
		return Observed, nil
	case "Likely":
		return Likely, nil
	case "Possible":
		return Possible, nil
	case "Unlikely":
		return Unlikely, nil
	case "Unknown":
		return UnknownCertainty, nil
	default:
		return 0, fmt.Errorf("invalid Certainty: %s", s)
	}
}

type Status int64

const (
	Actual Status = iota
	Exercise
	System
	Test
	Draft
)

func StringToStatus(s string) (Status, error) {
	switch s {
	case "Actual":
		return Actual, nil
	case "Exercise":
		return Exercise, nil
	case "System":
		return System, nil
	case "Test":
		return Test, nil
	case "Draft":
		return Draft, nil
	default:
		return 0, fmt.Errorf("invalid Status: %s", s)
	}
}

// Feed was generated 2024-05-22 11:48:09 by https://xml-to-go.github.io/ in Ukraine.
type Feed struct {
	XMLName   xml.Name  `xml:"feed"`
	Text      string    `xml:",chardata"`
	Xmlns     string    `xml:"xmlns,attr"`
	Cap       string    `xml:"cap,attr"`
	ID        string    `xml:"id"`
	Generator string    `xml:"generator"`
	Updated   time.Time `xml:"updated"`
	Author    struct {
		Text string `xml:",chardata"`
		Name string `xml:"name"`
	} `xml:"author"`
	Title string `xml:"title"`
	Link  struct {
		Text string `xml:",chardata"`
		Rel  string `xml:"rel,attr"`
		Href string `xml:"href,attr"`
	} `xml:"link"`
	Entry []struct {
		Text string `xml:",chardata"`
		ID   string `xml:"id"`
		Link struct {
			Text string `xml:",chardata"`
			Rel  string `xml:"rel,attr"`
			Href string `xml:"href,attr"`
		} `xml:"link"`
		Updated   time.Time `xml:"updated"`
		Published time.Time `xml:"published"`
		Author    struct {
			Text string `xml:",chardata"`
			Name string `xml:"name"`
		} `xml:"author"`
		Title     string    `xml:"title"`
		Summary   string    `xml:"summary"`
		Event     string    `xml:"event"`
		Sent      time.Time `xml:"sent"`
		Effective time.Time `xml:"effective"`
		Onset     time.Time `xml:"onset"`
		Expires   time.Time `xml:"expires"`
		Status    string    `xml:"status"`
		MsgType   string    `xml:"msgType"`
		Category  string    `xml:"category"`
		Urgency   string    `xml:"urgency"`
		Severity  string    `xml:"severity"`
		Certainty string    `xml:"certainty"`
		AreaDesc  string    `xml:"areaDesc"`
		Polygon   string    `xml:"polygon"`
		Geocode   []struct {
			Text      string `xml:",chardata"`
			ValueName string `xml:"valueName"`
			Value     string `xml:"value"`
		} `xml:"geocode"`
		Parameter []struct {
			Text      string `xml:",chardata"`
			ValueName string `xml:"valueName"`
			Value     string `xml:"value"`
		} `xml:"parameter"`
	} `xml:"entry"`
}

type AlertObj struct {
	nwsUrl       string
	alertCap     string
	updated      time.Time
	published    time.Time
	authorName   string
	title        string
	summary      string
	capEvent     string
	capEffective time.Time
	capExpires   time.Time
	capStatus    Status
	link         string
	capMsgType   MsgType
	capCategory  Category
	capUrgency   Urgency
	capSeverity  Severity
	capCertainty Certainty
	capAreadesc  string
	capPolygon   string
	capGeocode   map[string]string
	capParameter map[string]string
	updatedAt    time.Time
	createdAt    time.Time
	dbId         int
}

type QueryLogger struct{}

func (a *AlertObj) InsertAlert(pool *pgxpool.Pool) {
	paramsJson, err := json.Marshal(a.capParameter)
	if err != nil {
		panic(err)
	}
	geocodeJson, err := json.Marshal(a.capGeocode)
	if err != nil {
		panic(err)
	}
	// insert query
	insertQuery := `insert into alerts (nws_url, updated, published, author_name, title, 
										summary, cap_event, cap_effective, cap_expires, cap_status, link, cap_msgtype, 
                    					cap_category, cap_urgency, cap_severity, cap_certainty, cap_areadesc, cap_polygon, 
                    					cap_geocode, cap_parameter, created_at, updated_at) 
						VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22)`

	insertQuery = strings.ReplaceAll(insertQuery, "\n", "")
	insertQuery = strings.ReplaceAll(insertQuery, "\t", "")
	var polygon any = nil
	if a.capPolygon != "" {
		polygon = a.capPolygon
	}
	r, err := pool.Query(context.Background(), insertQuery,
		a.nwsUrl, a.updated, a.published, a.authorName, a.title,
		a.summary, a.capEvent, a.capEffective, a.capExpires, a.capStatus, a.link, a.capMsgType,
		a.capCategory, a.capUrgency, a.capSeverity, a.capCertainty, a.capAreadesc, polygon,
		geocodeJson, paramsJson, a.createdAt, a.updatedAt)
	if err != nil {
		panic(err)
	}
	defer r.Close()
	for r.Next() {
		log.Println("Got an row when inserting.")
	}
}

func (a *AlertObj) UpdateAlert(pool *pgxpool.Pool) {
	paramsJson, err := json.Marshal(a.capParameter)
	if err != nil {
		panic(err)
	}
	geocodeJson, err := json.Marshal(a.capGeocode)
	if err != nil {
		panic(err)
	}
	var polygon any = nil
	if a.capPolygon != "" {
		polygon = a.capPolygon
	}
	updateQuery :=
		`UPDATE alerts SET 
                  updated = $1, published = $2, author_name = $3, title = $4, summary = $5, cap_event = $6, 
                  cap_effective = $7, cap_expires = $8, cap_status = $9, link = $10, cap_msgtype = $11,
                  cap_category = $12, cap_urgency = $13, cap_severity = $14, cap_certainty = $15, cap_areadesc = $16, 
                  cap_polygon = $17, cap_geocode = $18, cap_parameter = $19, updated_at = $20
			 WHERE nws_url = $21;`
	updateQuery = strings.ReplaceAll(updateQuery, "\n", "")
	updateQuery = strings.ReplaceAll(updateQuery, "\t", "")
	r, err := pool.Query(context.Background(), updateQuery,
		a.updatedAt, a.published, a.authorName, a.title, a.summary, a.capEvent,
		a.capEffective, a.capExpires, a.capStatus, a.link, a.capMsgType,
		a.capCategory, a.capUrgency, a.capSeverity, a.capCertainty, a.capAreadesc,
		polygon, geocodeJson, paramsJson, a.updatedAt, a.nwsUrl)
	if err != nil {
		panic(err)
	}
	defer r.Close()
	for r.Next() {
		log.Println("Got an row when updating.")
	}
}

func GetUrlContents(url string) []byte {
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	content, err := io.ReadAll(response.Body)
	response.Body.Close()
	if err != nil {
		panic(err)
	}
	return content
}

func FormatPolygon(polygon string) string {
	polygon = strings.ReplaceAll(polygon, " ", "T")
	polygon = strings.ReplaceAll(polygon, ",", " ")
	polygon = strings.ReplaceAll(polygon, "T", ",")
	cords := strings.Split(polygon, ",")
	finalPolygon := ""
	for cord := range cords {
		parts := strings.Split(cords[cord], " ")
		finalPolygon += parts[1] + " " + parts[0] + ", "
	}
	finalPolygon = finalPolygon[:len(finalPolygon)-2]
	return "POLYGON((" + finalPolygon + "))"
}

func GetDatabasePool(databaseUrl string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(databaseUrl)
	if err != nil {
		return nil, err
	}
	config.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		if err := pgxgeos.Register(ctx, conn, geos.NewContext()); err != nil {
			return err
		}
		return nil
	}
	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}
	return pool, nil
}

func UpdateAlertCap(pool *pgxpool.Pool, nws_url string) {
	alertCap := GetUrlContents(nws_url)
	updateQuery := `update alerts set alert_cap = $1 where nws_url = $2`
	_, err := pool.Exec(context.Background(), updateQuery, alertCap, nws_url)
	if err != nil {
		panic(err)
	}

}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	databaseUrl := os.Getenv("DATABASE_URL")
	if databaseUrl == "" {
		panic("DATABASE_URL is not set")
	}
	pool, err := GetDatabasePool(databaseUrl)
	if err != nil {
		panic(err)
	}
	defer pool.Close()
	start := time.Now()
	fmt.Println("Starting Import at: ", start.Format(time.Stamp))
	alertUrl := "https://api.weather.gov/alerts/active.atom"
	rawXmlData := GetUrlContents(alertUrl)
	fmt.Println("Get rawXmlData", len(rawXmlData))
	feed := new(Feed)
	err = xml.Unmarshal(rawXmlData, feed)
	logger.Debug("Feed has been unmarshalled")
	if err != nil {
		panic(err)
	}
	fmt.Println("Feed Last Updated", feed.Updated)
	for _, entry := range feed.Entry {
		fmt.Errorf("Entry: %v", entry.Title)
		id := entry.ID
		fmt.Println("ID: ", id)
		alertCount := 0
		err = pool.QueryRow(context.Background(), "select count(1) from alerts where nws_url = $1", id).Scan(&alertCount)
		if err != nil && err != pgx.ErrNoRows {
			panic(err)
		}
		alert := new(AlertObj)
		alert.nwsUrl = entry.ID
		alert.updated = entry.Updated
		alert.published = entry.Published
		alert.authorName = entry.Author.Name
		alert.title = entry.Title
		alert.summary = entry.Summary
		alert.capEvent = entry.Event
		alert.capEffective = entry.Effective
		alert.capExpires = entry.Expires
		alert.link = entry.Link.Href
		alert.capStatus, err = StringToStatus(entry.Status)
		if err != nil {
			log.Fatal(err)
		}
		alert.capMsgType, err = StringToMsgType(entry.MsgType)
		if err != nil {
			log.Fatal()
		}
		alert.capCategory, err = StringToCategory(entry.Category)
		if err != nil {
			log.Fatal(err)
		}
		alert.capUrgency, err = StringToUrgency(entry.Urgency)
		if err != nil {
			log.Fatal(err)
		}
		alert.capSeverity, err = StringToSeverity(entry.Severity)
		if err != nil {
			log.Fatal(err)
		}
		alert.capCertainty, err = StringToCertainty(entry.Certainty)
		if err != nil {
			log.Fatal(err)
		}
		alert.capAreadesc = entry.AreaDesc
		alert.capGeocode = make(map[string]string)
		for _, geocode := range entry.Geocode {
			alert.capGeocode[geocode.ValueName] = geocode.Value
		}
		alert.capParameter = make(map[string]string)
		for _, parameter := range entry.Parameter {
			alert.capParameter[parameter.ValueName] = parameter.Value
		}
		if entry.Polygon != "" {
			alert.capPolygon = FormatPolygon(entry.Polygon)
		} else {
			alert.capPolygon = ""
		}
		now := time.Now()
		alert.updatedAt = now
		alert.createdAt = now
		fmt.Errorf("alert done creating: %v", alert.nwsUrl)
		// Insert into database
		if alertCount > 0 {
			alert.UpdateAlert(pool)
		} else {
			alert.InsertAlert(pool)
		}
		//UpdateAlertCap(pool, alert.nwsUrl)
	}
}
