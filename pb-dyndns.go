package main

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

type editRecordRequest struct {
	Secret  string `json:"secretapikey"`
	Key     string `json:"apikey"`
	Content string `json:"content"`
	Ttl     int    `json:"ttl"`
}

type editRecordResponse struct {
	Status string `json:"status"`
}

func currentIp() (string, error) {
	res, err := http.DefaultClient.Get("https://ip.wbrawner.com")
	if err != nil {
		return "", err
	}
	ipBytes := make([]byte, 15)
	_, err = res.Body.Read(ipBytes)
	if err != nil {
		return "", err
	}
	return string(ipBytes), nil
}

func updateRecord(domain string, subdomain string, recordType string, request editRecordRequest) error {
	url := fmt.Sprintf("https://porkbun.com/api/json/v3/dns/editByNameType/%s/%s/%s", domain, recordType, subdomain)
	body, err := json.Marshal(request)
	if err != nil {
		return err
	}
	bodyReader := bytes.NewReader(body)
	res, err := http.DefaultClient.Post(url, "application/json", bodyReader)
	if err != nil {
		return err
	}
	var response editRecordResponse
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return err
	}
	if response.Status != "SUCCESS" {
		return fmt.Errorf("failed to update record: %v", response.Status)
	}
	return nil
}

func readCredentials(path string) (editRecordRequest, error) {
	contents, _ := os.ReadFile(path)
	var request editRecordRequest
	err := json.Unmarshal(contents, &request)
	return request, err
}

func main() {
	domain := flag.String("domain", "", "domain name to update")
	recordType := flag.String("type", "A", "record type to update (e.g. A, CNAME, etc)")
	subdomain := flag.String("subdomain", "", "subdomain to update")
	credentialsFile := flag.String("credentials", "credentials.json", "path to credentials JSON file")
	flag.Parse()
	request, err := readCredentials(*credentialsFile)
	if err != nil {
		log.Fatalf("failed to read credentials JSON: %v", err)
	}
	ip, err := currentIp()
	if err != nil {
		log.Fatalf("failed to get current ip address: %v", err)
	}
	request.Content = ip
	request.Ttl = 600
	err = updateRecord(*domain, *subdomain, *recordType, request)
	if err != nil {
		log.Fatalf("failed to update record: %v", err)
	}
}
