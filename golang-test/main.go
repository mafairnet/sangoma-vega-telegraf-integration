package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

var vegaIP = "172.19.0.48"

func main() {
	client := &http.Client{}

	jar := &myjar{}
	jar.jar = make(map[string][]*http.Cookie)
	client.Jar = jar

	/* Authenticate */
	payloadRqst := "username=cti_dev&password=ct1_r34d3r&last=1"
	var payload = []byte(payloadRqst)
	req, err := http.NewRequest("POST", "http://"+vegaIP+"/vs_login", bytes.NewBuffer(payload))
	//req.SetBasicAuth("<username>", "<password>")
	q := req.URL.Query()
	/*q.Add("username", "cti_dev")
	q.Add("password", "ct1_r34d3r")*/
	req.URL.RawQuery = q.Encode()
	fmt.Printf("URL RawQuery %v\n", req.URL.RawQuery)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error : %s", err)
	}

	fmt.Printf("LOGIN Resp: %v\n", resp)

	/* Get Details */
	req, err = http.NewRequest("GET", "http://"+vegaIP+"/vsconfig?sid=0&form_name=95&dont_need_uri_decode=1&cli_command=show%20ports", nil)
	resp, err = client.Do(req)
	if err != nil {
		fmt.Printf("Error : %s", err)
	}

	bodyText, err := ioutil.ReadAll(resp.Body)
	s := string(bodyText)

	fmt.Printf("GET_CONFIG Resp: %v\n", s)

	fxo_results := strings.Split(s, "\n")

	fxo_total_channels := 0
	fxo_inuse_channels := 0
	fxo_available_channels := 0
	//fxo_blocked_channels := 0
	fxo_disconnected_channels := 0

	for _, channel := range fxo_results {
		if strings.Contains(channel, "ready") || strings.Contains(channel, "busy") || strings.Contains(channel, "offline") {
			//chan_data := strings.Split(channel, "")
			//print(chan_data)
			if len(channel) > 2 {
				if strings.Contains(channel, "ready") {
					fxo_available_channels = fxo_available_channels + 1
				}
				if strings.Contains(channel, "busy") {
					fxo_inuse_channels = fxo_inuse_channels + 1
				}
				if strings.Contains(channel, "offline") {
					fxo_disconnected_channels = fxo_disconnected_channels + 1
				}
				fxo_total_channels += 1
			}
		}
	}

	fmt.Printf("Total: %v\n", fxo_total_channels)
	fmt.Printf("InUse: %v\n", fxo_inuse_channels)
	fmt.Printf("Available: %v\n", fxo_available_channels)
	fmt.Printf("Disconnected: %v\n", fxo_disconnected_channels)
}
