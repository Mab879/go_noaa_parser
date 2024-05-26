package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"github.com/jackc/pgx/v5"
	"io"
	"log"
	"net/http"
	"os"
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

type Category int64

const (
	Geo Category = 0
	Met Category = 1
)

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
		Status    int       `xml:"status"`
		MsgType   int       `xml:"msgType"`
		Category  int       `xml:"category"`
		Urgency   int       `xml:"urgency"`
		Severity  int       `xml:"severity"`
		Certainty int       `xml:"certainty"`
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
	capStatus    int
	link         string
	capMsgtype   int
	capCategory  int
	capUrgency   int
	capSeverity  int
	capCertainty int
	capAreadesc  string
	capPolygon   string
	capGeocode   string
	capParameter string
	updatedAt    time.Time
	createdAt    time.Time
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

func main() {
	connection, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}
	defer func(connection *pgx.Conn, ctx context.Context) {
		err := connection.Close(ctx)
		if err != nil {
			panic(err)
		}
	}(connection, context.Background())
	start := time.Now()
	fmt.Println("Starting Import at: ", start.Format(time.Stamp))
	alertUrl := "https://api.weather.gov/alerts/active.atom"
	rawXmlData := GetUrlContents(alertUrl)
	fmt.Println("Get rawXmlData", len(rawXmlData))
	feed := new(Feed)
	err = xml.Unmarshal(rawXmlData, feed)
	if err != nil {
		panic(err)
	}
	fmt.Println("Feed Last Updated", feed.Updated)
	for _, entry := range feed.Entry {
		id := entry.ID
		var alertCount = 0
		err = connection.QueryRow(context.Background(), "select count(1) from alerts where nws_url = $1", id).Scan(&alertCount)
		if err != nil {
			panic(err)
		}
		alert := new(AlertObj)
		if alertCount != 0 {
			err = connection.QueryRow(context.Background(), "select * from alerts where nws_url = $1", id).Scan(&alert)
			continue
		}
		alert.nwsUrl = entry.ID
		alert.updated = entry.Updated
		alert.published = entry.Published
		alert.authorName = entry.Author.Name
		alert.title = entry.Title
		alert.summary = entry.Summary
		alert.capEvent = entry.Event
		alert.capEffective = entry.Effective
		alert.capExpires = entry.Expires
		alert.capStatus = entry.Status
		alert.capMsgtype = entry.MsgType
		alert.capCategory = entry.Category
		alert.capUrgency = entry.Urgency
		alert.capSeverity = entry.Severity
		alert.capCertainty = entry.Certainty
		alert.capAreadesc = entry.AreaDesc

	}
}
