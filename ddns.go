package main
import (
	"io"
	"fmt"
	"log"
	"net/http"
)

func main() {
	ip:= getIP()
	fmt.Println(ip)
	setDNS()
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

func setDNS()  {
	zoneId := ""
	//accountId := ""
	//recordId := ""
	apiToken  := ""
	//spfRecordId := ""
	//apiKey := ""
	//userEmail := ""
	//recordName := ""
	baseUrl := "https://api.cloudflare.com/client/v4/zones/"
	dnsRecordsApi := fmt.Sprint(baseUrl, zoneId, "/dns_records")
	fmt.Println(dnsRecordsApi)
	// req, err := http.NewRequest("GET", "https://api.cloudflare.com/client/v4/user/tokens/verify", nil)
	req, err := http.NewRequest("GET", dnsRecordsApi, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer " + apiToken)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		bodyString := string(bodyBytes)
		log.Println("dns", bodyString)
		fmt.Println("OK")
	}
}
