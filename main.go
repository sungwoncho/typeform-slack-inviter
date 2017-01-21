package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type payload struct {
	Responses []response `json:"responses"`
}

type response struct {
	Answers map[string]string `json:"answers"`
}

var typeFormAPIKeyFlag = flag.String("typeformAPIKey", "", "TypeForm API key")
var formUIDFlag = flag.String("formUID", "", "TypeForm form UID")
var slackAPITokenFlag = flag.String("slackAPIToken", "", "SlackAPIToken")
var intervalFlag = flag.Int("interval", 1, "Interval for checking TypeForm responses (in minutes)")

// fetchTypeformData fetches the typeform responses since the given timestamp
// and returns a `payload`
func fetchTypeformData(typeFormAPIKey string, formUID string, since time.Time) (payload, error) {
	sinceTimestamp := since.Unix()
	endpoint := fmt.Sprintf("https://api.typeform.com/v1/form/%s?key=%s&since=%d", formUID, typeFormAPIKey, sinceTimestamp)
	fmt.Println(endpoint)
	var ret payload

	res, err := http.Get(endpoint)
	if err != nil {
		return ret, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return ret, err
	}

	var data payload
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("error unmarshalling")
		return ret, err
	}

	return data, nil
}

// sendSlackInvitation sends Slack invitation to the given email address
func sendSlackInvitation(email string) error {
	fmt.Println("Inviting", email)
	endpoint := fmt.Sprintf("https://slack.com/api/users.admin.invite?token=%s&email=%s&resend=false", *slackAPITokenFlag, email)

	req, err := http.NewRequest("POST", endpoint, nil)
	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response body:", string(body))
	return nil
}

func run() {
	now := time.Now().UTC()
	targetTime := now.Add(-time.Duration(*intervalFlag) * time.Minute)

	data, err := fetchTypeformData(*typeFormAPIKeyFlag, *formUIDFlag, targetTime)
	if err != nil {
		panic(err)
	}

	fmt.Println("Need to invite", len(data.Responses))

	for _, response := range data.Responses {
		email := response.Answers["email_40622900"]

		if email == "" {
			email = response.Answers["email_41266127"]
		}

		if email != "" {
			sendSlackInvitation(email)
		}
	}
}

func main() {
	flag.Parse()

	t := time.NewTicker(time.Duration(*intervalFlag) * time.Minute)

	for {
		run()
		<-t.C
	}
}
