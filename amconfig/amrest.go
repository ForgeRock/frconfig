package amconfig

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	log "github.com/Sirupsen/logrus"
	"net/http"
	"github.com/spf13/viper"
	"io"
)

// OpenAMConnection to an openam server instance
type OpenAMConnection struct {
	BaseURL  string // base URL including /openam. Example: http://openam.example.com:8080/openam
	User     string
	Password string
	tokenId  string
	Realm    string
}

// AuthNResponse returned by OpenAM on authenticate request
type AuthNResponse struct {
	TokenID   string `json: "tokenId"`
	SuccessURL string `json:"successUrl"`
}

// Create an OpenAM connection based on viper config file
func GetOpenAMConnection() (am *OpenAMConnection,err error) {
	url := viper.GetString("default.openam.url")
	pass := viper.GetString("default.openam.password")
	user := viper.GetString("default.openam.user")
	return Open(url,user,pass)
}

func Open(url, user, password string) (am *OpenAMConnection,err error) {
	am = &OpenAMConnection{BaseURL:url, User: user, Password: password}
	err = am.Authenticate()

	return
}

// Authenticate to OpenAM. Set the AuthN token in the connection for subsequent requests
func (am *OpenAMConnection)Authenticate() error {

	// get session token

	url := am.requestURL("/json/authenticate")

	var jsonStr = []byte(`{}`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))

	req.Header.Set("X-OpenAM-Username", am.User)
	req.Header.Set("X-OpenAM-Password", am.Password)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	log.Debug("Authenticating to OpenAM ", url)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	var a AuthNResponse

	err = json.Unmarshal(body, &a)

	if resp.StatusCode != 200 {
		return fmt.Errorf("Failed to Authenticate. %v", resp.Status)

	}


	am.tokenId = a.TokenID

	return nil
}


// create a the request URL given the partial path
func (am *OpenAMConnection) requestURL(path string) string {
	return fmt.Sprint(am.BaseURL, path)
}

func debug(data []byte, err error) {
	if err == nil {
		fmt.Printf("%s\n\n", data)
	} else {
		log.Fatalf("%s\n\n", err)
	}
}

// Create a new request setting the iPro auth cookie and the content type
func (openam *OpenAMConnection)newRequest(method, url string, body io.Reader) *http.Request {
	//client := &http.Client{}


	req, err := http.NewRequest(method, openam.requestURL(url), body)
	if err != nil {
		log.Panicf("Could not create new request, err = %v", err)
	}
	ipro := http.Cookie{Name: "iPlanetDirectoryPro", Value: openam.tokenId}
	req.AddCookie(&ipro)
	req.Header.Set("Content-Type", "application/json")
	return req
}



/*

http://openam.forgerock.org/doc/bootstrap/dev-guide/index.html

Read a specific resource ListResourceTypes
curl \
--header "iPlanetDirectoryPro: AQIC5..." \
https://openam.example.com:8443/openam/json/myrealm/resourcetypes/12345a67-8f0b-123c-45de-6fab78cd01e3

Create a resouce type

curl \
--header "iPlanetDirectoryPro: AQIC5..." \
--request POST \
--header "Content-Type: application/json" \
--data '{
    "name": "My Resource Type",
    "actions": {
        "LEFT": true,
        "RIGHT": true,
        "UP": true,
        "DOWN": true
    },
    "patterns": [
        "http://device/location/*"
    ]
}'


update a resource type
curl \
--header "iPlanetDirectoryPro: AQIC5..." \
--request PUT \
--header "Content-Type: application/json" \
--data '{
    "uuid": "12345a67-8f0b-123c-45de-6fab78cd01e4",
    "name": "My Updated Resource Type",
    "actions": {
        "LEFT": false,
        "RIGHT": false,
        "UP": false,
        "DOWN": false
    },
    "patterns": [
        "http://device/location/*"
    ]
}' \

Delete

curl \
--request DELETE \
--header "iPlanetDirectoryPro: AQIC5..." \
https://openam.example.com:8443/openam/json/myrealm/resourcetypes/12345a67-8f0b-123c-45de-6fab78cd01e4
{}


*/