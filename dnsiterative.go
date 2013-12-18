package dnsiterative

import (
	"errors"
	"fmt"
	"github.com/miekg/dns"
	"math/rand"
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

func lookup(cl *dns.Client, name, server string, recType RecordType, toMatch ...string) (bool, error) {
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
			return lookup(cl, name, fmt.Sprintf("%s:53", ns.Ns[0:len(ns.Ns)-1]), recType, toMatch...)
		}
	} else {
		for _, rr := range response.Answer {
			if recType == A {
				if a, ok := rr.(*dns.A); ok {
					for _, m := range toMatch {
						if m == a.A.String() {
							return true, nil
						}
					}
				}
			} else if recType == CNAME {
				if cn, ok := rr.(*dns.CNAME); ok {
					for _, m := range toMatch {
						if m == cn.Target {
							return true, nil
						}
					}
				}

			}
		}
		return false, nil
	}
}

//
// This will check if a DNS name has at least 1 record (of recType) that matches ones of the values within
// `addressesToMatch`.
//
// Take care to pass proper DNS names (i.e 'google.com.' not 'google.com') for both `name` and `addressesToMatch`.
//
// Will randomly select a root server from `DnsRoots` every time it runs, but will not retry failures.
//
func DomainHasRecord(name string, recType RecordType, addressesToMatch ...string) (bool, error) {
	cl := new(dns.Client)
	cl.Net = "tcp"
	return lookup(cl, name, DnsRoots[rand.Intn(len(DnsRoots))], recType, addressesToMatch...)
}
