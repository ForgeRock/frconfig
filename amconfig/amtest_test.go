package amconfig

import (
	"testing"
	"os"
	//log "github.com/Sirupsen/logrus"
)

var base = "http://openam.test.com:8080/openam"
var amc = OpenAMConnection{BaseURL:base,  Password:"password",  User: "amadmin"}

// Basic functional test. Assume OpenAM is up and running at the URL below
func TestREST(t *testing.T) {

	if err := amc.Authenticate(); err != nil {
		t.Fatalf("Could not authenticate")
	}


	rt, err := amc.ListResourceTypes()

	if err != nil {
		t.Fatalf("Could not List resoruce types")
	}

	for _, v := range rt {
		t.Logf("%v", v)
	}

}

func TestXACML(t *testing.T) {

	if err := amc.Authenticate(); err != nil {
		t.Fatalf("Could not authenticate")
	}


	xacml, err := amc.ExportXacmlPolicies()

	t.Logf("\n Xacml = %v", xacml)

	if err != nil {
		t.Fatalf("Could not get policies", err)
	}

	f, err := os.Create("/var/tmp/xacml.xml")
	defer f.Close()

	if err != nil {
		t.Fatalf("Could create policy file", err)
	}
	_, err = f.WriteString(xacml)

	if err != nil {
		t.Fatalf("could not write policy file", err)
	}

}


func TestJSONExport(t *testing.T) {

	if err := amc.Authenticate(); err != nil {
		t.Fatalf("Could not authenticate")
	}

	json,err := amc.ExportJSONPolicies()

	if err != nil {
		t.Fatalf("Cant read json policies %v", err)
	}
	t.Log(json)
}