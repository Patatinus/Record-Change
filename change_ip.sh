#!/bin/bash
# You can use http://127.0.0.1:58371 if you are running
# this API on the same device
API_URL="https://Your-Server-URL"
API_TOKEN="API-token-you-set-in-env"


case $1 in

    "--change-ip")
        ip="$(curl -s https://ifconfig.net)"
        echo "Changing record's ip to \"$ip\""
        curl -X PATCH "$API_URL/currentIP/" \
            -H "Authorization: Bearer $API_TOKEN" \
            -H "Content-Type: application/json" \
            -d "{\"ip\": \"$ip\"}"
    ;;
    "--change-ip-manual")
        ip=$2
        echo "Changing record's ip to \"$ip\""
        curl -X PATCH "$API_URL/currentIP/" \
            -H "Authorization: Bearer $API_TOKEN" \
            -H "Content-Type: application/json" \
            -d "{\"ip\": \"$ip\"}"
    ;;
    "--change-cname")
        cname=$2
        echo "Changing record's ip to \"$cname\""
        curl -X PATCH "$API_URL/api/cname/" \
            -H "Authorization: Bearer $API_TOKEN" \
            -H "Content-Type: application/json" \
            -d "{\"ip\": \"$cname\"}"
    ;;
    "--get-current")
        curl -X GET "$API_URL/currentIP/" \
            -H "Authorization: Bearer $API_TOKEN"
    ;;
    *)
        echo "Usage:"
        echo "--change-ip                           automatically changes the record into an A record and sets ip to the current ip of the device."
        echo "--change-ip-manual <ip address>       automatically changes the record into an A record and sets ip to the given ip address."
        echo "--change-cname <fqdn>                 automatically changes the record into an CNAME record and sets target to the given FQDN address."
        echo "--get-current                         gets currently applied ip address or FQDN of the record."
    ;;
esac
