package main

import (
	"time"
	"strconv"
	"flag"
	"fmt"
	"os"
	"io"
	"text/template"
	"net/url"
	"net/http"
	"strings"
)

var t *template.Template

type card struct {
	ID string			// Device IP/Hostname
	Activity string		// HTML formatted brief message
	Title string		// Card title (expanded view)
	Description string	// HTML formatted (expanded view)
}

type message struct {
	Text string		// unformatted message text
	Color string	// red, green, yellow
	Card card
}

type postMsg struct {
	URL *url.URL
	Auth *string
	Body io.Reader
}

func httpPost(msg postMsg) (http.Response, error) {
	req, err := http.NewRequest("POST",msg.URL.String(),msg.Body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer " + *msg.Auth)

	if err != nil {
		fmt.Println("Unable to create http request")
		os.Exit(1)
	}

	fmt.Printf("Connecting to: %s\n", msg.URL.String())

	timeout := time.Duration(2 * time.Second)
	client := &http.Client{Timeout: timeout}
	res, err := client.Do(req)
	
	if err == nil {
		defer res.Body.Close()
	}

	return *res, err
}

// createHcNotifyURL Generates HipChat server notification URL
// for testing authentication or submitting JSON message data
func createHcNotifyURL(serverURL *string, roomID *int, test bool) (*url.URL, error) {
	out := *serverURL + "/v2/room/" + strconv.Itoa(*roomID) + "/notification"

	// Create test url
	// https://developer.atlassian.com/server/hipchat/about-the-hipchat-rest-api/
	if test {
		out = out + "?auth_test=true"
	}
	
	return url.Parse(out)
}

func main() {

	// cli params
	// General options
	templateFile := flag.String("template", "./message.tmpl", "Location of message template file")

	// HipChat connection parameters
	authTest := flag.Bool("test", false, "Does authenticaiton test to HipChat server. Must provide -server -room and -token")
	hcSrvURL := flag.String("server", "", "HipChat Server address (Required)")
	hcSrvRoomID := flag.Int("room", 0, "HipChat Room ID (Required)")
	hcSrvAuthToken := flag.String("token", "", "HipChat Room Authorization token (Required)")

	msgSeverity := flag.String("severity", "normal", "Message type <normal|warning|critical> (Required)")
	clientAddress := flag.String("client", "", "SNMP Trap device address in form of IP or Hostname")
	trapMessage := flag.String("msg", "", "SNMP Trap message")
	trapDescription := flag.String("desc", "", "SNMP Trap description (optional)")
	trapOID := flag.String("oid", "", "SNMP Trap OID value")

	flag.Parse()

	// validate provided parameters
	if flag.Parsed() {
		// required fields
		if *hcSrvURL == "" || *hcSrvRoomID == 0 || *hcSrvAuthToken == "" {
			flag.PrintDefaults()
			os.Exit(1)
		}

		// validate things when not doing auth test
		if !*authTest {
			// validate message type selection
			msgSeverityChoices := map[string]bool{"normal":true, "warning":true, "critical":true}
			if _, validChoice := msgSeverityChoices[strings.ToLower(*msgSeverity)]; !validChoice {
				flag.PrintDefaults()
				os.Exit(1)
			}

			// check if template file provided
			if _, err := os.Stat(*templateFile); os.IsNotExist(err) {
				fmt.Printf("Unable to find message template file in: %s\n", *templateFile)
				os.Exit(1)
			}
		}
	}

	// Test authentication to HipChat Server Room
	// Authentication token must have notify access
	if *authTest {
		url, _ := createHcNotifyURL(hcSrvURL,hcSrvRoomID,true)
		msg := postMsg{url,hcSrvAuthToken,nil}
		res, err := httpPost(msg)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if res.StatusCode != 202 {
			fmt.Printf("Authentication test failed. (%s)\n", res.Status)
		} else {
			fmt.Printf("Authentication test passed! (%s)\n", res.Status)
		}
	}

	if !*authTest {
		// create template
		hcMsgText := "SNMP Trap received from " + *clientAddress
		hcMsgHTML := "<b>" + *clientAddress + "</b> - " + *trapMessage
		hcTitle := *trapOID
		
		// define message color
		msgColor := "yellow"
		switch strings.ToLower(*msgSeverity) {
			case "normal":
				msgColor = "green"
			case "warning":
				msgColor = "yellow"
			case "critical":
				msgColor = "red"
		}

		c1 := card{*clientAddress,hcMsgHTML,hcTitle,*trapDescription}
		m1 := message{hcMsgText,msgColor,c1}

		t,err := template.ParseFiles(*templateFile)
		if err != nil {
			fmt.Println(err)
		}

		// Create IO Pipe to communicate between http and tempate
		read, write := io.Pipe()

		// Writing without a reader will deadlock so write in a goroutine
		go func() {
			defer write.Close()
			err = t.ExecuteTemplate(write,"message",m1)
			if err != nil {
				fmt.Println(err)
			}
		}()

		// Build HipChat message
		notifyURL, _ := createHcNotifyURL(hcSrvURL,hcSrvRoomID,false)
		notifyMsg := postMsg{notifyURL,hcSrvAuthToken,read}

		res,err := httpPost(notifyMsg)

		if err != nil {
			fmt.Println("Something went wrong...")
		}

		fmt.Printf("Got response: %s\n", res.Status)
	}
}
