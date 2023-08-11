package cloudfare

import (
	"context"
	"fmt"
	"jarvis-trading-bot/consts"
	"log"
	"os"

	"github.com/cloudflare/cloudflare-go"
)

func UpdateDNS(cnameValue string) {
	dnsZone := os.Getenv(consts.DnsZone)
	dnsRecord := os.Getenv(consts.DnsRecord)
	//apiToken := os.Getenv(consts.CloudfareApiToken)
	apiKey := os.Getenv(consts.CloudfareApiKey)
	apiEmail := os.Getenv(consts.CloudfareApiEmail)

	if dnsZone == "" || dnsRecord == "" {
		log.Fatal("DNS zone or DNS record is not set")
	}
	
	if apiKey == "" || apiEmail == "" {
		log.Fatal("API Key or API Email is not set")
	}

	//api, err := cloudflare.NewWithAPIToken(apiToken)
	api, err := cloudflare.New(apiKey, apiEmail)
	if err != nil {
		log.Fatal(err)
	}

	// Most API calls require a Context
	ctx := context.Background()

	// Fetch user details on the account
	u, err := api.UserDetails(ctx)
	if err != nil {
		log.Fatal(err)
	}
	// Print user details
	fmt.Println(u)

	// Fetch the zone ID
	zoneID, err := api.ZoneIDByName(dnsZone) // Assuming example.com exists in your Cloudflare account already
	if err != nil {
		log.Fatal(err)
	}

	zone, err := api.ZoneDetails(ctx, zoneID)
	if err != nil {
		log.Fatal(err)
	}

	// Print zone details
	fmt.Println(zone)

	// Fetch all DNS records for example.org
	records, _, err := api.ListDNSRecords(ctx, cloudflare.ZoneIdentifier(zoneID), cloudflare.ListDNSRecordsParams{})
	if err != nil {
		fmt.Println(err)
		return
	}
    
    idxRecord := -1
    for i, record := range records {
        if record.Name == dnsRecord {
            idxRecord = i
        }
    }
	// Print DNS records
	proxied := false
	
	if idxRecord == -1 {
	    fmt.Println("Creating DNS record....")
	    _, err = api.CreateDNSRecord(ctx, cloudflare.ZoneIdentifier(zoneID), cloudflare.CreateDNSRecordParams{
    		Type:    "CNAME",
    		Name:    dnsRecord,
    		Content: cnameValue,
    		Proxied: &proxied,
    		TTL:     3600,
    	})
	} else {
	    fmt.Println("Updating DNS record...")
	    err = api.UpdateDNSRecord(ctx, cloudflare.ZoneIdentifier(zoneID), cloudflare.UpdateDNSRecordParams{
	    	ID:      records[idxRecord].ID,
    		Type:    "CNAME",
    		Name:    dnsRecord,
    		Content: cnameValue,
    		Proxied: &proxied,
    		TTL:     3600,
    	})
	}
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Updated record %s to %s", dnsRecord, cnameValue)
}
