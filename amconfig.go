package main

import (
	"flag"
	log "github.com/Sirupsen/logrus"
	"github.com/forgerock/frconfig/amconfig"
	"os"
)

func main() {
	url := 		flag.String("url", "http://openam.test.com:8080/openam", "OpenAM URL")
	user := 	flag.String("user", "amadmin", "Admin User")
	password := 	flag.String("password", "password", "Admin Password")
	filePath := 	flag.String("path", "/tmp/policy.json", "Policy File")

	flag.Parse()

	am, err := amconfig.Open(*url, *user, *password)

	if err != nil {
		log.Fatalf("Can't create connection: %v", err)
		os.Exit(1)
	}


	am.ExportPoliciesToJSONFile(*filePath)
	am.ImportPoliciesFromFile(*filePath)


}
