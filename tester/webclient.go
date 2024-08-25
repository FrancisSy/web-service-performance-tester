package tester

import (
	"bytes"
	"fmt"
	"log"
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
			Timeout:   30 * time.Second,
		},
	}
}

func (w *WebClient) Get(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	res, err := w.client.Do(req)
	s := res.Status + " " + url + " " + fmt.Sprintf("%.3fs", w.transport.Duration().Seconds())
	if Is2xxSuccessful(res) {
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

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	res, err := w.client.Do(req)
	s := res.Status + " " + url + " " + fmt.Sprintf("%.3fs", w.transport.Duration().Seconds())
	if Is2xxSuccessful(res) {
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

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	res, err := w.client.Do(req)
	s := res.Status + " " + url + " " + fmt.Sprintf("%.3fs", w.transport.Duration().Seconds())
	if Is2xxSuccessful(res) {
		color.Green(s)
	} else {
		color.HiRed(s)
	}

	return res, err
}

func (w *WebClient) GetWithQueryParams(url, ep string) (*http.Response, error) {
	url += "?" + ep
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	res, err := w.client.Do(req)
	s := res.Status + " " + url + " " + fmt.Sprintf("%.3fs", w.transport.Duration().Seconds())
	if Is2xxSuccessful(res) {
		color.Green(s)
	} else {
		color.HiRed(s)
	}

	return res, err
}

func (w *WebClient) Post(url string, body []byte) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	res, err := w.client.Do(req)
	s := res.Status + " " + url + " " + fmt.Sprintf("%.3fs", w.transport.Duration().Seconds())
	if Is2xxSuccessful(res) {
		color.Green(s)
	} else {
		color.HiRed(s)
	}

	return res, err
}

func (w *WebClient) Patch(url string, body []byte) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPatch, url, bytes.NewReader(body))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	res, err := w.client.Do(req)
	s := res.Status + " " + url + " " + fmt.Sprintf("%.3fs", w.transport.Duration().Seconds())
	if Is2xxSuccessful(res) {
		color.Green(s)
	} else {
		color.HiRed(s)
	}

	return res, err
}

func SetHeaders(r *http.Request, m map[string]string) {
	for k, v := range m {
		r.Header.Set(k, v)
	}
}

func Is2xxSuccessful(r *http.Response) bool {
	status := r.StatusCode
	return status >= 200 && status <= 299
}

func Is3xxRedirection(r *http.Response) bool {
	status := r.StatusCode
	return status >= 300 && status <= 399
}

func Is4xxClientError(r *http.Response) bool {
	status := r.StatusCode
	return status >= 400 && status <= 499
}

func Is5xxServerError(r *http.Response) bool {
	status := r.StatusCode
	return status >= 500 && status <= 599
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
