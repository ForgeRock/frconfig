package amconfig

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	log "github.com/Sirupsen/logrus"
	"net/http"
	"net/http/httputil"

	"github.com/davecgh/go-spew/spew"
	"github.com/ghodss/yaml"
)

// Policy in OpenAMConnection
type Policy struct {
	Name             string      `json:"name"`
	Active           bool        `json:"active"`
	ApplicationName  string      `json:"applicationName"`
	ActionValues     interface{} `json:"actionValues"`
	Resources        []string    `json:"resources"`
	Description      string      `json:"description"`
	Subject          interface{} `json:"subject"`
	Condition        interface{} `json:"condition"`
	ResourceTypeUUID string      `json:"resourceTypeUuid"`
	CreatedBy        string      `json:"createdBy"`
	CreationDate     string      `json:"creationDate"`
	LastModifiedBy   string      `json:"lastModifiedBy"`
	LastModifiedDate string      `json:"lastModifiedDate"`
}

// A PolicyResultList is a set of Policies
type PolicyResultList struct {
	Result                []Policy `json:"result"`
	ResultCount           int64    `json:"resultCount"`
	PagedResultsCookie    string   `json:"pagedResultsCookie`
	RemainingPagedResults int64    `json:"remainingPagedResults"`
}

// ListPolicy lists all OpenAM policies for a realm
func ListPolicy(openam *OpenAMConnection) ([]Policy, error) {

	client := &http.Client{}
	req := openam.newRequest("GET", "/json/policies?_queryFilter=true")

	dump, err := httputil.DumpRequestOut(req, true)

	fmt.Printf("dump req is %s", dump)

	//debug(httputil.DumpResponse(response, true))

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	//fmt.Println("response Status:", resp.Status)
	//fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)

	fmt.Println("response Body:", string(body))

	var result PolicyResultList

	err = json.Unmarshal(body, &result)

	if err != nil {
		log.Fatalf("cant get result type", err)
	}

	//fmt.Printf("result = %s", result)

	spew.Dump(result)

	return result.Result, err
}

func PolicytoJSON(policies []Policy) {
	

	for _, p := range policies {
		s, err := json.Marshal(p)
		if err != nil {
			fmt.Println("error %v", err)
		} else {

			fmt.Println("json:", string(s))
		}
		
		y,err := yaml.Marshal(p)
		if err != nil {
			fmt.Println("error %v", err)
		} else {

			fmt.Println("yaml:\n", string(y))
		}

	}

}

// Export all the policies as a XACML policy set
func (openam *OpenAMConnection) ExportXacmlPolicies() (string, error) {
	req := openam.newRequest("GET", "/xacml/policies")

	client := &http.Client{}

	resp, err := client.Do(req)
	defer resp.Body.Close()

	if err != nil {
		return "", err
	}



	body, err := ioutil.ReadAll(resp.Body)

	return string(body),err

}

// Export all the policies as a JSON policy set
func (openam *OpenAMConnection) ExportJSONPolicies() (string, error) {
	req := openam.newRequest("GET", "/json/policies?_queryFilter=true")

	client := &http.Client{}

	resp, err := client.Do(req)
	defer resp.Body.Close()

	if err != nil {
		return "", err
	}


	body, err := ioutil.ReadAll(resp.Body)

	return string(body),err

}


// Script query - to get Uuid
// http://openam.test.com:8080/openam/json/scripts?_pageSize=20&_sortKeys=name&_queryFilter=true&_pagedResultsOffset=0




