package crest

import (
	"net/http"
	"io/ioutil"
	log "github.com/Sirupsen/logrus"
	"encoding/json"
	"github.com/ghodss/yaml"
	"io"
	"fmt"
)


// A query response returned from CREST
type CRESTResult struct {
	Result                 []interface{}  `json:"result"`
	ResultCount            int64          `json:"resultCount"`
	PagedResultsCookie     string         `json:"pagedResultsCookie`
	RemainingPagedResults  int64          `json:"remainingPagedResults"`
	TotalPagedResults      int64          `json:"totalPagedResults"`
	TotalPagedResultPolicy string         `json:"totalPagedResultsPolicy"`
}

type FRObject struct {
	Kind  string `json:"kind"`
	Items *[]interface{}  `json:"spec"`
}


// Make a CREST request. http is already set up with the method, url and any cookies/ headers
func GetCRESTResult(req *http.Request) (CRESTResult, error) {
	client := &http.Client{}

	resp, err := client.Do(req)
	defer resp.Body.Close()

	var result CRESTResult

	if err != nil {
		return result, err
	}


	body, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return result, fmt.Errorf("Bad Response %v %v\n", resp.StatusCode, string(body))
	}

	//log.Println("response Body:", string(body))

	err = json.Unmarshal(body, &result)

	if err != nil {
		log.Fatalf("cant unmarshal REST result: %v", err)
		return result, err
	}

	//for k,v := range result.Result {
	//	log.Infof("k= %v  v= %#v", k,v)
	//	//var m map[string]interface{}
	//	//json.Unmarshal(v,&m)
	//	//log.Infof("m =%v", m)
	//	//for k,v := range v {
	//	//	log.Infof("k,v =", k, v)
	//	//}
	//}

	return result, nil
}



// Read in FR Object from a stream. The stream is assumed to be JSON or YAML
// We use the YAML decoder because YAML is a superset of JSON
// So this works with both types
func ReadFRConfig(in io.Reader) (obj *FRObject, err error) {
	obj = &FRObject{}

	if b,err := ioutil.ReadAll(in); err == nil {
		err = yaml.Unmarshal(b,&obj)
		if err != nil {
			fmt.Printf("Cant unmarshall file err=%v", err)
		}
		// fmt.Printf("Kind = %s, items=%v", obj.Kind, obj.Items)
	}
	return
}

