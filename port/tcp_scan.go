package port

// port tcp scan

import (
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

type PortScanner struct {
	host      string
	timeout   time.Duration
	threads   int
	startport int
	endport   int
}

func (p *PortScanner) SetPort(rawport string) (err error) {
	s := strings.Split(rawport, "-")
	startport, err := strconv.Atoi(s[0])
	if err != nil {
		return
	}
	var endport int = 0
	if len(s) == 1 {
		endport = startport + 1
	} else if len(s) == 2 {
		endport, err = strconv.Atoi(s[1])
		if err != nil {
			fmt.Printf("cant parse scan port with %v\n", rawport)
			return
		}
	} else {
		err = errors.New("cant parse scan port with" + rawport)
		return
	}

	//解析完了之后开看规则
	if startport <= 0 || startport >= 65535 {
		err = errors.New("The port to start scanning is out of range")
		return
	}
	if endport <= 1 || endport >= 65536 {
		err = errors.New("The port to end scanning is out of range")
		return
	}
	if endport <= startport {
		err = errors.New("the start port of the scan cannot be larger than the end port")
		return
	}

	p.startport = startport
	p.endport = endport
	log.Printf("set ports range: %d-%d.\n", p.startport, p.endport)
	return
}

func (p *PortScanner) SetThreads(threads int) {
	p.threads = threads
	log.Printf("set threads %d.\n", p.threads)
}
func (p *PortScanner) SetHosts(host string) error {
	addr, err := net.ResolveIPAddr("ip4", host)
	if err != nil {
		return err
	}
	p.host = addr.String()
	log.Printf("set host %s.\n", p.host)
	return nil
}
func (p *PortScanner) SetTimeout(timeout time.Duration) {
	p.timeout = timeout
	log.Printf("set timeout %s.\n", p.timeout)
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
	for port := p.startport; port <= p.endport; port++ {
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
	return rv
}
