package client

import (
	"net/url"
	"testing"
)

var TestBaseUrls = [...]string{
	"https://ndo.host.cisco",
	"https://ndo.host.cisco/",
	"https://ndo.host.cisco//",
	"https://ndo.host.cisco///",
}

func AssertFullUrl(t *testing.T, baseUrl string, platform string, method string, path string, expected string) {
	url, err := url.Parse(baseUrl)
	if err != nil {
		t.Fatal(err)
	}
	ndclient := &Client{
		BaseURL:  url,
		platform: platform,
	}

	actual, err := ndclient.MakeFullUrl(method, path)
	if actual != expected || err != nil {
		t.Errorf(`MakeFullUrl("%s", "%s") %s = %q, %v, expected %#q`, method, path, platform, actual, err, expected)
	}
}

func TestMakeFullUrl_Login(t *testing.T) {
	expected := "https://ndo.host.cisco/login"
	for _, baseUrl := range TestBaseUrls {
		AssertFullUrl(t, baseUrl, "nd", "POST", "/login", expected)
		AssertFullUrl(t, baseUrl, "mso", "POST", "/login", expected)
	}
}

func TestMakeFullUrl_Get(t *testing.T) {
	expected := "https://ndo.host.cisco/templates/123"
	expected_nd := "https://ndo.host.cisco/mso/templates/123"
	paths := [...]string {
		"templates/123",
		"/templates/123",
		"///templates/123",
	}
	for _, baseUrl := range TestBaseUrls {
		for _, path := range paths {
			AssertFullUrl(t, baseUrl, "nd", "GET", path, expected_nd)
			AssertFullUrl(t, baseUrl, "mso", "GET", path, expected)
		}
	}
}

func TestMakeFullUrl_Patch(t *testing.T) {
	expected := "https://ndo.host.cisco/templates/123?validate=false"
	expected_nd := "https://ndo.host.cisco/mso/templates/123?validate=false"
	path := "/templates/123"
	for _, baseUrl := range TestBaseUrls {
		AssertFullUrl(t, baseUrl, "nd", "PATCH", path, expected_nd)
		AssertFullUrl(t, baseUrl, "mso", "PATCH", path, expected)
	}
}

func TestMakeFullUrl_PatchExtraQuery(t *testing.T) {
	expected := "https://ndo.host.cisco/templates/123?extra=query&validate=false"
	expected_nd := "https://ndo.host.cisco/mso/templates/123?extra=query&validate=false"
	path := "templates/123?extra=query"
	for _, baseUrl := range TestBaseUrls {
		AssertFullUrl(t, baseUrl, "nd", "PATCH", path, expected_nd)
		AssertFullUrl(t, baseUrl, "mso", "PATCH", path, expected)
	}
}
