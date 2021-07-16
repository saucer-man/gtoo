package utils

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"time"
)

var proxyUrl = "http://127.0.0.1:8080"

var proxy, _ = url.Parse(proxyUrl)

var tr = &http.Transport{
	// Proxy:           http.ProxyURL(proxy),
	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
}

var Client = &http.Client{
	Transport: tr,
	Timeout:   time.Second * 20, //超时时间
}
