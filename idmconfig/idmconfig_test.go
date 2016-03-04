package idmconfig


import (
	"testing"
//	"os"
//log "github.com/Sirupsen/logrus"
)


// functional OpenIDM test. Assumes IDM is running

var base = "http://localhost:8080/openidm"

var c = IDMConnection{BaseURL:base,  Password:"openidm-admin",  User: "openidm-admin"}

// Basic functional test. Assume OpenAM is up and running at the URL below
func TestREST(t *testing.T) {

	if err := c.Authenticate(); err != nil {
		t.Fatalf("Could not authenticate, error is %v", err)
	}

	v,err := c.GET("/repo/internal/user/openidm-admin")

	if err != nil {
		t.Fatal(err)
	}
	t.Log(v)

	v2,err := c.POST("/system?_action=availableConnectors", []byte(``))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(v2)

	v3,err := c.GET("/config")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(v3)


}