package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"

	"os"

	"github.com/gin-gonic/gin"
)

type data struct {
	IP string `json:"ip"`
}

type CloudflareResponse struct {
	Result RecordData `json:"result"`
}

type RecordData struct {
	IP string `json:"content"`
}

type UpdateDNSRecord struct {
	Content string `json:"content"`
}

func getIP(c *gin.Context) {

	// curl https://api.cloudflare.com/client/v4/zones/$ZONE_ID/dns_records/$DNS_RECORD_ID \
	// -H "Authorization: Bearer $CLOUDFLARE_API_TOKEN"

	zoneID := os.Getenv("CF_ZONE_ID")
	apiToken := os.Getenv("CF_API_TOKEN")

	recordID := os.Getenv("CF_DNS_RECORD_ID")

	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records/%s", zoneID, recordID)

	req, err := http.NewRequest("GET", url, nil)

	req.Header.Add("Authorization", "Bearer "+apiToken)
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println("Error executing request:", err)
		return
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println("Error reading response body:", err)
		return
	}

	var responseData CloudflareResponse

	// Unmarshal converts the raw byte body into our map
	if err := json.Unmarshal(body, &responseData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println("Error parsing JSON:", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"content": responseData.Result.IP})
	fmt.Printf("Extracted Data: %#v\n", responseData.Result.IP)
}

func patchIP(c *gin.Context) {

	var newIP data
	if err := c.ShouldBindJSON(&newIP); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	candidateIP := newIP.IP
	cleanIP := strings.TrimSpace(candidateIP)

	ip := net.ParseIP(cleanIP)
	if ip == nil && ip.To4() == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Provided string is not a valid IPv4 address"})
		return
	}

	zoneID := os.Getenv("CF_ZONE_ID")
	apiToken := os.Getenv("CF_API_TOKEN")
	recordID := os.Getenv("CF_DNS_RECORD_ID")

	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records/%s", zoneID, recordID)

	payloadData := UpdateDNSRecord{
		Content: cleanIP,
	}

	jsonData, err := json.Marshal(payloadData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonData))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error executing request:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode == 200 {
		c.JSON(http.StatusOK, gin.H{"message": "Successfully changed the IP to " + cleanIP})
		fmt.Println("Successfully updated DNS record!")
		// Optional: You could unmarshal the response body here to confirm the new IP
	} else {
		c.JSON(resp.StatusCode, gin.H{"error": "Failed to update.", "response": string(body)})
		fmt.Printf("Failed to update. Status: %d, Response: %s\n", resp.StatusCode, string(body))
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
