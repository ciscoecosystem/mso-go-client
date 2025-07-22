package tests

import (
	"encoding/json"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/ciscoecosystem/mso-go-client/client"
	"github.com/ciscoecosystem/mso-go-client/container"
	"github.com/stretchr/testify/assert"
)

func TestClientAuthenticate(t *testing.T) {

	client := GetTestClient()
	err := client.Authenticate()
	if err != nil {
		t.Error(err)
	}

	fmt.Printf("err is %v", err)

	if client.AuthToken.Token == "{}" {
		t.Error("Token is empty")
	}

	fmt.Printf("Got Token %v", client.AuthToken.Token)
}

func GetTestClient() *client.Client {
	return client.GetClient("https://173.36.219.193", "admin", client.Password("ins3965!ins3965!"), client.Insecure(true))
}

func TestParallelGetSchemas(t *testing.T) {
	cl := GetTestClient()
	err := cl.Authenticate()
	if err != nil {
		t.Error(err)
	}
	schId := "6878807a072d2d88bec9b3b3" // Test_Schema
	schUrl := "api/v1/schemas/" + schId
	_, err = cl.GetViaURL(schUrl)

	assert := assert.New(t)
	assert.Equal(err, nil)

	numRequests := 6
	resps := make(map[int]*container.Container)
	errs := []error{}

	numObjs := 100
	numBatches := numObjs / numRequests

	fmt.Printf("Requesting %v objects in %v batches in %v requests per batch", numObjs, numBatches, numRequests)

	for b := 1; b <= numBatches; b++ {
		wgReqs := sync.WaitGroup{}
		// Create the workers
		for w := 1; w <= numRequests; w++ {
			wgReqs.Add(numRequests)
			go func(reqN int) {
				defer wgReqs.Done()
				var err error
				resps[reqN], err = cl.GetViaURL(schUrl)
				fmt.Printf("Batch: %v Request: %v GetViaURL err = [%v]\n", b, reqN, err)
				errs = append(errs, err)
			}(w)
		}
		wgReqs.Wait()
		//		time.Sleep(2 * time.Second)
		time.Sleep(200000000) // 2*10^8 nano seconds = 200 ms
	}
	assert.Equal(err, nil)
	fmt.Printf("len(resps) = %v\n", len(resps))
}

func TestParallelGetSchemasMso(t *testing.T) {
	cl := GetTestClient()
	err := cl.Authenticate()
	if err != nil {
		t.Error(err)
	}
	schId := "6878807a072d2d88bec9b3b3" // for Test_Schema
	schUrl := "mso/api/v1/schemas/" + schId
	_, err = cl.GetViaURL(schUrl)

	assert := assert.New(t)
	assert.Equal(err, nil)

	numRequests := 6
	resps := make(map[int]*container.Container)
	errs := []error{}

	numObjs := 120
	numBatches := numObjs / numRequests

	fmt.Printf("Requesting %v objects in %v batches in %v requests per batch", numObjs, numBatches, numRequests)

	for b := 1; b <= numBatches; b++ {
		wgReqs := sync.WaitGroup{}
		// Create the workers
		for w := 1; w <= numRequests; w++ {
			wgReqs.Add(1)
			go func(reqN int) {
				defer wgReqs.Done()
				var err error
				resps[reqN], err = cl.GetViaURL(schUrl)
				fmt.Printf("Batch: %v Request: %v GetViaURL err = [%v]\n", b, reqN, err)
				errs = append(errs, err)
			}(w)
		}
		wgReqs.Wait()
		//		time.Sleep(2 * time.Second)
		time.Sleep(200000000) // 2*10^8 nano seconds = 200 ms
	}
	assert.Equal(err, nil)
	fmt.Printf("len(resps) = %v\n", len(resps))
}

func TestParallelPatchSchemas(t *testing.T) {
	cl := GetTestClient()
	err := cl.Authenticate()
	if err != nil {
		t.Error(err)
	}

	return

	schemaID := "6878807a072d2d88bec9b3b3"
	schUrl := "api/v1/schemas/" + schemaID

	assert := assert.New(t)

	_, err = cl.GetViaURL(schUrl)

	numBatches := 3
	numRequests := 3

	for b := 0; b < numBatches; b++ {
		wgReqs := sync.WaitGroup{}
		// Create the workers
		for w := 1; w <= numRequests; w++ {
			wgReqs.Add(1)
			bdNum := b*numRequests + w
			go func(bdN int) {
				defer wgReqs.Done()
				var err error
				bdName := fmt.Sprintf("BD%v", bdN)
				desc := fmt.Sprintf("new descr %v", 300+bdN)
				err = patchBDDescr(cl, schemaID, "Tmpl1", bdName, desc)
				assert.Equal(err, nil)
				_, err = cl.GetViaURL(schUrl)
				fmt.Printf("Batch: %v Request: %v GetViaURL err = [%v]\n", b, w, err)
			}(bdNum)
		}
		wgReqs.Wait()
		//		time.Sleep(2 * time.Second)
		time.Sleep(200000000) // 2*10^8 nano seconds = 200 ms
	}
	assert.Equal(err, nil)
}

func doPatchRequest(msoClient *client.Client, path string, payloadCon *container.Container) error {
	req, err := msoClient.MakeRestRequest("PATCH", path, payloadCon, true)
	if err != nil {
		return err
	}

	cont, _, err := msoClient.Do(req)
	if err != nil {
		return err
	}

	err = client.CheckForErrors(cont, "PATCH")
	if err != nil {
		return err
	}

	return nil
}

func addPatchPayloadToContainer(payloadContainer *container.Container, op, path string, value interface{}) error {

	payloadMap := map[string]interface{}{"op": op, "path": path, "value": value}

	payload, err := json.Marshal(payloadMap)
	if err != nil {
		return err
	}

	jsonContainer, err := container.ParseJSON([]byte(payload))
	if err != nil {
		return err
	}

	err = payloadContainer.ArrayAppend(jsonContainer.Data())
	if err != nil {
		return err
	}

	return nil
}

func patchBDDescr(cl *client.Client, schemaID string, templateName string, bdName string, desc string) error {
	basePath := fmt.Sprintf("/templates/%s/bds/%s", templateName, bdName)
	payloadCon := container.New()
	payloadCon.Array()

	err := addPatchPayloadToContainer(payloadCon, "replace", fmt.Sprintf("%s/description", basePath), desc)
	if err != nil {
		return err
	}

	err = doPatchRequest(cl, fmt.Sprintf("api/v1/schemas/%s", schemaID), payloadCon)
	if err != nil {
		return err
	}
	return nil
}
