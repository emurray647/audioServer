package system_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/emurray647/audioServer/internal/model"
	"github.com/stretchr/testify/assert"
)

const (
	urlBase = "http://audioapi_test:8080"
)

func doRequest(method, path string, reader io.Reader) (io.ReadCloser, error) {

	req, err := http.NewRequest(method, fmt.Sprintf("%s%s", urlBase, path), reader)
	if err != nil {
		return nil, err
	}
	req.Close = true

	fmt.Println(req)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
}

func TestAPI(t *testing.T) {
	var statusMessage model.StatusMessage
	var listResult []model.AudioFileDetails
	var infoResult model.AudioFileDetails

	buffer, err := os.ReadFile("/samples/cantina.wav")
	checkError(err)

	// upload a single file
	reader := bytes.NewReader(buffer)
	body, err := doRequest("POST", "/files?name=cantina.wav", reader)
	checkError(err)
	err = json.NewDecoder(body).Decode(&statusMessage)
	checkError(err)
	assert.Equal(t, 200, statusMessage.StatusCode)
	assert.Equal(t, true, statusMessage.Success)

	// doing the same request again should give us a 409
	body, err = doRequest("POST", "/files?name=cantina.wav", reader)
	checkError(err)
	err = json.NewDecoder(body).Decode(&statusMessage)
	body.Close()
	checkError(err)
	assert.Equal(t, 409, statusMessage.StatusCode)
	assert.Equal(t, false, statusMessage.Success)

	// try uploading nonsense: should give 400
	reader = bytes.NewReader([]byte{0, 1, 2, 3, 4, 5, 6, 7})
	body, err = doRequest("POST", "/files?name=nonsense.wav", reader)
	checkError(err)
	err = json.NewDecoder(body).Decode(&statusMessage)
	checkError(err)
	body.Close()
	assert.Equal(t, 400, statusMessage.StatusCode)
	assert.Equal(t, false, statusMessage.Success)

	// list ing should give use one result
	body, err = doRequest("GET", "/list", nil)
	checkError(err)
	err = json.NewDecoder(body).Decode(&listResult)
	checkError(err)
	assert.Equal(t, 1, len(listResult))
	assert.Equal(t, "cantina.wav", listResult[0].Name)

	// now we can delete
	body, err = doRequest("DELETE", "/files?name=cantina.wav", nil)
	err = json.NewDecoder(body).Decode(&statusMessage)
	checkError(err)
	assert.Equal(t, 200, statusMessage.StatusCode)
	assert.Equal(t, true, statusMessage.Success)

	// now we should see zero files
	body, err = doRequest("GET", "/list", nil)
	checkError(err)
	err = json.NewDecoder(body).Decode(&listResult)
	checkError(err)
	assert.Equal(t, 0, len(listResult))

	// upload three files
	buffer, err = os.ReadFile("/samples/cantina.wav")
	checkError(err)
	reader = bytes.NewReader(buffer)
	body, err = doRequest("POST", "/files?name=cantina.wav", reader)
	checkError(err)
	err = json.NewDecoder(body).Decode(&statusMessage)
	checkError(err)
	assert.Equal(t, 200, statusMessage.StatusCode)
	assert.Equal(t, true, statusMessage.Success)

	buffer, err = os.ReadFile("/samples/gettysburg.wav")
	checkError(err)
	reader = bytes.NewReader(buffer)
	body, err = doRequest("POST", "/files?name=gettysburg.wav", reader)
	checkError(err)
	err = json.NewDecoder(body).Decode(&statusMessage)
	checkError(err)
	assert.Equal(t, 200, statusMessage.StatusCode)
	assert.Equal(t, true, statusMessage.Success)

	buffer, err = os.ReadFile("/samples/tchaiv.mp3")
	checkError(err)
	reader = bytes.NewReader(buffer)
	body, err = doRequest("POST", "/files?name=tchaiv.mp3", reader)
	checkError(err)
	err = json.NewDecoder(body).Decode(&statusMessage)
	checkError(err)
	assert.Equal(t, 200, statusMessage.StatusCode)
	assert.Equal(t, true, statusMessage.Success)

	// there should be 3 when we list
	body, err = doRequest("GET", "/list", nil)
	checkError(err)
	err = json.NewDecoder(body).Decode(&listResult)
	checkError(err)
	assert.Equal(t, 3, len(listResult))

	// check info
	body, err = doRequest("GET", "/info?name=gettysburg.wav", nil)
	checkError(err)
	err = json.NewDecoder(body).Decode(&infoResult)
	checkError(err)
	assert.Equal(t, "gettysburg.wav", infoResult.Name)
	assert.Equal(t, 22050, *infoResult.SampleRate)
	assert.Equal(t, 10.0039, *infoResult.Duration)

	// check some list filters
	body, err = doRequest("GET", "/list?format=wav", nil)
	checkError(err)
	err = json.NewDecoder(body).Decode(&listResult)
	checkError(err)
	assert.Equal(t, 2, len(listResult))
	for _, res := range listResult {
		assert.Equal(t, "wav", res.Format)
	}

	body, err = doRequest("GET", "/list?maxduration=300", nil)
	checkError(err)
	err = json.NewDecoder(body).Decode(&listResult)
	checkError(err)
	assert.Equal(t, 2, len(listResult))
	for _, res := range listResult {
		assert.Less(t, *res.Duration, 300.0)
	}

}
