package amconfig

import( "fmt"
	log "github.com/Sirupsen/logrus"
	"net/http"
	"net/http/httputil"

	"github.com/davecgh/go-spew/spew"
"io/ioutil"
"encoding/json"
	"io"
	"github.com/forgerock/frconfig/crest"
)

// ResourceType is an OpenAM policy resource type
type ResourceType struct {
	UUID             string      `json: "uuid"`
	Name             string      `json: "name"`
	Description      string      `json: "description"`
	Patterns         []string    `json: "patterns"`
	Actions          interface{} `json: "actions"`
	CreatedBy        string      `json: "createdBy"`
	CreationDate     int64       `json: "creationDate"`
	LastModifiedBy   string      `json: "lastModifiedBy"`
	LastModifiedDate int64       `json: "lastModifiedDate"`
}

type ResourceTypeResult struct {
	Result                []ResourceType `json: "result"`
	ResultCount           int64          `json: "resultCount"`
	PagedResultsCookie    string         `json: "pagedResultsCookie`
	RemainingPagedResults int64          `json: "remainingPagedResults"`
}


// ListResourceTypes returns the OpenAM policy resource types
func (openam *OpenAMConnection)ListResourceTypes() ([]ResourceType, error) {

	client := &http.Client{}
	req := openam.newRequest("GET", "/json/resourcetypes?_queryFilter=true", nil)

	dump, err := httputil.DumpRequestOut(req, true)

	fmt.Printf("dump req is %s", dump)

	resp, err := client.Do(req)
	defer resp.Body.Close()

	debug(httputil.DumpResponse(resp, true))


	if err != nil {
		return nil, err
	}

	body, _ := ioutil.ReadAll(resp.Body)


	var result ResourceTypeResult

	err = json.Unmarshal(body, &result)

	if err != nil {
		log.Fatalf("cant get result type", err)
	}

	//fmt.Printf("result = %s", result)

	spew.Dump(result)

	return result.Result, err
}


func CreateFRObjects(in io.Reader) (err error) {

	obj,err := crest.ReadFRConfig(in)

	if err != nil {
		return err
	}

	var am *OpenAMConnection

	// TOOD: Use constants here...
	switch obj.Kind {
	case "policy":
		am,err = GetOpenAMConnection()
		err = am.CreatePolicies(obj,true,true)
	default:
		err = fmt.Errorf("Unknown object type %s", obj.Kind)
	}
	return
}
