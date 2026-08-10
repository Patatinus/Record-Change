package main

import (
	"fmt"
	"net"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"os"

	"github.com/cloudflare/cloudflare-go/v7"
	"github.com/cloudflare/cloudflare-go/v7/dns"
	"github.com/cloudflare/cloudflare-go/v7/option"

	"github.com/gin-gonic/gin"
)

type data struct {
	Content string `json:"content"`
}

type recordSRV struct {
	Content string `json:"content"`
	Port    string `json:"port"`
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

func patchA(c *gin.Context) {

	var newA data
	if err := c.ShouldBindJSON(&newA); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cleanIP := strings.TrimSpace(newA.Content)

	ip := net.ParseIP(cleanIP)
	if ip == nil || ip.To4() == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Provided string is not a valid IPv4 address"})
		fmt.Printf("Invalid ip address: %s\n", cleanIP)
		return
	}

	zoneID := os.Getenv("CF_ZONE_ID")
	apiToken := os.Getenv("CF_API_TOKEN")
	recordID := os.Getenv("CF_DNS_RECORD_ID")
	recordFQDN := os.Getenv("RECORD_FQDN")

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
				Name:    cloudflare.F(recordFQDN),
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

func patchCNAME(c *gin.Context) {

	var newCName data
	if err := c.ShouldBindJSON(&newCName); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cleanCName := strings.TrimSpace(newCName.Content)

	if !isValidFQDN(cleanCName) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Provided string is not a valid FQDN"})
		fmt.Printf("Invalid FQDN: %s\n", cleanCName)
		return
	}

	zoneID := os.Getenv("CF_ZONE_ID")
	apiToken := os.Getenv("CF_API_TOKEN")
	recordID := os.Getenv("CF_DNS_RECORD_ID")
	recordFQDN := os.Getenv("RECORD_FQDN")

	client := cloudflare.NewClient(
		option.WithAPIToken(apiToken),
	)

	recordResponse, err := client.DNS.Records.Edit(
		c.Request.Context(),
		recordID,
		dns.RecordEditParams{
			ZoneID: cloudflare.F(zoneID),
			Body: dns.CNAMERecordParam{
				Type:    cloudflare.F(dns.CNAMERecordTypeCNAME),
				Content: cloudflare.F(cleanCName),
				Name:    cloudflare.F(recordFQDN),
			},
		},
	)
	if err == nil {
		c.JSON(http.StatusOK, gin.H{"message": "Successfully changed the target domain to " + cleanCName})
		fmt.Printf("Updated Record: %+v\n", recordResponse)
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update the DNS record."})
		fmt.Printf("Cloudflare API Error: %+v\n", err)
	}
}

func patchSRV(c *gin.Context) {

	var newSRV recordSRV
	if err := c.ShouldBindJSON(&newSRV); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cleanSRV := strings.TrimSpace(newSRV.Content)
	intPort, err := strconv.Atoi(strings.TrimSpace(newSRV.Port))

	if !isValidFQDN(cleanSRV) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Provided string is not a valid FQDN"})
		fmt.Printf("Invalid FQDN: %s\n", cleanSRV)
		return
	}

	if err != nil || intPort < 0 || intPort > 65535 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Provided integer is not a valid port number"})
		fmt.Printf("Invalid port: %d\n", intPort)
		return
	}

	zoneID := os.Getenv("CF_ZONE_ID")
	apiToken := os.Getenv("CF_API_TOKEN")
	recordID := os.Getenv("CF_DNS_RECORD_ID")
	recordFQDN := os.Getenv("RECORD_FQDN")

	serviceType := os.Getenv("SERVICE_TYPE")

	SRVRecordName := fmt.Sprintf("%s.%s", serviceType, recordFQDN)

	client := cloudflare.NewClient(
		option.WithAPIToken(apiToken),
	)

	srvData := dns.SRVRecordDataParam{
		Target:   cloudflare.F(cleanSRV),
		Port:     cloudflare.F(float64(intPort)),
		Priority: cloudflare.F(float64(0)),
		Weight:   cloudflare.F(float64(0)),
	}

	payload := dns.SRVRecordParam{
		Type:    cloudflare.F(dns.SRVRecordTypeSRV),
		Name:    cloudflare.F(SRVRecordName),
		Proxied: cloudflare.F(false),
		Data:    cloudflare.F(srvData),
	}

	recordResponse, err := client.DNS.Records.Edit(
		c.Request.Context(),
		recordID,
		dns.RecordEditParams{
			ZoneID: cloudflare.F(zoneID),
			Body:   payload,
		},
	)

	if err == nil {
		c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Successfully changed the target domain to %s with the port %d", cleanSRV, intPort)})
		fmt.Printf("Updated Record: %+v\n", recordResponse)
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update the DNS record."})
		fmt.Printf("Cloudflare API Error: %+v\n", err)
	}
}

// regex was made by Google Gemini
var fqdnRegex = regexp.MustCompile(`^(?:[a-zA-Z0-9_](?:[a-zA-Z0-9-_]{0,61}[a-zA-Z0-9_])?\.)+[a-zA-Z0-9-]{2,}(?:\.)?$`)

func isValidFQDN(domain string) bool {

	if len(domain) < 4 || len(domain) > 253 {
		return false
	}

	return fqdnRegex.MatchString(domain)
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

	router.GET("/api/current/", checkAuth(), getIP)
	router.PATCH("/api/a/", checkAuth(), patchA)
	router.PATCH("/api/cname/", checkAuth(), patchCNAME)
	router.PATCH("/api/srv/", checkAuth(), patchSRV)

	router.Run("0.0.0.0:8080")
}
