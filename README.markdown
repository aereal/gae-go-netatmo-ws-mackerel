# gae-go-netatmo-ws-mackerel

Fetch metrics from [Netatmo Weather Station][] and post them to [Mackerel][] on [GAE/Go][].

## Overview

This web application collects metrics from Netatmo and send them to Mackerel if requested on `/postMetrics`.

If `cron.yaml` is deployed, GAE cron jobs run `/postMetrics` every 10 minutes.

## Environment varibles

These variables should be written in `secret.yaml` (that is ignored and you should create it).

- `NETATMO_EMAIL` - Email address of Netatmo account
- `NETATMO_PASSWORD` - Password of Netatmo account
- `NETATMO_APP_ID` - application id that registered on [Netatmo Connect][]
- `NETATMO_APP_SECRET` - application secret that registered on [Netatmo Connect][]
- `MACKEREL_APIKEY` - API key that issued on [Mackerel][]
- `MACKEREL_SERVICE_NAME` - The service name on [Mackerel][]

## Deployment

```sh
gcloud app deploy app.yaml
gcloud app deploy cron.yaml
```

[Mackerel]: https://mackerel.io/
[Netatmo Weather Station]: https://www.netatmo.com/product/weather/weatherstation
[GAE/Go]: https://cloud.google.com/appengine/docs/standard/go/
[Netatmo Connect]: https://dev.netatmo.com/
