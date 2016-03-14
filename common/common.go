package common

import (
	"net/http"
	"io/ioutil"
	log "github.com/Sirupsen/logrus"

	"encoding/json"
)

// todo: many of these fields can probably be package private
type CRESTResult struct {
	Result                 []interface{}  `json: "result"`
	ResultCount            int64          `json: "resultCount"`
	PagedResultsCookie     string         `json: "pagedResultsCookie`
	RemainingPagedResults  int64          `json: "remainingPagedResults"`
	TotalPagedResults      int64          `json: "totalPagedResults"`
	TotalPagedResultPolicy string         `json: "totalPagedResultsPolicy"`
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

	log.Println("response Body:", string(body))

	err = json.Unmarshal(body, &result)

	if err != nil {
		log.Fatalf("cant unmarshal REST result: %v", err)
		return result, err
	}

	for k,v := range result.Result {
		log.Infof("k= %v  v= %#v", k,v)
		//var m map[string]interface{}
		//json.Unmarshal(v,&m)
		//log.Infof("m =%v", m)
		//for k,v := range v {
		//	log.Infof("k,v =", k, v)
		//}
	}


	return result, nil
}