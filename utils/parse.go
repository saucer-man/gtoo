package utils

import (
	"fmt"
	"strconv"

	"net/url"
	"strings"

	"golang.org/x/net/publicsuffix"
)

//GetPorts 从输入中解析
func GetPorts(selection string) ([]int, error) {
	ports := []int{}
	if selection == "" {
		return ports, nil
	}

	ranges := strings.Split(selection, ",")
	for _, r := range ranges {
		r = strings.TrimSpace(r)
		if strings.Contains(r, "-") {
			parts := strings.Split(r, "-")
			if len(parts) != 2 {
				return nil, fmt.Errorf("Invalid port selection segment: '%s'", r)
			}

			p1, err := strconv.Atoi(parts[0])
			if err != nil {
				return nil, fmt.Errorf("Invalid port number: '%s'", parts[0])
			}

			p2, err := strconv.Atoi(parts[1])
			if err != nil {
				return nil, fmt.Errorf("Invalid port number: '%s'", parts[1])
			}

			if p1 > p2 {
				return nil, fmt.Errorf("Invalid port range: %d-%d", p1, p2)
			}

			for i := p1; i <= p2; i++ {
				ports = append(ports, i)
			}

		} else {
			if port, err := strconv.Atoi(r); err != nil {
				return nil, fmt.Errorf("Invalid port number: '%s'", r)
			} else {
				ports = append(ports, port)
			}
		}
	}
	return ports, nil
}

//URL embeds net/url and adds extra fields ontop
type URL struct {
	Subdomain, Domain, TLD, Port string
	ICANN                        bool
	*url.URL
}

//ParseURL mirrors net/url.Parse except instead it returns
//a tld.URL, which contains extra fields.
func ParseURL(s string) (*URL, error) {
	url, err := url.Parse(s)
	if err != nil {
		return nil, err
	}
	if url.Host == "" {
		return nil, fmt.Errorf("can not find host from %s", s)
	}
	dom, port := domainPort(url.Host)
	//etld+1
	etld1, err := publicsuffix.EffectiveTLDPlusOne(dom)
	_, icann := publicsuffix.PublicSuffix(strings.ToLower(dom))
	if err != nil {
		return nil, err
	}
	//convert to domain name, and tld
	i := strings.Index(etld1, ".")
	domName := etld1[0:i]
	tld := etld1[i+1:]
	//and subdomain
	sub := ""
	if rest := strings.TrimSuffix(dom, "."+etld1); rest != dom {
		sub = rest
	}
	return &URL{
		Subdomain: sub,
		Domain:    domName,
		TLD:       tld,
		Port:      port,
		ICANN:     icann,
		URL:       url,
	}, nil
}

func domainPort(host string) (string, string) {
	for i := len(host) - 1; i >= 0; i-- {
		if host[i] == ':' {
			return host[:i], host[i+1:]
		} else if host[i] < '0' || host[i] > '9' {
			return host, ""
		}
	}
	//will only land here if the string is all digits,
	//net/url should prevent that from happening
	return host, ""
}
