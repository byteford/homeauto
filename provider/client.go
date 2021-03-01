package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// Client -
type Client struct {
	HostURL    string
	HTTPClient *http.Client
	Token      string
}

//NewClient -
func NewClient(host, token *string) (*Client, error) {
	if token == nil {
		return nil, fmt.Errorf("no token")
	}
	c := Client{
		HTTPClient: &http.Client{},
		HostURL:    "",
		Token:      fmt.Sprintf("Bearer %s", *token),
	}

	if host != nil {
		c.HostURL = *host
	}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/", c.HostURL), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", c.Token)
	_, err = c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

//GetLight -
func (c *Client) GetLight(lightID string) (*LightItem, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/states/%s", c.HostURL, lightID), nil)
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)

	if err != nil {
		return nil, err
	}
	light := LightItem{}
	err = json.Unmarshal(body, &light)
	if err != nil {
		return nil, err
	}
	return &light, nil
}

//StartLight -
func (c *Client) StartLight(lightItem LightItem) (*LightItem, error) {
	rb, err := json.Marshal(lightItem)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/states/%s", c.HostURL, lightItem.EntityID), strings.NewReader(string(rb)))
	if err != nil {
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, fmt.Errorf("err in request: %v", err.Error(), req.URL, c.Token)
	}
	light := LightItem{}
	err = json.Unmarshal(body, &light)
	if err != nil {
		return nil, err
	}

	return &light, nil
}

// DelLight -
func (c *Client) DelLight(lightID string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/states/%s", c.HostURL, lightID), nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)

	if err != nil {
		return err
	}
	return nil
}
func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	req.Header.Set("Authorization", c.Token)

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	return body, err
}
