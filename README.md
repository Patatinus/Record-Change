# Record-Change
[![Go Build Check](https://github.com/Patatinus/Record-Change/actions/workflows/go.yml/badge.svg)](https://github.com/Patatinus/Record-Change/actions)
## About The Project
This is a simple API to edit and view the context of a record on Cloudflare. It utilizes Cloudflare's API and is built in Go.

This project can be used to easily change a domain's IP address. Great for self hosted services with dynamic IP addresses. Like lending a subdomain for your friend's self hosted Minecraft server or website. CNAME is to mask services like e4mc that gives a free and relatively bad looking domains.

Only works with IPv4 and CNAME for now.

_Don't tell anyone, but I made this for Cloudflare only because their API is free to use._


## Getting Started

### Prerequisites
- Docker
- A device that can run Docker
- A domain that is managed by Cloudflare

### Installation

1. Clone the repo
   ```sh
   git clone https://github.com/Patatinus/Record-Change.git
   ```
2. Copy the .env.example file
   ```sh
   cp .env.example .env
   ```
3. Enter your Cloudflare API, Zone ID and Domain ID in `.env`
   ```txt
   CF_API_TOKEN=cloudflare-api-token
   CF_ZONE_ID=cloudflare-zone-id
   CF_DNS_RECORD_ID=cloudflare-dns-record-id
   ```
4. Generate a token to ensure authentication and enter it in `.env`. If you ever suspect it is compromised, generate a new one
   ```sh
   openssl rand -base64 32 # copy the output token
   ```
   ```txt
   API_TOKEN=the-token-your-server-will-use-to-authorise
   ```
5. Run the API  with Docker
   ```sh
   docker compose up 
   # for the first run "docker compose up --build"
   # it runs on 58371 port, which practically nobody else uses
   ```





## Usage

App can be used by editing the change_ip.sh and adding the user generated API token.

### curl examples
- Request to change the record's IP
```sh
curl -X PATCH "https://Your-Server-URL/currentIP/" \
            -H "Authorization: Bearer $API_TOKEN" \
            -H "Content-Type: application/json" \
            -d "{\"ip\": \"$ip\"}"
``` 
- 200 outcome:
```json
{"message":"Successfully changed the IP to 10.10.10.10"}
```
- 400 outcome due to bad IPv4 address:
```json
{"error":"Provided string is not a valid IPv4 address"}
```
- Request to get the record's current IP
```sh
curl -X GET "https://Your-Server-URL/currentIP/" \
            -H "Authorization: Bearer $API_TOKEN"
```
- 200 outcome:
```json
{"content":"10.10.10.10"}
```
- Request to change the record into CNAME and assigning a FQDN
```sh
curl -X PATCH "https://Your-Server-URL/api/cname/" \
            -H "Authorization: Bearer $API_TOKEN" \
            -H "Content-Type: application/json" \
            -d "{\"ip\": \"$cname\"}"
```
- 200 outcome:
```json
{"message":"Successfully changed the target domain to very.real.domain"}
```
- 400 outcome due to bad IPv4 address:
```json
{"error":"Provided string is not a valid FQDN"}
```

### Usage of change_ip.sh
1. Give execute permissions to the script
```sh
sudo chmod +x change_ip.sh
```
2. Fill in API_URL and API_TOKEN. Make sure API_URL ends with `/currentIP/`
```sh
# You can use http://127.0.0.1:58371 if you are running
# this API on the same device
API_URL="https://Your-Server-URL"
API_TOKEN="API-token-you-set-in-env"
```
3. Run the script to see the parameters, then you are ready to go
```sh
./change_ip.sh
```





## License

Distributed under the MIT License. See `LICENSE` for more information.






## Contact

Project Link: [https://github.com/Patatinus/Record-Change](https://github.com/Patatinus/Record-Change)



