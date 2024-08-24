package tester

import (
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/fatih/color"
)

type WebClient struct {
	client    http.Client
	transport *Transport
}

func InitWebClient() *WebClient {
	t := InitTransport()
	return &WebClient{
		transport: t,
		client: http.Client{
			Transport: t,
		},
	}
}

func (w *WebClient) Get(url string) (*http.Response, error) {
	res, err := w.client.Get(url)
	s := res.Status + " " + url + " " + fmt.Sprintf("%.3fs", w.transport.Duration().Seconds())
	if res.StatusCode == 200 {
		color.Green(s)
	} else {
		color.HiRed(s)
	}

	return res, err
}

func (w *WebClient) GetWithPathParam(url, p string) (*http.Response, error) {
	if strings.LastIndex(url, "/") != len(url)-1 {
		url += "/" + p
	} else {
		url += p
	}

	res, err := w.client.Get(url)
	s := res.Status + " " + url + " " + fmt.Sprintf("%.3fs", w.transport.Duration().Seconds())
	if res.StatusCode == 200 {
		color.Green(s)
	} else {
		color.HiRed(s)
	}

	return res, err
}

func (w *WebClient) GetWithPathParams(url string, p []interface{}) (*http.Response, error) {
	if strings.LastIndex(url, "/") != len(url)-1 {
		url = fmt.Sprintf(url+"/", p...)
	} else {
		url = fmt.Sprintf(url, p...)
	}

	res, err := w.client.Get(url)
	s := res.Status + " " + url + " " + fmt.Sprintf("%.3fs", w.transport.Duration().Seconds())
	if res.StatusCode == 200 {
		color.Green(s)
	} else {
		color.HiRed(s)
	}

	return res, err
}

func (w *WebClient) GetWithQueryParams(url, ep string) (*http.Response, error) {
	url += "?" + ep
	res, err := w.client.Get(url)
	s := res.Status + " " + url + " " + fmt.Sprintf("%.3fs", w.transport.Duration().Seconds())
	if res.StatusCode == 200 {
		color.Green(s)
	} else {
		color.HiRed(s)
	}

	return res, err
}

type Transport struct {
	rtp       http.RoundTripper
	dialer    *net.Dialer
	connStart time.Time
	connEnd   time.Time
	reqStart  time.Time
	reqEnd    time.Time
}

func InitTransport() *Transport {
	t := &Transport{
		dialer: &net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		},
	}

	t.rtp = &http.Transport{
		Proxy:               http.ProxyFromEnvironment,
		Dial:                t.dial,
		TLSHandshakeTimeout: 10 * time.Second,
	}

	return t
}

func (t *Transport) RoundTrip(r *http.Request) (*http.Response, error) {
	t.reqStart = time.Now()
	res, err := t.rtp.RoundTrip(r)
	t.reqEnd = time.Now()
	return res, err

}

func (t *Transport) dial(network, addr string) (net.Conn, error) {
	t.connStart = time.Now()
	conn, err := t.dialer.Dial(network, addr)
	t.connEnd = time.Now()
	return conn, err
}

func (t *Transport) Duration() time.Duration {
	return t.ReqDuration() - t.ConnDuration()
}

func (t *Transport) ConnDuration() time.Duration {
	return t.connEnd.Sub(t.connStart)
}

func (t *Transport) ReqDuration() time.Duration {
	return t.reqEnd.Sub(t.reqStart)
}
