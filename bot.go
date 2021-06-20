package main

import (
	"strings"
	"time"

	"github.com/jspc/bottom"
	"github.com/lrstanley/girc"
	"github.com/olekukonko/tablewriter"
)

type Bot struct {
	bottom bottom.Bottom
	google Google
}

func New(user, password, server string, verify bool, g Google) (b Bot, err error) {
	b.google = g

	b.bottom, err = bottom.New(user, password, server, verify)
	if err != nil {
		return
	}

	b.bottom.Client.Handlers.Add(girc.CONNECTED, func(c *girc.Client, e girc.Event) {
		c.Cmd.Join(Chan)
	})

	router := bottom.NewRouter()
	router.AddRoute(`(?i)show\s+my\s+diary`, b.diary)

	b.bottom.Middlewares.Push(router)

	return
}

func (b Bot) diary(_, channel string, _ []string) (err error) {
	events, err := b.google.Today()
	if err != nil {
		return
	}

	sb := strings.Builder{}

	table := tablewriter.NewWriter(&sb)
	table.SetHeader([]string{"", "Event", "Location"})

	for _, line := range events {
		table.Append([]string{line.StartEnd(), line.Title, line.LocationString()})
	}

	table.Render()

	b.bottom.Client.Cmd.Messagef(channel, "Diary for %s", time.Now().In(b.google.timezone).Format("2006-01-02"))

	for _, line := range strings.Split(sb.String(), "\n") {
		b.bottom.Client.Cmd.Message(channel, line)
	}

	return
}
