package main

import (
	"log"
	"sync"

	"github.com/nicklaw5/helix"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	followerCountDesc = prometheus.NewDesc(
		"followers_total",
		"Number of Followers",
		[]string{"username"}, nil,
	)
)

type FollowerCollector struct {
	client *helix.Client
	users  map[string]string // key is the id, value is the username
}

func NewFollowersCollector(client *helix.Client, usernames []string) (*FollowerCollector, error) {
	resp, err := client.GetUsers(&helix.UsersParams{
		Logins: usernames,
	})
	if err != nil {
		return nil, err
	}

	users := make(map[string]string, len(usernames))
	for _, user := range resp.Data.Users {
		users[user.ID] = user.Login
	}

	return &FollowerCollector{
		client: client,
		users:  users,
	}, nil
}

func (fc FollowerCollector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(fc, ch)
}

func (fc FollowerCollector) Collect(ch chan<- prometheus.Metric) {
	var wg sync.WaitGroup

	wg.Add(len(fc.users))
	for id := range fc.users {
		go fc.fetchFollowerCount(id, ch, &wg)
	}

	wg.Wait()
}

func (fc FollowerCollector) fetchFollowerCount(id string, ch chan<- prometheus.Metric, wg *sync.WaitGroup) {
	defer wg.Done()

	out, err := fc.client.GetUsersFollows(&helix.UsersFollowsParams{
		ToID:  id,
		First: 1,
	})
	if err != nil {
		log.Println(err)
	}

	ch <- prometheus.MustNewConstMetric(
		followerCountDesc,
		prometheus.CounterValue,
		float64(out.Data.Total),
		fc.users[id],
	)
}
