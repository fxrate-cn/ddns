# cloudflare DDNS golang script

A tiny golang that updates a cloudflare dns record with your current ip.


In the current setup, it is expected that you clone this repo to any directory `${DIRECTORY}`. For example '/home/user/dev'

You need to put an `.env` file into `${DIRECTORY}/ddns` (for example /home/${user}/dev/ddns) that looks like `.env-example`.
You need input:
* apiToken: get you token from https://dash.cloudflare.com/profile/api-tokens
* zoneId: Zone ID of your domain
* name: subdomain name like "www.fxrate.cn"
  
Run `go mod tidy` to install dependecies. Run `go mod tidy` to build application.
You can then register a cronjob executing `${DIRECTORY}/ddns/ddns` (for example /home/${user}/dev/ddns/ddns) in an interval of choice.

```
0 6 * * * /home/${user}/dev/ddns/ddns
```

