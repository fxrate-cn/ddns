package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	ip := getIP()
	setDNS(ip)
}

func getIP() string {
	bodyBytes := request("GET", "http://checkip.amazonaws.com", nil)
	bodyString := string(bodyBytes)
	return bodyString
}

// use godot package to load/read the .env file and
// return the value of the key
func goDotEnvVariable(key string) string {
	// load .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	return os.Getenv(key)
}

type DNSRecords struct {
	Result []struct {
		ID        string `json:"id"`
		ZoneID    string `json:"zone_id"`
		ZoneName  string `json:"zone_name"`
		Name      string `json:"name"`
		Type      string `json:"type"`
		Content   string `json:"content"`
		Proxiable bool   `json:"proxiable"`
		Proxied   bool   `json:"proxied"`
		TTL       int    `json:"ttl"`
		Locked    bool   `json:"locked"`
		Meta      struct {
			AutoAdded           bool   `json:"auto_added"`
			ManagedByApps       bool   `json:"managed_by_apps"`
			ManagedByArgoTunnel bool   `json:"managed_by_argo_tunnel"`
			Source              string `json:"source"`
		} `json:"meta"`
		Comment    any       `json:"comment"`
		Tags       []any     `json:"tags"`
		CreatedOn  time.Time `json:"created_on"`
		ModifiedOn time.Time `json:"modified_on"`
		Priority   int       `json:"priority,omitempty"`
	} `json:"result"`
	Success    bool  `json:"success"`
	Errors     []any `json:"errors"`
	Messages   []any `json:"messages"`
	ResultInfo struct {
		Page       int `json:"page"`
		PerPage    int `json:"per_page"`
		Count      int `json:"count"`
		TotalCount int `json:"total_count"`
		TotalPages int `json:"total_pages"`
	} `json:"result_info"`
}

type Record struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Content string `json:"content"`
	Proxied bool   `json:"proxied"`
}

func setDNS(ip string) {
	zoneId := goDotEnvVariable("zoneId")
	name := goDotEnvVariable("name")
	baseUrl := "https://api.cloudflare.com/client/v4/zones/"
	dnsRecordsApi := fmt.Sprint(baseUrl, zoneId, "/dns_records")
	bodyBytes := request("GET", dnsRecordsApi, nil)
	var result DNSRecords
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		fmt.Println("Can not unmarshal JSON")
	}
	var id = ""
	// Loop through the data and find the record id of domain name
	for _, rec := range result.Result {
		if rec.Name == name {
			id = rec.ID
		}
	}
	updateApi := dnsRecordsApi + "/" + id
	record := &Record{Name: name, Content: ip, Proxied: true}
	jsonValue, _ := json.Marshal(record)
	updateResult := request("PATCH", updateApi, bytes.NewBuffer(jsonValue))
	log.Print(string(updateResult))
}

func request(method string, url string, reader io.Reader) []byte {
	apiToken := goDotEnvVariable("apiToken")
	req, err := http.NewRequest(method, url, reader)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+apiToken)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		return bodyBytes
	}
	return nil
}
