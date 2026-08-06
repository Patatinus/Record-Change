#!/bin/bash
API_URL="https://Your-Server-URL/currentIP/"
API_TOKEN="API-token-you-set-in-env"


case $1 in

    "--change")
        ip="$(curl -s https://ifconfig.net)"
        echo "Changing record's ip to \"$ip\""
        curl -X PATCH "$API_URL" \
            -H "Authorization: Bearer $API_TOKEN" \
            -H "Content-Type: application/json" \
            -d "{\"ip\": \"$ip\"}"
    ;;
    "--change-manual")
        ip=$2
        echo "Changing record's ip to \"$ip\""
        curl -X PATCH "$API_URL" \
            -H "Authorization: Bearer $API_TOKEN" \
            -H "Content-Type: application/json" \
            -d "{\"ip\": \"$ip\"}"
    ;;
    "--get-current")
        (curl -X GET "$API_URL" \
            -H "Authorization: Bearer $API_TOKEN")
    ;;
    *)
        echo "Usage:"
        echo "--change                              automatically changes the ip to the current ip of the device."
        echo "--change-manual <ip address>          manually changes the ip to the current ip of the device."
        echo "--get-current                         gets currently applied ip address of the record."
    ;;
esac
