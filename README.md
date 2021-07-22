# Twitch prometheus exporter

This exporter will allow you to export statistics of your Twitch account into prometheus, like follower and viewer count.

## Build

To build this project you'll need to have [golang](https://golang.org/) installed.
After that just clone this project and build it with `go build ./cmd/twitch-prometheus/...` while in the directory of your local clone.
Alternatively if you rather use docker, you can just run `docker build -t <tagname> .` in the main directory of this project.

## Usage

To set this up you'll first off have to register an application with Twitch in order to use their API.
In order to do this, head over to <https://dev.twitch.tv/console/apps/create> and fill in the form, for the OAuth redirect url just fill in `https://localhost`, as this isn't used by the exporter anyway.
Category I would recommend to fill in `Analytics Tool`.
After having done this, press the Manage button for your new application.
You should see a Client ID field, now press the New Secret button right below this.
Now create a new config.yml based on config.example.yml and fill in the Client ID and Client Secret you have just received.
Next up just fill in the port you'll want the server to listen on and fill in a list of twitch users to gather statistics from and you're ready to go.

## Prometheus

Finally you'll want to configure prometheus to actually query the exporter.
This is basically just adding the following to your prometheus config (likely located at `/etc/prometheus/prometheus.yml`).
Of course in case you decide to run on a different port you'll have to adjust this.

```yml
scrape_configs:
  - job_name: "twitch"
    static_configs:
      - targets: ['127.0.0.1:9001']
```

## Exported metrics

And finally, here's a brief example of what metrics you can expect.
Do of course understand that this data comes directly from the Twitch API, which might be a bit off sometimes.

```bash
# HELP twitch_followers_total Number of Followers
# TYPE twitch_followers_total counter
twitch_followers_total{username="schoentoon"} 27
# HELP twitch_viewers_total Number of Viewers
# TYPE twitch_viewers_total counter
twitch_viewers_total{username="schoentoon"} 10
```
