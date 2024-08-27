package webtest

import (
	"bytes"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/fatih/color"
)

type WebClient struct {
	client    http.Client
	headers   *map[string]string
	transport *Transport
}

func InitWebClient() *WebClient {
	t := InitTransport()
	return &WebClient{
		transport: t,
		headers:   nil,
		client: http.Client{
			Transport: t,
		},
	}
}

func (w *WebClient) Get(url string) (*http.Response, error) {
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	return w.execute(req, url)
}

func (w *WebClient) Post(url string, body []byte) (*http.Response, error) {
	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	return w.execute(req, url)
}

func (w *WebClient) Patch(url string, body []byte) (*http.Response, error) {
	req, _ := http.NewRequest(http.MethodPatch, url, bytes.NewReader(body))
	return w.execute(req, url)
}

func (w *WebClient) Put(url string, body []byte) (*http.Response, error) {
	req, _ := http.NewRequest(http.MethodPut, url, bytes.NewReader(body))
	return w.execute(req, url)
}

func (w *WebClient) Delete(url, p string) (*http.Response, error) {
	req, _ := http.NewRequest(http.MethodDelete, url+p, nil)
	return w.execute(req, url)
}

func (w *WebClient) Headers(m *map[string]string) *WebClient {
	w.headers = m
	return w
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

func (w *WebClient) execute(req *http.Request, url string) (*http.Response, error) {
	if w.headers != nil {
		SetHeaders(req, *w.headers)
	}

	res, err := w.client.Do(req)
	s := res.Status + " " + url + " " + fmt.Sprintf("%.3fs", w.transport.Duration().Seconds())
	if Is2xxSuccessful(res) {
		color.Green(s)
	} else if Is3xxRedirection(res) {
		color.Yellow(s)
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
