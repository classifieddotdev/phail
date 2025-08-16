// goblin.go
package main

import (
	"crypto/tls"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

type Target struct {
	Name         string `json:"name"`
	Target       string `json:"target"`
	PrimaryIP    string `json:"primary_ip"`
	FailoverIP   string `json:"failover_ip"`
	DNSName      string `json:"dns_name"`
	PingInterval int    `json:"ping_interval"`
	IgnoreSSL    bool   `json:"ignoreSSL"`
	Retries      int    `json:"retries"`
	RetryInterval int   `json:"retry_interval"`
}

type Config struct {
	Targets []Target `json:"targets"`
}

func main() {
	data, err := os.ReadFile("config.json")
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		log.Fatalf("Error parsing config: %v", err)
	}

	for _, t := range cfg.Targets {
		target := t
		go monitorTarget(target)
	}

	select {} // this stops go from exiting
}

func monitorTarget(t Target) {
	client := &http.Client{
		Timeout: 3 * time.Second,
	}
	if t.IgnoreSSL {
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}

	for {
		healthy := false
		for i := 0; i < t.Retries; i++ {
			resp, err := client.Get(t.Target)
			if err == nil && resp.StatusCode == 200 {
				healthy = true
				break
			}
			time.Sleep(time.Duration(t.RetryInterval) * time.Second)
		}

		var currentIP string
		if healthy {
			currentIP = t.PrimaryIP
		} else {
			currentIP = t.FailoverIP
		}

		if err := updateDNSRecord(t.DNSName, currentIP, healthy); err != nil {
			log.Printf("[%s] Error updating DNS: %v", t.DNSName, err)
		}

		time.Sleep(time.Duration(t.PingInterval) * time.Second)
	}
}