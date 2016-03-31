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

// Represents a serialized object on disk. We add some metadata and object type
// information so we can deserialize to the right object
type FRObject struct {
	Kind  		string 			`json:"kind"` // Object type
	Metadata 	map[string]string  	`json:"metadata"`
	Items 		*[]interface{}  	`json:"spec"`
}

// A function that knows how to create an Object of a certain type.
// For example, create a policy object in OpenAM
type CreateObjectHandler func (obj *FRObject, overwrite, continueOnError bool) (error)

var handlerMap =  make (map[string]CreateObjectHandler)

// register Create Object handlers for the given object types
func RegisterCreateObjectHandler(objectKinds []string, handler CreateObjectHandler) {
	for _,kind := range objectKinds {
		handlerMap[kind] = handler
	}
}


// Make a CREST request. http is already set up with the method, url and any cookies/ headers
func GetCRESTResult(req *http.Request) (CRESTResult, error) {
	client := &http.Client{}

	resp, err := client.Do(req)

	log.Debugf("Request %+v", req)
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

// Marshal an FRObject to a json or yaml string

func  (obj *FRObject) Marshal(format string) (str string,err error) {

	var  b []byte

	switch format {
	case "yaml":
		b,err = yaml.Marshal(obj)
	case "json":
		b,err = json.MarshalIndent(obj, "", "  ")
	default:
		return "", fmt.Errorf("Unrecognized output format %s", format)

	}

	return string(b),err
}


// Read in an Object from the reader and create it in the ForgeRock stack
// todo: Check for overwrite,
func CreateFRObjects(in io.Reader, overwrite,continueOnError bool) (err error) {

	obj,err := ReadFRConfig(in)

	if err != nil {
		return err
	}

	if handler, ok := handlerMap[obj.Kind]; ok {
		return handler(obj,overwrite,continueOnError)
	}
	return fmt.Errorf("Don't know how to create an object of type %s", obj.Kind)

	//
	//var am *OpenAMConnection
	//
	//// is object type meant for OpenAM?
	//if strings.HasPrefix(obj.Kind,"am.") {
	//	am,err = GetOpenAMConnection()
	//	if err != nil {
	//		return
	//	}
	//}
	//
	//fmt.Printf("Handling object %s", obj.Kind)
	//switch obj.Kind {
	//case POLICY:
	//	err = am.CreatePolicies(obj,overwrite,continueOnError)
	//default:
	//	err = fmt.Errorf("Unknown object type %s", obj.Kind)
	//}
	//return
}
