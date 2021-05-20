package sangoma

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strings"
	"net/url"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/inputs"
)

type Sangoma struct {
	DevicesIP  	[]string	`toml:"devices_ip"`
	DevicesName	[]string	`toml:"devices_name"`
	Username     	string		`toml:"username"`
	Password 	string		`toml:"password"`
}

var SangomaConfig = `
  ##Sample Config
  #devices_ip = ["0.0.0.0"]
  #devices_name = [""]
  #username = "user"
  #password = "password"
`

type myjar struct {
	jar map[string][]*http.Cookie
}

func (p *myjar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	p.jar[u.Host] = cookies
}

func (p *myjar) Cookies(u *url.URL) []*http.Cookie {
	return p.jar[u.Host]
}

func (s *Sangoma) SampleConfig() string {
	return SangomaConfig
}

func (s *Sangoma) Description() string {
	return "Collects basic metrics from Sangoma Vega 60G Device"
}

func (s *Sangoma) Gather(acc telegraf.Accumulator) error {

	for k, device := range s.DevicesIP {

		client := &http.Client{}

		jar := &myjar{}
		jar.jar = make(map[string][]*http.Cookie)
		client.Jar = jar

		/* Authenticate */
		payloadRqst := "username=" + s.Username + "&password=" + s.Password + "&last=1"
		var payload = []byte(payloadRqst)
		req, _ := http.NewRequest("POST", "http://"+device+"/vs_login", bytes.NewBuffer(payload))
		q := req.URL.Query()
		req.URL.RawQuery = q.Encode()

		resp, _ := client.Do(req)
		/*if err != nil {
			fmt.Printf("Error : %s", err)
		}*/

		/* Get Details */
		req, _ = http.NewRequest("GET", "http://"+device+"/vsconfig?sid=0&form_name=95&dont_need_uri_decode=1&cli_command=show%20ports", nil)
		resp, _ = client.Do(req)
		/*if err != nil {
			fmt.Printf("Error : %s", err)
		}*/

		bodyText, _ := ioutil.ReadAll(resp.Body)
		body := string(bodyText)

		fxo_results := strings.Split(body, "\n")

		fxo_total_channels := 0
		fxo_inuse_channels := 0
		fxo_available_channels := 0
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

		//Fields association
		fields := make(map[string]interface{})
        	fields[s.DevicesName[k] + ".fxo_total_channels"] = fxo_total_channels 
		fields[s.DevicesName[k] + ".fxo_inuse_channels"] = fxo_inuse_channels
		fields[s.DevicesName[k] + ".fxo_available_channels"] = fxo_available_channels
		fields[s.DevicesName[k] + ".fxo_disconnected_channels"] = fxo_disconnected_channels

		tags := make(map[string]string)

		acc.AddFields("sangoma", fields, tags)
	}
	return nil
}

func init() {
	inputs.Add("sangoma", func() telegraf.Input { return &Sangoma{} })
}
