package homeauto

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// Client is the struct that terraform passes as the identity of the server
type Client struct {
	HostURL    string
	HTTPClient *http.Client
	Token      string
}

//NewClient creates Client and makes sure that the Authorization is correct
// host = hosts url is http://127.0.0.1:8123
// token = The bearer token used for authentication
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

//GetLight is used to get the currest state of a light
//lightID = the ID of a light (eg. light.virtual_light_10 )
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

//StartLight runs a POST request on the api to create a light (is light already exists it updates it)
// Takes in the state you want the light to be, returns the state the light ends up in
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

// DelLight will delete the light and return an error if the delete fails
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

// doRequest is used to add authorization and send the request to the server
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
