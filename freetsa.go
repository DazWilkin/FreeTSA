package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	freeTSA = "https://freetsa.org"
)

// Client represents a FreeTSA client
type Client struct {
	c *http.Client
}

// NewClient creates a new FreeTSA Client
func NewClient(c *http.Client) *Client {
	return &Client{
		c: c,
	}
}

// Request submits a DER-encoded request to FreeTSA
func (c *Client) Request(der []byte) ([]byte, error) {
	rqst, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s/tsr", freeTSA),
		bytes.NewReader(der),
	)
	if err != nil {
		return nil, err
	}

	rqst.Header.Set(
		"Content-Type",
		"application/timestamp-query",
	)

	resp, err := c.c.Do(rqst)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%d: %s", resp.StatusCode, resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
