package client

import (
	"crypto/tls"
	"io"
	"net/http"
)

func newClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
	}
}

func Do(t string, url string, header *map[string]string, data io.Reader) (*http.Response, error) {
	c := newClient()
	req, err := http.NewRequest(t, url, data)
	if err != nil {
		return nil, err
	}
	for k, v := range *header {
		req.Header.Add(k, v)
	}
	return c.Do(req)
}

func Get(url string, header *map[string]string, data io.Reader) (resp *http.Response, err error) {
	return Do("GET", url, header, data)
}

func Post(url string, header *map[string]string, data io.Reader) (resp *http.Response, err error) {
	return Do("POST", url, header, data)
}

func Put(url string, header *map[string]string, data io.Reader) (resp *http.Response, err error) {
	return Do("PUT", url, header, data)
}

func Delete(url string, header *map[string]string, data io.Reader) (resp *http.Response, err error) {
	return Do("DELETE", url, header, data)
}
