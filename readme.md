# Cloudflare DNS Client
DNS Client for Cloudflare API.
Checks your current external IP, and updates the A-record for the specified domain using the Cloudflare API.
You need to have your domain registered with [Cloudflare DNS](https://www.cloudflare.com/en-gb/)



Note that for this to work for game servers, I strongly suggest you disable the proxy functionality for your domain with Cloudflare DNS.

Read more here: [How Cloudflare works](https://developers.cloudflare.com/fundamentals/concepts/how-cloudflare-works/), (Disable proxy on DNS Records)[https://developers.cloudflare.com/fundamentals/setup/manage-domains/pause-cloudflare/#disable-proxy-on-dns-records]


This client is a prime candidate for an automated cronjob, see below example
```crontab

*/5 * * * * /PATH_TO_BIN/cloudflare-dns-client -domain example.org -token ABCDE -zone ABCDE >> /PATH_TO_LOG/dns.log
# This updates your A-record every five minutes
```
See how to format your crontab on [crontab guru](https://crontab.guru/) 
Requires three parameters:

1. -domain This is the domain you want to update (example.org)
2. -token This is your CloudFlare API Token, retrievable from the CloudFlare dashboard. See [CloudFlare Docs](https://developers.cloudflare.com/fundamentals/api/get-started/create-token/) 
3. -zone This is your Zone ID for the domain, retrievable from CloudFlare dashboard. See [CloudFlare Docs](https://developers.cloudflare.com/fundamentals/setup/find-account-and-zone-ids/#find-zone-and-account-ids) 

