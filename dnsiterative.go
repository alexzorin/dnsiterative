package dnsiterative

import (
	"errors"
	"fmt"
	"math/rand"

	"github.com/miekg/dns"
)

type RecordType string

const (
	A     RecordType = "A"
	CNAME RecordType = "CNAME"
)

var (
	DnsRoots = []string{
		"a.root-servers.net:53", "b.root-servers.net:53",
		"c.root-servers.net:53", "d.root-servers.net:53",
		"e.root-servers.net:53", "f.root-servers.net:53",
		"g.root-servers.net:53", "h.root-servers.net:53",
		"i.root-servers.net:53", "j.root-servers.net:53",
		"k.root-servers.net:53", "l.root-servers.net:53",
		"m.root-servers.net:53",
	}
	ErrNoNameservers = errors.New("No nameservers registered for that domain")
	ErrUnhandled     = errors.New("Unknown error")
)

// A matcher will match if both the Type (A,CNAME) and Value match exactly
// with the record returned by the DNS server
type Matcher struct {
	Type  RecordType
	Value string
}

func (m *Matcher) matches(val dns.RR) bool {
	switch m.Type {
	case A:
		if a, ok := val.(*dns.A); ok {
			return a.A.String() == m.Value
		} else {
			return false
		}
	case CNAME:
		if cn, ok := val.(*dns.CNAME); ok {
			return cn.Target == m.Value
		} else {
			return false
		}
	}
	return false
}

func lookup(cl *dns.Client, name, server string, matchers ...Matcher) (bool, error) {
	msg := new(dns.Msg)
	msg.Id = dns.Id()
	msg.RecursionDesired = false

	msg.Question = []dns.Question{
		dns.Question{
			name, dns.TypeA, dns.ClassINET,
		},
	}

	response, _, err := cl.Exchange(msg, server)
	if err != nil {
		return false, err
	}
	if len(response.Answer) == 0 {
		if len(response.Ns) == 0 {
			return false, ErrNoNameservers
		} else {
			ns, ok := response.Ns[rand.Intn(len(response.Ns))].(*dns.NS)
			if !ok {
				return false, ErrUnhandled
			}
			return lookup(cl, name, fmt.Sprintf("%s:53", ns.Ns[0:len(ns.Ns)-1]), matchers...)
		}
	} else {
		for _, matcher := range matchers {
			for _, rr := range response.Answer {
				if matcher.matches(rr) {
					return true, nil
				}
			}
		}
		return false, nil
	}
}

//
// This will check if a DNS name has at least 1 record that meets the requirements outlined in `matchers`
//
// Take care to pass proper DNS names (i.e 'google.com.' not 'google.com') for both `name` and `addressesToMatch`.
//
// Will randomly select a root server from `DnsRoots` every time it runs, but will not retry failures.
//
func DomainHasRecord(name string, matchers ...Matcher) (bool, error) {
	cl := new(dns.Client)
	return lookup(cl, name, DnsRoots[rand.Intn(len(DnsRoots))], matchers...)
}
