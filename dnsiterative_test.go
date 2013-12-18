package dnsiterative

import (
	"math/rand"
	"testing"
	"time"
)

func TestLookupASuccess(t *testing.T) {
	rand.Seed(time.Now().Unix())
	hasRecord, err := DomainHasRecord("serversaurus.com.au.", Matcher{A, "111.223.231.133"}, Matcher{CNAME, "barney.serversaurus.com.au"})
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	if !hasRecord {
		t.Log("Domain should have had record, but didn't")
		t.Fail()
	}
}

func TestLookupAFail(t *testing.T) {
	hasRecord, err := DomainHasRecord("google.com.", Matcher{A, "1.1.1.1"})
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	if hasRecord {
		t.Log("Domain should not have had record")
		t.Fail()
	}
}

func TestLookupCNAMESuccess(t *testing.T) {
	rand.Seed(time.Now().Unix())
	hasRecord, err := DomainHasRecord("www.reddit.com.", Matcher{CNAME, "reddit.com.edgesuite.net."})
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	if !hasRecord {
		t.Log("Domain should have had record, but didn't")
		t.Fail()
	}
}
