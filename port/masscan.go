package port

import (
	"bufio"
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"

	log "github.com/sirupsen/logrus"
)

type Masscan struct {
	SystemPath string
	Args       []string
}

// SetSystemPath masscan可执行文件路径,默认不需要设置
func (m *Masscan) SetSystemPath(systemPath string) {
	if systemPath != "" {
		m.SystemPath = systemPath
	}
}
func (m *Masscan) SetArgs(arg ...string) {
	m.Args = arg
}

// SetTargets 扫描IP
func (m *Masscan) SetTargets(targets string) {
	m.Args = append(m.Args, targets)
}

// SetTargetFile 扫描IP
func (m *Masscan) SetTargetFile(targetFile string) {
	m.Args = append(m.Args, "-iL")
	m.Args = append(m.Args, targetFile)
}

// SetPorts 扫描端口范围
func (m *Masscan) SetPorts(ports string) {
	m.Args = append(m.Args, "-p")
	m.Args = append(m.Args, ports)
}

// SetRate 速率
func (m *Masscan) SetRate(rate int) {
	m.Args = append(m.Args, "--rate")
	m.Args = append(m.Args, strconv.Itoa(rate))
}

// SetExclude 隔离扫描名单
func (m *Masscan) SetExclude(exclude string) {
	m.Args = append(m.Args, "--exclude")
	m.Args = append(m.Args, exclude)
}

// SetOutput 设置输出文件的路径
func (m *Masscan) SetOutput(filepath string) {
	m.Args = append(m.Args, "-oX")
	m.Args = append(m.Args, filepath)
}

// Start scanning
func (m *Masscan) Run() error {
	cmd := exec.Command(m.SystemPath, m.Args...)
	log.Info(cmd.Args)
	var stdoutBuf, stderrBuf bytes.Buffer
	stdoutIn, _ := cmd.StdoutPipe()
	stderrIn, _ := cmd.StderrPipe()
	var errStdout, errStderr error
	stdout := io.MultiWriter(os.Stdout, &stdoutBuf)
	stderr := io.MultiWriter(os.Stderr, &stderrBuf)
	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("masscan start failed: %v", err)
	}
	go func() {
		_, errStdout = io.Copy(stdout, stdoutIn)
	}()
	go func() {
		_, errStderr = io.Copy(stderr, stderrIn)
	}()
	err = cmd.Wait()
	if err != nil {
		return fmt.Errorf("masscan wait failed: %v", err)
	}
	if errStdout != nil || errStderr != nil {
		return errors.New("failed to capture stdout or stderr")
	}
	return nil
}

// Parse scans result.
func (m *Masscan) Parse(filepath string) ([]Host, error) {
	var hosts []Host
	fi, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer fi.Close()
	r := bufio.NewReader(fi)
	decoder := xml.NewDecoder(r)
	for {
		t, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if t == nil {
			break
		}
		switch se := t.(type) {
		case xml.StartElement:
			if se.Name.Local == "host" {
				var host Host
				err := decoder.DecodeElement(&host, &se)
				if err == io.EOF {
					break
				}
				if err != nil {
					return nil, err
				}
				hosts = append(hosts, host)
			}
		default:
		}
	}
	return hosts, nil
}

func New() *Masscan {
	return &Masscan{
		SystemPath: "masscan",
	}
}

type Address struct {
	Addr     string `xml:"addr,attr"`
	AddrType string `xml:"addrtype,attr"`
}
type State struct {
	State     string `xml:"state,attr"`
	Reason    string `xml:"reason,attr"`
	ReasonTTL string `xml:"reason_ttl,attr"`
}
type Host struct {
	XMLName xml.Name `xml:"host"`
	Endtime string   `xml:"endtime,attr"`
	Address Address  `xml:"address"`
	Ports   Ports    `xml:"ports>port"`
}
type Ports []struct {
	Protocol string   `xml:"protocol,attr"`
	Portid   string   `xml:"portid,attr"`
	State    State    `xml:"state"`
	Service  MService `xml:"service"`
}
type MService struct {
	Name   string `xml:"name,attr"`
	Banner string `xml:"banner,attr"`
}
