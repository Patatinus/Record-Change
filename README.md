# Record-Change
[![Go Build Check](https://github.com/Patatinus/Record-Change/actions/workflows/go.yml/badge.svg)](https://github.com/Patatinus/Record-Change/actions)
## About The Project
This is a simple API to edit and view the context of a record on Cloudflare. It utilizes Cloudflare's API and is built in Go.

This project can be used to easily change a domain's IP address. Great for self hosted services with dynamic IP address, like lending a subdomain for your friend's self hosted Minecraft server or website. SRV and CNAME are to mask services like e4mc that give free and relatively bad looking domains.

Only works with A, CNAME and SRV records for now.

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
3. Enter your Cloudflare API Token, Zone ID, Record ID, RECORD_FQDN (FQDN you want to use), SERVICE_TYPE (for e4mc like services) and API_PORT in `.env`
   ```txt
   CF_API_TOKEN=cloudflare-api-token
   CF_ZONE_ID=cloudflare-zone-id
   CF_DNS_RECORD_ID=cloudflare-dns-record-id

   RECORD_FQDN=www.example.com
   SERVICE_TYPE=_service._type #_minecraft._tcp for minecraft

   API_PORT=58371
   ```
4. Generate a token to ensure authentication and enter it in `.env`. If you ever suspect it is compromised, generate a new one
   ```sh
   openssl rand -base64 32 # copy the output token
   ```
   ```txt
   API_TOKEN=the-token-your-server-will-use-to-authorise
   ```
5. Run the API with Docker
   ```sh
   sudo docker compose up -d
   # for the first run "sudo docker compose up --build -d"
   # it runs on 58371 port by default, you can change it in .env
   ```





## Usage

App can be used by editing the change_ip.sh and adding the user generated API token.

### curl examples
- Request to change the record's IP
```sh
curl -X PATCH "https://Your-Server-URL/api/a/" \
            -H "Authorization: Bearer $API_TOKEN" \
            -H "Content-Type: application/json" \
            -d "{\"content\": \"$ip\"}"
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
curl -X GET "https://Your-Server-URL/api/current/" \
            -H "Authorization: Bearer $API_TOKEN"
```
- 200 outcome:
```json
{"content":"10.10.10.10"}
```
- Request to change the record into CNAME and assign a FQDN
```sh
curl -X PATCH "https://Your-Server-URL/api/cname/" \
            -H "Authorization: Bearer $API_TOKEN" \
            -H "Content-Type: application/json" \
            -d "{\"content\": \"$FQDN\"}"
```
- 200 outcome:
```json
{"message":"Successfully changed the target domain to very.real.domain"}
```
- 400 outcome due to bad FQDN:
```json
{"error":"Provided string is not a valid FQDN"}
```
- Request to change the record into SRV and assign a FQDN and a port
```sh
curl -X PATCH "https://Your-Server-URL/api/srv/" \
            -H "Authorization: Bearer $API_TOKEN" \
            -H "Content-Type: application/json" \
            -d "{\"content\": \"$FQDN\", \"port\": \"$PORT\"}"
```
- 200 outcome:
```json
{"message":"Successfully changed the target domain to www.example2.com with the port 25565"}
```
- 400 outcome due to bad port number:
```json
{"error":"Provided integer is not a valid port number"}
```

### Usage of change_ip.sh
1. Give execute permissions to the script
```sh
sudo chmod +x change_ip.sh
```
2. Fill in API_URL and API_TOKEN. Make sure API_URL doesn't end with a "/"
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



