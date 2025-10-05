// goblin.go
package main

import (
	"crypto/tls"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
	"github.com/containrrr/shoutrrr"
	"github.com/containrrr/shoutrrr/pkg/router"
)

func printBanner() {
	log.Println(`
	██████╗ ██╗  ██╗ █████╗ ██╗██╗     
	██╔══██╗██║  ██║██╔══██╗██║██║     
	██████╔╝███████║███████║██║██║     
	██╔═══╝ ██╔══██║██╔══██║██║██║     
	██║     ██║  ██║██║  ██║██║███████╗
	╚═╝     ╚═╝  ╚═╝╚═╝  ╚═╝╚═╝╚══════╝
	
	P o o r - m a n ’ s   H i g h l y  
	  A v a i l a b l e   I n t e r n e t  
	     L o a d  B a l a n c e r
	
	Started at: ` + time.Now().UTC().Format(time.RFC3339) + "\n")
}

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
	MonitorOnly bool   `json:"monitor_only"`
}

type Config struct {
	Targets         []Target `json:"targets"`
	ShoutrrrURLs    []string `json:"shoutrrr_urls"`
}

func main() {
	log.SetOutput(os.Stdout)
	printBanner()
	data, err := os.ReadFile("config.json")
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		log.Fatalf("Error parsing config: %v", err)
	}
	var notifier *router.ServiceRouter
	if len(cfg.ShoutrrrURLs) > 0 {
		notifier, err = shoutrrr.CreateSender(cfg.ShoutrrrURLs...)
		if err != nil {
			log.Fatalf("Failed to create shoutrrr sender: %v", err)
		}
	}


	for _, t := range cfg.Targets {
		target := t
		go monitorTarget(target, notifier)
	}

	select {} // this stops go from exiting
}

func monitorTarget(t Target, notifier *router.ServiceRouter) {
	client := &http.Client{
		Timeout: 3 * time.Second,
	}
	if t.IgnoreSSL {
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}

	var lastHealthy bool
	var initialized bool

	for {
		healthy := false
		for i := 0; i < t.Retries; i++ {
			resp, err := client.Get(t.Target)
			if err == nil {
				resp.Body.Close()
				if resp.StatusCode == 200 {
					healthy = true
					break
				}
			}
			time.Sleep(time.Duration(t.RetryInterval) * time.Second)
		}

		var currentIP string
		if healthy {
			currentIP = t.PrimaryIP
		} else {
			currentIP = t.FailoverIP
		}

		if t.MonitorOnly {
			if !initialized || lastHealthy != healthy {
				notify(t.DNSName, healthy, notifier)
				lastHealthy = healthy
				initialized = true
			}
		} else {
			if err := updateDNSRecord(t.DNSName, currentIP, healthy, notifier); err != nil {
				log.Printf("[%s] Error updating DNS: %v", t.DNSName, err)
			}
		}

		

		time.Sleep(time.Duration(t.PingInterval) * time.Second)
	}
}