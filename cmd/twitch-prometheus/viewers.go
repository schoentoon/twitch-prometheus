package main

import (
	"log"

	"github.com/nicklaw5/helix"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	viewersCountDesc = prometheus.NewDesc(
		"viewers_total",
		"Number of Viewers",
		[]string{"username"}, nil,
	)
)

type ViewerCollector struct {
	client  *helix.Client
	userIDs []string
}

func NewViewersCollector(client *helix.Client, usernames []string) (*ViewerCollector, error) {
	resp, err := client.GetUsers(&helix.UsersParams{
		Logins: usernames,
	})
	if err != nil {
		return nil, err
	}

	userIDs := make([]string, 0, len(usernames))
	for _, user := range resp.Data.Users {
		userIDs = append(userIDs, user.ID)
	}

	return &ViewerCollector{
		client:  client,
		userIDs: userIDs,
	}, nil
}

func (vc ViewerCollector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(vc, ch)
}

func (vc ViewerCollector) Collect(ch chan<- prometheus.Metric) {
	out, err := vc.client.GetStreams(&helix.StreamsParams{
		UserIDs: vc.userIDs,
	})
	if err != nil {
		log.Println(err)
	}

	for _, stream := range out.Data.Streams {
		ch <- prometheus.MustNewConstMetric(
			viewersCountDesc,
			prometheus.CounterValue,
			float64(stream.ViewerCount),
			stream.UserName,
		)
	}
}
