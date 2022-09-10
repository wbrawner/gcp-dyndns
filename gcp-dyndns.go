package main

import (
	"context"
	_ "embed"
	"flag"
	"log"
	"net/http"
	"strings"

	"google.golang.org/api/dns/v1"
	"google.golang.org/api/option"
)

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

func updateRecord(project string, zone string, domainName string, ipAddr string, credentialsFile string) error {
	ctx := context.Background()
	client, err := dns.NewService(ctx, option.WithCredentialsFile(credentialsFile))
	if err != nil {
		return err
	}
	deletion, err := client.ResourceRecordSets.Get(project, zone, domainName, "A").Do()
	if err != nil {
		return err
	}
	if !strings.HasSuffix(domainName, ".") {
		domainName += "."
	}
	addition := dns.ResourceRecordSet{
		Name:    domainName,
		Rrdatas: []string{ipAddr},
		Ttl:     60,
		Type:    "A",
	}
	change := dns.Change{
		Additions: []*dns.ResourceRecordSet{
			&addition,
		},
		Deletions: []*dns.ResourceRecordSet{
			deletion,
		},
	}
	_, err = client.Changes.Create(project, zone, &change).Do()
	return err
}

func main() {
	project := flag.String("project", "", "GCP project name")
	zone := flag.String("zone", "", "GCP zone name")
	domain := flag.String("domain", "", "domain name to update")
	credentials := flag.String("credentials", "", "path to credentials JSON file")
	flag.Parse()
	ip, err := currentIp()
	if err != nil {
		log.Fatalf("failed to get current ip address: %v", err)
	}
	err = updateRecord(*project, *zone, *domain, ip, *credentials)
	if err != nil {
		log.Fatalf("failed to update record: %v", err)
	}
}
