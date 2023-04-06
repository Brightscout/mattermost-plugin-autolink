package autolinkclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/brightscout/mattermost-plugin-autolink/server/autolink"
)

const (
	autolinkPluginID       = "mattermost-autolink"
	AutolinkNameQueryParam = "autolinkName"
)

type PluginAPI interface {
	PluginHTTP(*http.Request) *http.Response
}

type Client struct {
	http.Client
}

type pluginAPIRoundTripper struct {
	api PluginAPI
}

func (p *pluginAPIRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	resp := p.api.PluginHTTP(req)
	if resp == nil {
		return nil, fmt.Errorf("failed to make interplugin request")
	}
	return resp, nil
}

func NewClientPlugin(api PluginAPI) *Client {
	client := &Client{}
	client.Transport = &pluginAPIRoundTripper{api}
	return client
}

func (c *Client) Add(links ...autolink.Autolink) error {
	for _, link := range links {
		linkBytes, err := json.Marshal(link)
		if err != nil {
			return err
		}

		req, err := http.NewRequest(http.MethodPost, "/"+autolinkPluginID+"/api/v1/link", bytes.NewReader(linkBytes))
		if err != nil {
			return err
		}

		resp, err := c.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			respBody, _ := ioutil.ReadAll(resp.Body)
			return fmt.Errorf("unable to install autolink. Error: %v, %v", resp.StatusCode, string(respBody))
		}
	}

	return nil
}

func (c *Client) Delete(links ...string) error {
	for _, link := range links {
		queryParams := url.Values{
			AutolinkNameQueryParam: {link},
		}

		req, err := http.NewRequest(http.MethodDelete, "/"+autolinkPluginID+"/api/v1/link", nil)
		if err != nil {
			return err
		}

		req.URL.RawQuery = queryParams.Encode()

		resp, err := c.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			respBody, _ := ioutil.ReadAll(resp.Body)
			return fmt.Errorf("unable to install autolink. Error: %v, %v", resp.StatusCode, string(respBody))
		}
	}

	return nil
}

func (c *Client) Get(autolinkName string) ([]autolink.Autolink, error) {
	queryParams := url.Values{
		AutolinkNameQueryParam: {autolinkName},
	}

	req, err := http.NewRequest(http.MethodGet, "/"+autolinkPluginID+"/api/v1/link", nil)
	if err != nil {
		return nil, err
	}

	req.URL.RawQuery = queryParams.Encode()

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unable to install autolink. Error: %v, %v", resp.StatusCode, string(respBody))
	}

	var response []autolink.Autolink
	if err = json.Unmarshal(respBody, &response); err != nil {
		return nil, err
	}

	return response, nil
}
