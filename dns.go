// dns.go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/cloudflare/cloudflare-go/v5"
	"github.com/cloudflare/cloudflare-go/v5/dns"
	"github.com/cloudflare/cloudflare-go/v5/option"
	"github.com/containrrr/shoutrrr/pkg/router"
	"github.com/containrrr/shoutrrr/pkg/types"
)

var (
	cfZoneID    = os.Getenv("CF_ZONE_ID")
	cfAuthToken = os.Getenv("CF_API_TOKEN")
)

type DNSResponse struct {
	Result []DNSRecord `json:"result"`
}

type DNSRecord struct {
	Name    string `json:"name"`
	TTL     int    `json:"ttl"`
	Type    string `json:"type"`
	Content string `json:"content"`
	Proxied bool   `json:"proxied"`
	ID      string `json:"id"`
}

func updateDNSRecord(name, ip string, healthy bool, notifier *router.ServiceRouter) error {
	//first get the record id
	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records?name=%s", cfZoneID, name)
	req, err := http.NewRequest(
		"GET",
		url,
		nil,
	)
	if (err != nil) {
		return err
	}
	
	req.Header.Set("Authorization", "Bearer "+cfAuthToken)
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	body, _ := io.ReadAll(resp.Body)
	
	var result DNSResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return(err)
	}
	
	dnsRecord := result.Result[0]
	
	if dnsRecord.Content == ip {
		// no update needed
		return nil
	}
	
	if healthy {
		log.Printf("[%s] Target is healthy. Pointing DNS to primary IP", name)
	} else {
		log.Printf("[%s] Target is unhealthy. Pointing DNS to failover IP", name)
	}
	
	
	// do DNS update if the current and new ip don't match
	client := cloudflare.NewClient(
		option.WithAPIToken(cfAuthToken),
	)
	timestamp := time.Now().UTC().Format(time.RFC3339)
	prefix := "[FAILOVER]"
	if healthy {
	    prefix = "[HEALTHY]"
	}
	message := fmt.Sprintf("%s Updated by PHAIL at %s", prefix, timestamp)
	_, err = client.DNS.Records.Edit(
		context.TODO(),
		dnsRecord.ID,
		dns.RecordEditParams{
			ZoneID: cloudflare.F(cfZoneID),
			Body: dns.ARecordParam{
				Name: cloudflare.F(dnsRecord.Name),
				TTL: cloudflare.F(dns.TTL1),
				Type: cloudflare.F(dns.ARecordTypeA),
				Content: cloudflare.F(ip),
				Comment: cloudflare.F(message),
				Proxied: cloudflare.F(dnsRecord.Proxied),
			},
		},
	)
	
	if err != nil {
		return err
	}

	//send shoutrrr
	params := &types.Params{
    	"title": name,
	}
	if notifier != nil {
		notifier.Send(message, params)
	}
	
	return nil
}