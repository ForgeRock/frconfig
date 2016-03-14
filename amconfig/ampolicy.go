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
	"os"
	crest "github.com/forgerock/frconfig/common"
	"github.com/rackspace/gophercloud/openstack/networking/v2/extensions/fwaas/policies"
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

func PolicytoYAML(policies []Policy) {
	

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

// Export all the policies as a JSON policy set string
func (openam *OpenAMConnection) ExportJSONPolicies() (string, error) {
	req := openam.newRequest("GET", "/json/policies?_queryFilter=true")

	result,err := crest.GetCRESTResult(req)
	if err != nil {
		log.Fatalf("Could not get policies, err=%v",err)
		return "",err
	}

	log.Infof("Crest result = %+v", result)

	b,err := json.Marshal(result.Result)
	if err != nil {
		log.Fatalf("can't marshal json: %v ", err)
		return "",nil
	}

	return string(b),err

}

func (openam *OpenAMConnection) ExportPoliciesToJSONFile(filePath string) error {
	p,err := openam.ExportJSONPolicies()
	if err != nil {
		return err
	}
	f, err := os.Create(filePath)
	defer f.Close()

	if err != nil {
		log.Fatalf("Could create policy file %v err = %v", filePath, err)
		return err
	}


	_, err = f.WriteString(p)

	if err != nil {
		log.Fatalf("could not write policy file: %v", err)
		return err
	}
	return nil
}


type  PolicyArray  []interface{}

func (openam *OpenAMConnection) ImportPoliciesFromFile(filePath string)  error {
	f,err := os.Open(filePath)
	defer f.Close()
	if err != nil {
		log.Errorf("Can't open file %v, err=%v", filePath, err)
	}
	//r := bufio.NewReader(f)

	bytes, err := ioutil.ReadAll(f)

	if err != nil {
		log.Errorf("Can't read policy file. Err = %v", err)
		return err
	}

	var p PolicyArray

	err = json.Unmarshal(bytes, &p)

	if err != nil {
		log.Fatalf("Can't unmarshal json file, Err=%v", err)
	}

	for _,v := range p {
		log.Debugf("policy=%v",  v)
	}
	return err

}

// Export all the policies as a JSON policy set string
func (openam *OpenAMConnection) ImportJSONPolicies(p PolicyArray) error {
	req := openam.newRequest("POST", "/json/policies?_action=create")

	result,err := crest.GetCRESTResult(req)
	if err != nil {
		log.Fatalf("Could not get policies, err=%v",err)
		return "",err
	}

	log.Infof("Crest result = %+v", result)

	b,err := json.Marshal(result.Result)
	if err != nil {
		log.Fatalf("can't marshal json: %v ", err)
		return "",nil
	}

	return string(b),err

}

// Script query - to get Uuid
// http://openam.test.com:8080/openam/json/scripts?_pageSize=20&_sortKeys=name&_queryFilter=true&_pagedResultsOffset=0

