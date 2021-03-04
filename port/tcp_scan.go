package port

// port tcp scan

import (
	"fmt"
	"gtoo/utils"
	"net"
	"sort"
	"sync"
	"time"
)

type PortScanner struct {
	host    string
	timeout time.Duration
	threads int
	ports   []int
}

func (p *PortScanner) SetPort(rawport string) (err error) {
	p.ports, err = utils.GetPorts(rawport)
	return err
}

func (p *PortScanner) SetThreads(threads int) {
	p.threads = threads
}

func (p *PortScanner) SetHosts(host string) error {
	addr, err := net.ResolveIPAddr("ip4", host)
	if err != nil {
		return err
	}
	p.host = addr.String()
	return nil
}
func (p *PortScanner) SetTimeout(timeout time.Duration) {
	p.timeout = timeout

}
func (p PortScanner) IsOpen(port int) bool {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", p.host, port), p.timeout)
	if err != nil {
		return false
	}

	defer conn.Close()

	return true
}

func (p PortScanner) GetOpenedPort() []int {
	rv := []int{}
	l := sync.Mutex{}
	sem := make(chan bool, p.threads) // sem是为了开线程
	for port := range p.ports {
		sem <- true
		go func(port int) {
			if p.IsOpen(port) {
				l.Lock()
				rv = append(rv, port)
				l.Unlock()
			}
			<-sem
		}(port)
	}
	for i := 0; i < cap(sem); i++ {
		sem <- true
	}
	sort.Ints(rv)
	return rv
}
