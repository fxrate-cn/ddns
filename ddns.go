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
	fmt.Println(ip)
	setDNS(ip)
}

func getIP() string {
	resp, err := http.Get("http://checkip.amazonaws.com")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		bodyString := string(bodyBytes)
		log.Println(bodyString)
		return bodyString
	} else {
		return ""
	}
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
	apiToken := goDotEnvVariable("apiToken")
	name := goDotEnvVariable("name")
	baseUrl := "https://api.cloudflare.com/client/v4/zones/"
	dnsRecordsApi := fmt.Sprint(baseUrl, zoneId, "/dns_records")
	req, err := http.NewRequest("GET", dnsRecordsApi, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+apiToken)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		var result DNSRecords
		if err := json.Unmarshal(bodyBytes, &result); err != nil { // Parse []byte to the go struct pointer
			fmt.Println("Can not unmarshal JSON")
		}

		var id = ""
		// Loop through the data node for the FirstName
		for _, rec := range result.Result {
			if rec.Name == name {
				id = rec.ID
			}
		}
		updateApi := dnsRecordsApi + "/" + id
		record := Record{}
		record.Name = name
		record.Content = ip
		record.Proxied = true
		jsonValue, _ := json.Marshal(record)
		req, err := http.NewRequest("PATCH", updateApi, bytes.NewBuffer(jsonValue))
		if err != nil {
			log.Fatal(err)
		}
		req.Header.Set("Authorization", "Bearer "+apiToken)
		req.Header.Set("Content-Type", "application/json")
		resp, err := client.Do(req)
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(string(bodyBytes))
		}

	}
}
