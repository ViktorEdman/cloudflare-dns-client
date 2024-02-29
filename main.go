package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	zoneid string
	token  string
	domain string
)

func main() {
	zoneflag := flag.String("zone", "", "Specifies zone id with cloudflare")
	tokenflag := flag.String("token", "", "Specifies API token with cloudflare")
	domainflag := flag.String("domain", "", "Specifies domain name with cloudflare")
	flag.Parse()
	flag.VisitAll(func(f *flag.Flag) {
		if f.Value.String() == "" {
			fmt.Println(f.Name + " not set!!!")
			flag.Usage()
			os.Exit(1)
		}
	})
	zoneid = *zoneflag
	token = *tokenflag
	domain = *domainflag
	Run(zoneid, token, domain)
}

func Run(zoneid string, token string, domain string) {
	url := "https://api.cloudflare.com/client/v4/zones/" + zoneid + "/dns_records"
	fmt.Println(strings.Repeat("-", 26))
	fmt.Println(time.Now())
	var externalIp, dnsIp, entryId string
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		externalIp = retrieveExternalIp()
	}()
	go func() {
		defer wg.Done()
		dnsIp, entryId = retrieveDnsIp(url, token)
	}()

	wg.Wait()
	fmt.Println("External IP:", externalIp)
	if len(dnsIp) == 0 {
		log.Fatal("Failed to fetch dns ip!")
	}
	if dnsIp == externalIp {
		fmt.Println("Match! No update required")
		fmt.Println(strings.Repeat("-", 26))
		os.Exit(0)
	}
	updateDnsIp(token, url, entryId, externalIp)
	fmt.Println(strings.Repeat("-", 26))
}

func updateDnsIp(token string, url string, entryId string, IP string) {
	URL := url + "/" + entryId
	payload := map[string]string{"content": IP}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}
	req, _ := http.NewRequest("PATCH", URL, bytes.NewBuffer(jsonPayload))
	req.Header.Add("Authorization", "Bearer "+token)
	res, resErr := http.DefaultClient.Do(req)
	if resErr != nil {
		log.Fatal(resErr)
	}
	body, _ := io.ReadAll(res.Body)

	var patchResponse PatchResponse
	marshalError := json.Unmarshal(body, &patchResponse)
	if marshalError != nil {
		log.Fatal(marshalError)
	}
	fmt.Println("IP updated to", patchResponse.Result.Content)
}

func retrieveExternalIp() string {
	ipCheckRes, _ := http.Get("https://checkip.amazonaws.com")
	ipCheckBody, _ := io.ReadAll(ipCheckRes.Body)
	externalIp := strings.TrimSpace(string(ipCheckBody))
	return externalIp
}

func retrieveDnsIp(url string, token string) (string, string) {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", "Bearer "+token)
	res, responseError := http.DefaultClient.Do(req)

	if responseError != nil {
		log.Fatal("Request failed ", responseError)
	}
	if res.StatusCode == http.StatusUnauthorized {
		log.Fatal("Authentication with cloudflare failed! Check your token")
	}

	body, bodyError := io.ReadAll(res.Body)

	if bodyError != nil {
		log.Fatal("Could not read response body ", bodyError)
	}

	var cloudflareRes GetResponse
	jsonError := json.Unmarshal(body, &cloudflareRes)
	if jsonError != nil {
		log.Fatal("Could not parse JSON ", jsonError)
	}
	defer res.Body.Close()
	var dnsIp string
	var entryId string
	for i := 0; i < len(cloudflareRes.Result); i++ {
		result := cloudflareRes.Result[i]
		if result.Name == domain && result.Type == "A" {
			fmt.Println(
				"Dns IP:",
				result.Content,
			)
			dnsIp = strings.TrimSpace(result.Content)
			entryId = strings.TrimSpace(result.ID)

		}
	}
	return dnsIp, entryId
}
