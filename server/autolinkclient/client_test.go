package autolinkclient

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/mattermost/mattermost-server/v6/plugin/plugintest"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/mattermost/mattermost-plugin-autolink/server/autolink"
)

func TestRoundTripper(t *testing.T) {
	mockPluginAPI := &plugintest.API{}

	mockPluginAPI.On("PluginHTTP", mock.AnythingOfType("*http.Request")).Return(&http.Response{StatusCode: http.StatusOK})

	roundTripper := pluginAPIRoundTripper{api: mockPluginAPI}
	req, err := http.NewRequest("POST", "url", nil)
	require.Nil(t, err)
	resp, err := roundTripper.RoundTrip(req)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	mockPluginAPI2 := &plugintest.API{}
	mockPluginAPI2.On("PluginHTTP", mock.AnythingOfType("*http.Request")).Return(nil)

	roundTripper2 := pluginAPIRoundTripper{api: mockPluginAPI2}
	req2, err := http.NewRequest("POST", "url", nil)
	require.Nil(t, err)
	resp2, err := roundTripper2.RoundTrip(req2)
	require.Nil(t, resp2)
	require.Error(t, err)
}

func TestAddAutolinks(t *testing.T) {
	mockPluginAPI := &plugintest.API{}

	mockPluginAPI.On("PluginHTTP", mock.AnythingOfType("*http.Request")).Return(&http.Response{StatusCode: http.StatusOK, Body: http.NoBody})

	client := NewClientPlugin(mockPluginAPI)
	err := client.Add(autolink.Autolink{})
	require.Nil(t, err)
}

func TestAddAutolinksErr(t *testing.T) {
	mockPluginAPI := &plugintest.API{}

	mockPluginAPI.On("PluginHTTP", mock.AnythingOfType("*http.Request")).Return(nil)

	client := NewClientPlugin(mockPluginAPI)
	err := client.Add(autolink.Autolink{})
	require.Error(t, err)
}

func TestDeleteAutolinks(t *testing.T) {
	for _, tc := range []struct {
		name     string
		errFound bool
	}{
		{
			name: "delete the autolink",
		},
		{
			name:     "got error",
			errFound: true,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			mockPluginAPI := &plugintest.API{}

			if tc.errFound {
				mockPluginAPI.On("PluginHTTP", mock.AnythingOfType("*http.Request")).Return(nil)
			} else {
				mockPluginAPI.On("PluginHTTP", mock.AnythingOfType("*http.Request")).Return(&http.Response{StatusCode: http.StatusOK, Body: http.NoBody})
			}

			client := NewClientPlugin(mockPluginAPI)
			err := client.Delete("")

			if tc.errFound {
				require.Error(t, err)
			} else {
				require.Nil(t, err)
			}
		})
	}
}

func TestGetAutolinks(t *testing.T) {
	for _, tc := range []struct {
		name     string
		errFound bool
	}{
		{
			name: "get the autolink",
		},
		{
			name:     "got error",
			errFound: true,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			mockPluginAPI := &plugintest.API{}

			body := ioutil.NopCloser(strings.NewReader("{}"))

			if tc.errFound {
				mockPluginAPI.On("PluginHTTP", mock.AnythingOfType("*http.Request")).Return(nil)
			} else {
				mockPluginAPI.On("PluginHTTP", mock.AnythingOfType("*http.Request")).Return(&http.Response{StatusCode: http.StatusOK, Body: body})
			}

			client := NewClientPlugin(mockPluginAPI)
			_, err := client.Get("")

			if tc.errFound {
				require.Error(t, err)
			} else {
				require.Nil(t, err)
			}
		})
	}
}
