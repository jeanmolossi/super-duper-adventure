package solr

import (
	"fmt"
)

type SelectResponse struct {
	Results *DocumentCollection
	Status  int
	QTime   int
	// TODO: Debug info
}

func (r SelectResponse) String() string {
	return fmt.Sprintf("SelectResponse: %d Results, Status: %d, QTime: %d", r.Results.Len(), r.Status, r.QTime)
}

func SelectResponseFromHTTPResponse(b []byte) (*SelectResponse, error) {
	j, err := BytesToJSON(&b)
	if err != nil {
		return nil, err
	}

	resp, err := BuildResponse(j)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

type ErrorResponse struct {
	Message string
	Status  int
}

func (r ErrorResponse) String() string {
	return fmt.Sprintf("Solr Error: [code: %d, msg: \"%s\"]", r.Status, r.Message)
}

func SolrErrorResponse(m map[string]interface{}) (bool, *ErrorResponse) {
	// check for existance of "error" key
	if _, found := m["error"]; found {
		err := m["error"].(map[string]interface{})
		return true, &ErrorResponse{
			Message: err["msg"].(string),
			Status:  int(err["code"].(float64)),
		}
	}
	return false, nil
}

type UpdateResponse struct {
	Success bool
}

func (r UpdateResponse) String() string {
	if r.Success {
		return fmt.Sprintf("UpdateResponse: OK")
	}
	return fmt.Sprintf("UpdateResponse: FAIL")
}

func chunk(s []interface{}, sz int) [][]interface{} {
	var r [][]interface{}
	j := len(s)
	for i := 0; i < j; i += sz {
		r = append(r, s[i:i+sz])
	}
	return r
}
