package webtest

import (
	"bytes"
	"fmt"
	"net/http"

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
