package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

var (
	outputFmt = "15:04"
)

type Entry struct {
	Start    time.Time
	End      time.Time
	Title    string
	Location string
	AllDay   bool
}

func (e Entry) String() string {
	sb := strings.Builder{}

	sb.WriteString(e.StartEnd())
	sb.WriteString(fmt.Sprintf(" - %s", e.Title))

	if e.Location != "" {
		sb.WriteString(fmt.Sprintf(" (%s)", e.LocationString()))
	}

	return sb.String()
}

func (e Entry) StartEnd() string {
	if e.AllDay {
		return "all day"
	}

	return fmt.Sprintf("%s - %s", e.Start.Format(outputFmt), e.End.Format(outputFmt))
}

func (e Entry) LocationString() string {
	if e.Location == "" {
		return ""
	}

	return fmt.Sprintf("@ %s", e.Location)
}

type Google struct {
	calendar *calendar.Service
	timezone *time.Location
}

func NewGoogle(creds, tokens, tz string) (g Google, err error) {
	g.timezone, err = time.LoadLocation(tz)
	if err != nil {
		return
	}

	ctx := context.Background()

	b, err := os.ReadFile(creds)
	if err != nil {
		return
	}

	config, err := google.ConfigFromJSON(b, calendar.CalendarReadonlyScope)
	if err != nil {
		return
	}

	client, err := getClient(tokens, config)
	if err != nil {
		return
	}

	g.calendar, err = calendar.NewService(ctx, option.WithHTTPClient(client))

	return
}

func (g Google) Today() (e []Entry, err error) {
	bod := g.bod()
	min := bod.Format(time.RFC3339)
	max := bod.Add(24 * time.Hour).Format(time.RFC3339)

	events, err := g.calendar.Events.List("primary").ShowDeleted(false).
		SingleEvents(true).TimeMin(min).TimeMax(max).MaxResults(10).OrderBy("startTime").Do()

	if err != nil {
		return
	}

	e = make([]Entry, len(events.Items))

	for idx, item := range events.Items {
		e[idx] = Entry{
			Title:    item.Summary,
			Location: item.Location,
		}

		switch item.Start.DateTime {
		case "":
			e[idx].AllDay = true

			e[idx].Start, err = time.Parse("2006-01-02", item.Start.Date)
			if err != nil {
				return
			}

		default:
			e[idx].Start, err = time.Parse(time.RFC3339, item.Start.DateTime)
			if err != nil {
				return
			}

			e[idx].End, err = time.Parse(time.RFC3339, item.End.DateTime)
			if err != nil {
				return
			}
		}
	}

	return
}

// Return Beginning of Day
func (g Google) bod() (t time.Time) {
	t = time.Now().In(g.timezone)
	t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, g.timezone)

	return
}

func getClient(tokens string, config *oauth2.Config) (c *http.Client, err error) {
	tok, err := tokenFromFile(tokens)
	if err != nil {
		return
	}

	c = config.Client(context.Background(), tok)

	return
}

func tokenFromFile(file string) (tok *oauth2.Token, err error) {
	f, err := os.Open(file)
	if err != nil {
		return
	}

	defer f.Close()

	tok = &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)

	return
}
