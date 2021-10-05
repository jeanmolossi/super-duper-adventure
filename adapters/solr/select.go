package solr

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Query struct {
	Params     URLParamMap
	Rows       int
	Start      int
	Sort       string
	DefType    string
	Debug      bool
	OmitHeader bool
}

func (q *Query) String() string {
	var s []string

	if len(q.Params) > 0 {
		s = append(s, EncodeURLParamMap(&q.Params))
	}

	if q.Rows != 0 {
		s = append(s, fmt.Sprintf("rows=%d", q.Rows))
	}

	if q.Start != 0 {
		s = append(s, fmt.Sprintf("start=%d", q.Start))
	}

	if q.Sort != "" {
		s = append(s, fmt.Sprintf("sort=%s", q.Sort))
	}

	if q.DefType != "" {
		s = append(s, fmt.Sprintf("defType=%s", q.DefType))
	}

	if q.Debug {
		s = append(s, fmt.Sprintf("debugQuery=true"))
	}

	if q.OmitHeader {
		s = append(s, fmt.Sprintf("omitHeader=true"))
	}

	return strings.Join(s, "&")
}

type Connection struct {
	URL     string
	Version []int
}

func SolrSelectString(c *Connection, q string, handlerName string) string {
	return fmt.Sprintf("%s/%s?wt=json&%s", c.URL, handlerName, q)
}

func SolrUpdateString(c *Connection, commit bool) string {
	s := fmt.Sprintf("%s/update", c.URL)
	if commit {
		return fmt.Sprintf("%s?commit=true", s)
	}
	return s
}

func BytesToJSON(b *[]byte) (*interface{}, error) {
	var container interface{}
	err := json.Unmarshal(*b, &container)

	if err != nil {
		return nil, err
	}

	return &container, nil
}

func JSONToBytes(m map[string]interface{}) (*[]byte, error) {
	b, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	return &b, nil
}

func BuildResponse(j *interface{}) (*SelectResponse, error) {

	// look for a response element, bail if not present
	responseRoot := (*j).(map[string]interface{})
	response := responseRoot["response"]
	if response == nil {
		return nil, fmt.Errorf("Supplied interface appears invalid (missing response)")
	}

	// begin Response creation
	r := SelectResponse{}

	// do status & qtime, if possible
	rHeader := (*j).(map[string]interface{})["responseHeader"].(map[string]interface{})
	if rHeader != nil {
		r.Status = int(rHeader["status"].(float64))
		r.QTime = int(rHeader["QTime"].(float64))
	}

	// now do docs, if they exist in the response
	docs := response.(map[string]interface{})["docs"].([]interface{})
	if docs != nil {
		// the total amount of results, irrespective of the amount returned in the response
		numFound := int(response.(map[string]interface{})["numFound"].(float64))

		// and the amount actually returned
		numResults := len(docs)

		coll := DocumentCollection{}
		coll.NumFound = numFound

		var ds []Document

		for i := 0; i < numResults; i++ {
			ds = append(ds, Document{docs[i].(map[string]interface{})})
		}

		coll.Collection = ds
		r.Results = &coll
	}

	return &r, nil
}

func (c *Connection) CustomSelect(q *Query, handlerName string) (*SelectResponse, error) {
	body, err := HTTPGet(SolrSelectString(c, q.String(), handlerName))

	if err != nil {
		return nil, err
	}

	r, err := SelectResponseFromHTTPResponse(body)

	if err != nil {
		return nil, err
	}

	return r, nil
}

func (c *Connection) CustomSelectRaw(q string, handlerName string) (*SelectResponse, error) {
	body, err := HTTPGet(SolrSelectString(c, q, handlerName))

	if err != nil {
		return nil, err
	}

	r, err := SelectResponseFromHTTPResponse(body)

	if err != nil {
		return nil, err
	}

	return r, nil
}

func (c *Connection) SelectRaw(q string) (*SelectResponse, error) {
	resp, err := c.CustomSelectRaw(q, "select")
	return resp, err
}

func (c *Connection) Update(m map[string]interface{}, commit bool) (*UpdateResponse, error) {

	// encode "json" to a byte array & check
	payload, err := JSONToBytes(m)
	if err != nil {
		return nil, err
	}

	// perform request
	resp, err := HTTPPost(
		SolrUpdateString(c, commit),
		[][]string{{"Content-Type", "application/json"}},
		payload)

	if err != nil {
		return nil, err
	}

	// decode the response & check
	decoded, err := BytesToJSON(&resp)
	if err != nil {
		return nil, err
	}

	hasErr, report := SolrErrorResponse((*decoded).(map[string]interface{}))
	if hasErr {
		return nil, fmt.Errorf(fmt.Sprintf("%s", *report))
	}

	return &UpdateResponse{true}, nil
}
