package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type HttpClient interface {
	Get(url string, res interface{}, param ...Param) error
	Post(url string, payload interface{}, res interface{}, header ...map[string]string) error
}

type Param struct {
	Query  map[string]string
	Header map[string]string
}

type Client struct {
	http.Client
}

func New() HttpClient {
	return &Client{
		Client: http.Client{},
	}
}

func (c *Client) Get(url string, res interface{}, param ...Param) error {
	var body []byte
	var err error
	if len(param) > 0 {
		body, err = c.get(url, param[0])
	} else {
		body, err = c.get(url)
	}
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return fmt.Errorf("decode resp to struct err: %s, resp body: %s", err.Error(), string(body))
	}
	return nil
}

func (c *Client) Post(url string, payload interface{}, res interface{}, header ...map[string]string) error {
	var body []byte
	var err error
	if len(header) > 0 {
		body, err = c.post(url, payload, header[0])
	} else {
		body, err = c.post(url, payload)
	}
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return fmt.Errorf("decode resp to struct err: %s, resp body: %s", err.Error(), string(body))
	}
	return nil
}

func (c *Client) get(url string, param ...Param) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	// deal with headers and query param if passed
	if len(param) > 0 {
		for k, v := range param[0].Query {
			q := req.URL.Query()
			q.Add(k, v)
			req.URL.RawQuery = q.Encode()
		}
		for k, v := range param[0].Header {
			req.Header.Add(k, v)
		}
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http status code %d", resp.StatusCode)
	}
	return ioutil.ReadAll(resp.Body)
}

func (c *Client) post(url string, payload interface{}, header ...map[string]string) ([]byte, error) {
	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return nil, err
	}
	if len(header) > 0 {
		for k, v := range header[0] {
			req.Header.Add(k, v)
		}
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http status code %d", resp.StatusCode)
	}
	return ioutil.ReadAll(resp.Body)
}
