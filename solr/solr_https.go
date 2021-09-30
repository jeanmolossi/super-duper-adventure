package solr

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

func HTTPGet(httpUrl string) ([]byte, error) {

	r, err := http.Get(httpUrl)

	if err != nil {
		return nil, err
	}

	defer r.Body.Close()

	if err != nil {
		return nil, err
	}

	// read the response and check
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func HTTPPost(url string, headers [][]string, payload *[]byte) ([]byte, error) {
	// setup post client
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewReader(*payload))

	// add headers
	if len(headers) > 0 {
		for i := range headers {
			req.Header.Add(headers[i][0], headers[i][1])
		}
	}

	// perform request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if err != nil {
		return nil, err
	}

	// read response, check & return
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
