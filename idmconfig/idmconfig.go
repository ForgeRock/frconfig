package idmconfig




import (
	//"bytes"
	//"encoding/json"
	"fmt"
	"io/ioutil"
	log "github.com/Sirupsen/logrus"
	"net/http"
	"errors"
	"bytes"
	"net/http/httputil"
	"github.com/spf13/viper"
	"github.com/forgerock/frconfig/crest"
)


const (
	IDM_CONFIGURATION = "idm.conf"
)


func init() {
	crest.RegisterCreateObjectHandler( []string{IDM_CONFIGURATION}, CreateObjects)

}


type IDMConnection struct {
	BaseURL string
	User string
	Password string
	httpClient *http.Client
	session_jwt *http.Cookie
}


// Create an OpenAM connection based on viper config file
func GetOpenIDMConnection() (idm *IDMConnection,err error) {
	url := viper.GetString("default.openidm.url")
	pass := viper.GetString("default.openidm.password")
	user := viper.GetString("default.openidm.user")
	idm = &IDMConnection{BaseURL:url, User:user, Password:pass}
	err = idm.Authenticate()
	return
}


// Authenticate to OpenIDM, and get the session jwt for subsequent REST operations
func (idm *IDMConnection) Authenticate() error {
	// get session token

	url := idm.requestURL("/info/login")

	req, err := http.NewRequest("GET", url, nil)

	req.Header.Set("X-OpenIDM-Username", idm.User)
	req.Header.Set("X-OpenIDM-Password", idm.Password)
	req.Header.Set("X-OpenIDM-NoSession", "false")

	req.Header.Set("Content-Type", "application/json")

	idm.httpClient = &http.Client{}
	resp, err := idm.httpClient.Do(req)
	if err != nil  {
		return err
	}
	defer resp.Body.Close()

	// not clear we need the response body...
	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return errors.New( fmt.Sprintf("status = %v msg=%v", resp.StatusCode, string(body)))
	}


	for _,cookie := range resp.Cookies() {
		if cookie.Name == "session-jwt" {
			idm.session_jwt = cookie
		}
	}

	if idm.session_jwt == nil {
		msg := fmt.Sprintf("No session cookie found in response Body: %v", string(body))
		log.Error(msg)
		return errors.New(msg)
	}
	return nil
}

func  (idm *IDMConnection) setHeaders(req *http.Request) {

	req.Header.Set("X-OpenIDM-Username", idm.User)
	req.Header.Set("X-OpenIDM-Password", idm.Password)
	req.Header.Set("X-OpenIDM-NoSession", "false")

	req.Header.Set("Content-Type", "application/json")

}

func (idm *IDMConnection) GET(path string) (string,error) {
	req, err := http.NewRequest("GET", idm.requestURL(path), nil)

	req.AddCookie(idm.session_jwt)
	resp, err := idm.httpClient.Do(req)
	if err != nil  {
		return "",err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	return string(body),nil

}


func (idm *IDMConnection) POST(path string, body []byte) (string,error) {
	req, err := http.NewRequest(http.MethodPost, idm.requestURL(path), bytes.NewBuffer(body))

	if err != nil {
		return "",err
	}
	//req.AddCookie(idm.session_jwt)

	idm.setHeaders(req)

	dump, err := httputil.DumpRequestOut(req, true)
	log.Infof("Request %s", dump)
	if err != nil {
		log.Infof("Request %s", dump)
	} else {
		log.Infof("error %v", err)
	}
	resp, err := idm.httpClient.Do(req)
	if err != nil  {
		log.Error(err)
		return "",err
	}
	defer resp.Body.Close()

	b, _ := ioutil.ReadAll(resp.Body)
	log.Infof("resp = %s status = %v", b, resp.StatusCode)
	return string(b),nil

}

// create a the request URL given the partial path
func (idm *IDMConnection) requestURL(path string) string {
	return fmt.Sprint(idm.BaseURL, path)
}

func  (idm *IDMConnection) GetConfig(format string) (config string,err error) {
	path := "/config?_queryFilter=_id+sw+\"\"&prettyPrint=true"

	req, err := http.NewRequest(http.MethodGet, idm.requestURL(path), nil)
	idm.setHeaders(req)
	result,err := crest.GetCRESTResult(req)

	if err != nil {
		return "",fmt.Errorf("Cant get config, err = %v", err)
	}
	// todo: fix the type issues...
	obj := &crest.FRObject{Kind: "idm.config", Items: &result.Result}

	return obj.Marshal(format)
}

func CreateObjects(obj *crest.FRObject, overwrite, continueOnError bool) (err error){
	idm, err := GetOpenIDMConnection()

	switch( obj.Kind ) {
	case IDM_CONFIGURATION:
		fmt.Printf("Create idm conf %v", idm)
	}
	return nil
}

// command ./cli.sh configexport --user openidm-admin:openidm-admin ./tmp
// reading configurations
// GET /config.
// loop thru  _id values, skipping org.apache
// GET /config/_id
// dump it out to file of the same name


