package homeauto

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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
// client = The http client used for the calls
func NewClient(host, token string, client *http.Client) *Client {
	return &Client{
		HTTPClient: client,
		HostURL:    host,
		Token:      fmt.Sprintf("Bearer %s", token),
	}
}

//GetLight is used to get the currest state of a light
//lightID = the ID of a light (eg. light.virtual_light_10 )
func GetLight(lightID string, client Client) (LightItem, error) {
	url := fmt.Sprintf("%s/api/states/%s", client.HostURL, lightID)
	body, err := doRequest(client, http.MethodGet, url, nil)
	var light LightItem
	err = json.Unmarshal(body, &light)
	if err != nil {
		return LightItem{}, err
	}

	return light, nil
}

//StartLight runs a POST request on the api to create a light (is light already exists it updates it)
// Takes in the state you want the light to be, returns the state the light ends up in
func StartLight(lightItem LightItem, client Client) error {
	rb, err := json.Marshal(lightItem)
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%s/api/states/%s", client.HostURL, lightItem.EntityID)
	_, err = doRequest(client, http.MethodPost, url, rb)
	if err != nil {
		return err
	}
	return nil
}

// DeleteLight will delete the light and return an error if the delete fails
func DeleteLight(lightID string, client Client) error {
	url := fmt.Sprintf("%s/api/states/%s", client.HostURL, lightID)
	_, err := doRequest(client, http.MethodDelete, url, nil)
	if err != nil {
		return err
	}
	return nil
}
func doRequest(client Client, method, url string, rb []byte) ([]byte, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(rb))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", client.Token)
	res, err := client.HTTPClient.Do(req)
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
	return body, nil
}
