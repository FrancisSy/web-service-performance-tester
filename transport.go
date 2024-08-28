package webtest

import (
	"net"
	"net/http"
	"time"
)

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
