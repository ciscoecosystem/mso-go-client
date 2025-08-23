package client

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/ciscoecosystem/mso-go-client/container"
	"github.com/ciscoecosystem/mso-go-client/models"
)

const (
	CACHE_TIMEOUT = 60 // 60 seconds
)

type msoApi struct {
	readTs  time.Time
	resp    *container.Container
	writeTs time.Time
}

var msoApiCache map[string]msoApi
var muApiCache sync.RWMutex // mutex lock for upating the map

// init of the package
func init() {
	msoApiCache = make(map[string]msoApi)
}

// getFromCache: check the API cache and return the stored resp
// if it is with the timeout
func (c *Client) getFromCache(endpoint string) *container.Container {
	defer muApiCache.RUnlock()
	muApiCache.RLock()
	updEndpoint := strings.Replace(endpoint, "mso/", "", 1)
	if api, ok := msoApiCache[updEndpoint]; ok {
		curTs := time.Now()
		rDiff := curTs.Sub(api.readTs)
		wDiff := curTs.Sub(api.writeTs)
		log.Printf("[DEBUG] getFromCache readTs %v writeTs: %v rDiff %v wDiff %v\n", api.readTs, api.writeTs, rDiff.Seconds(), wDiff.Seconds())
		if rDiff.Seconds() >= CACHE_TIMEOUT || wDiff.Seconds() <= CACHE_TIMEOUT {
			return nil
		}
		log.Printf("[DEBUG] Found GET response in cache for schema endpoint: %v\n", updEndpoint)
		return api.resp
	}
	return nil
}

// storeInCache: store the given response in the API cache
func (c *Client) storeInCache(endpoint string, resp *container.Container) {
	updEndpoint := strings.Replace(endpoint, "mso/", "", 1)
	var re = regexp.MustCompile(`^api/v1/schemas/(.*)$`)
	matches := re.FindStringSubmatch(updEndpoint)

	if len(matches) != 2 {
		return
	}

	defer muApiCache.Unlock()

	muApiCache.Lock()
	if api, ok := msoApiCache[updEndpoint]; ok {
		curTs := time.Now()
		wDiff := curTs.Sub(api.writeTs)
		if wDiff.Seconds() <= CACHE_TIMEOUT {
			log.Printf("[DEBUG] Skip storing endpoint %v due to recent writeTs: %v\n", updEndpoint, api.writeTs)
			return
		}
	}

	api := msoApi{
		readTs:  time.Now(),
		resp:    resp,
		writeTs: time.Now().Add(-180 * time.Second),
	}

	log.Printf("[DEBUG] Caching GET endpoint:: %s readTs %v writeTs %v", updEndpoint, api.readTs, api.writeTs)
	msoApiCache[updEndpoint] = api
}

// invalidateCache: invalidate the cache
func (c *Client) updateCacheForWrite(endpoint string) {
	updEndpoint := strings.Replace(endpoint, "mso/", "", 1)
	var re = regexp.MustCompile(`^api/v1/schemas/(.*)(\?)?`)
	matches := re.FindStringSubmatch(updEndpoint)
	if len(matches) != 2 && len(matches) != 3 {
		return
	}

	defer muApiCache.Unlock()
	schEndPoint := "api/v1/schemas/" + matches[1]
	muApiCache.Lock()
	if api, ok := msoApiCache[schEndPoint]; ok {
		api.writeTs = time.Now()
		api.resp = nil
		msoApiCache[schEndPoint] = api
		log.Printf("[DEBUG] Update writeTs %v in cache for schema endpoint: %v\n", api.writeTs, schEndPoint)
	}
}

func (c *Client) GetViaURL(endpoint string) (*container.Container, error) {
	cobj := c.getFromCache(endpoint)

	if cobj != nil {
		c.storeInCache(endpoint, cobj)
		return cobj, nil
	}
	req, err := c.MakeRestRequest("GET", endpoint, nil, true)

	if err != nil {
		return nil, err
	}

	obj, _, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	if obj == nil {
		return nil, errors.New("Empty response body")
	}
	err = CheckForErrors(obj, "GET")

	if err != nil {
		return obj, err
	}

	c.storeInCache(endpoint, obj)

	return obj, nil
}

func (c *Client) GetPlatform() string {
	return c.platform
}

func (c *Client) Put(endpoint string, obj models.Model) (*container.Container, error) {
	jsonPayload, err := c.PrepareModel(obj)

	if err != nil {
		return nil, err
	}
	req, err := c.MakeRestRequest("PUT", endpoint, jsonPayload, true)
	if err != nil {
		return nil, err
	}

	c.Mutex.Lock()
	cont, _, err := c.Do(req)
	c.Mutex.Unlock()
	if err != nil {
		return nil, err
	}

	return cont, CheckForErrors(cont, "PUT")
}

func (c *Client) Save(endpoint string, obj models.Model) (*container.Container, error) {

	jsonPayload, err := c.PrepareModel(obj)

	if err != nil {
		return nil, err
	}
	req, err := c.MakeRestRequest("POST", endpoint, jsonPayload, true)
	if err != nil {
		return nil, err
	}

	cont, _, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	return cont, CheckForErrors(cont, "POST")
}

// CheckForErrors parses the response and checks of there is an error attribute in the response
func CheckForErrors(cont *container.Container, method string) error {

	if cont.Exists("code") && cont.Exists("message") {
		return errors.New(fmt.Sprintf("%s%s", cont.S("message"), cont.S("info")))
	} else if cont.Exists("error") {
		return errors.New(fmt.Sprintf("%s %s", models.StripQuotes(cont.S("error").String()), models.StripQuotes(cont.S("error_code").String())))
	}
	return nil
}

func (c *Client) DeletebyId(url string) error {

	req, err := c.MakeRestRequest("DELETE", url, nil, true)
	if err != nil {
		return err
	}

	_, resp, err1 := c.Do(req)
	if err1 != nil {
		return err1
	}
	if resp != nil {
		if resp.StatusCode == 204 || resp.StatusCode == 200 {
			return nil
		} else {
			return fmt.Errorf("Unable to delete the object")
		}
	}

	return nil
}

func (c *Client) PatchbyID(endpoint string, objList ...models.Model) (*container.Container, error) {

	contJs := container.New()
	contJs.Array()
	for _, obj := range objList {
		jsonPayload, err := c.PrepareModel(obj)
		if err != nil {
			return nil, err
		}
		contJs.ArrayAppend(jsonPayload.Data())

	}
	log.Printf("[DEBUG] Patch Request Container: %v\n", contJs)
	// URL encoding
	baseUrl, _ := url.Parse(endpoint)
	qs := url.Values{}
	qs.Add("validate", "false")
	baseUrl.RawQuery = qs.Encode()

	req, err := c.MakeRestRequest("PATCH", baseUrl.String(), contJs, true)
	if err != nil {
		return nil, err
	}

	c.Mutex.Lock()
	cont, _, err := c.Do(req)
	c.Mutex.Unlock()
	if err != nil {
		return nil, err
	}

	return cont, CheckForErrors(cont, "PATCH")
}

func (c *Client) PrepareModel(obj models.Model) (*container.Container, error) {
	con, err := obj.ToMap()
	if err != nil {
		return nil, err
	}

	payload := &container.Container{}
	if err != nil {
		return nil, err
	}

	for key, value := range con {
		payload.Set(value, key)
	}
	return payload, nil
}
