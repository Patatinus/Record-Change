package main

import (
	"fmt"
	"net"
	"net/http"
	"strings"

	"os"

	"github.com/cloudflare/cloudflare-go/v7"
	"github.com/cloudflare/cloudflare-go/v7/dns"
	"github.com/cloudflare/cloudflare-go/v7/option"

	"github.com/gin-gonic/gin"
)

type data struct {
	IP string `json:"ip"`
}

func getIP(c *gin.Context) {
	zoneID := os.Getenv("CF_ZONE_ID")
	apiToken := os.Getenv("CF_API_TOKEN")
	recordID := os.Getenv("CF_DNS_RECORD_ID")

	client := cloudflare.NewClient(
		option.WithAPIToken(apiToken),
	)
	recordResponse, err := client.DNS.Records.Get(
		c.Request.Context(),
		recordID,
		dns.RecordGetParams{
			ZoneID: cloudflare.F(zoneID),
		},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read the DNS record."})
		fmt.Printf("Cloudflare API Error: %+v\n", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"content": recordResponse.Content})
	fmt.Printf("%+v\n", recordResponse)

}

func patchIP(c *gin.Context) {

	var newIP data
	if err := c.ShouldBindJSON(&newIP); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cleanIP := strings.TrimSpace(newIP.IP)

	ip := net.ParseIP(cleanIP)
	if ip == nil || ip.To4() == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Provided string is not a valid IPv4 address"})
		fmt.Printf("Invalid ip address: %s\n", cleanIP)
		return
	}

	zoneID := os.Getenv("CF_ZONE_ID")
	apiToken := os.Getenv("CF_API_TOKEN")
	recordID := os.Getenv("CF_DNS_RECORD_ID")

	client := cloudflare.NewClient(
		option.WithAPIToken(apiToken),
	)

	recordResponse, err := client.DNS.Records.Edit(
		c.Request.Context(),
		recordID,
		dns.RecordEditParams{
			ZoneID: cloudflare.F(zoneID),
			Body: dns.ARecordParam{
				Type:    cloudflare.F(dns.ARecordTypeA),
				Content: cloudflare.F(cleanIP),
			},
		},
	)
	if err == nil {
		c.JSON(http.StatusOK, gin.H{"message": "Successfully changed the IP to " + cleanIP})
		fmt.Printf("Updated Record: %+v\n", recordResponse)
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update the DNS record."})
		fmt.Printf("Cloudflare API Error: %+v\n", err)
	}
}

func checkAuth() gin.HandlerFunc {
	return func(c *gin.Context) {

		expectedToken := os.Getenv("API_TOKEN")
		authHeader := c.GetHeader("Authorization")

		if authHeader != "Bearer "+expectedToken {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Invalid or missing token"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	router.GET("/currentIP/", checkAuth(), getIP)
	router.PATCH("/currentIP/", checkAuth(), patchIP)

	router.Run("0.0.0.0:8080")
}
