package main

import (
	//"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	log "github.com/Sirupsen/logrus"
	//"net/http"
	"os"
    //"encoding/json"
	"gopkg.in/yaml.v2"
   am "github.com/wstrange/frconfig/amconfig"
)

//var fileFlag = flag.String("file","config.yaml", "Configuration file to load")


func main() {


	openam := am.OpenAMConnection {
		BaseURL: "http://openam.example.com:8080/openam" ,User: "amadmin", Password: "password" ,
	}
			
	err := openam.Authenticate()

	
	//resultTypes, err  := am.ListResourceTypes(&openam)
	policies, err  := am.ListPolicy(&openam)

	am.PolicytoJSON(policies)
	
	//fmt.Printf("resultTypes = %#v", resultTypes)
	for _, rs := range policies {
		fmt.Printf("resultType =%#v\n\n", rs)
	}
	
	
	fileFlag := flag.String("file", "config.yaml", "Configuration file to load")

	flag.Parse()

	_, err = os.Stat(*fileFlag)
	if err != nil {
		log.Fatal("Config file is missing: ", *fileFlag)
	}

	data, err := ioutil.ReadFile(*fileFlag)
	if err != nil {
		log.Fatalf("Error reading configuration file %s err %v", *fileFlag, err)
	}


	t := T{}

	err = yaml.Unmarshal([]byte(data), &t)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	//fmt.Printf("--- t:\n%v\n\n", t)

	// d, err := yaml.Marshal(&t)
	// if err != nil {
	// 	log.Fatalf("error: %v", err)
	// }
	// //fmt.Printf("--- t dump:\n%s\n\n", string(d))

	// m := make(map[interface{}]interface{})

	// err = yaml.Unmarshal([]byte(data), &m)
	// if err != nil {
	// 	log.Fatalf("error: %v", err)
	// }
	// //fmt.Printf("--- m:\n%v\n\n", m)

	// d, err = yaml.Marshal(&m)
	// if err != nil {
	// 	log.Fatalf("error: %v", err)
	// }
	//fmt.Printf("--- m dump:\n%s\n\n", string(d))
}
