package solr

import (
	"fmt"
	"net/url"
)

type URLParamMap map[string][]string

func EncodeURLParamMap(m *URLParamMap) string {
	parameters := url.Values{}

	for k, v := range *m {
		l := len(v)
		for x := 0; x < l; x++ {
			parameters.Add(k, v[x])
		}
	}

	queryEncoded := parameters.Encode()
	return queryEncoded
}

func (c *Connection) Select(q *Query) (*SelectResponse, error) {
	resp, err := c.CustomSelect(q, "select")
	return resp, err
}

func Init(host string, port int, core string) (*Connection, error) {

	if len(host) == 0 {
		return nil, fmt.Errorf("Invalid hostname (must be length >= 1)")
	}

	if port <= 0 || port > 65535 {
		return nil, fmt.Errorf("Invalid port (must be 1..65535")
	}

	mountedUrl := fmt.Sprintf("http://%s:%d/solr/%s", host, port, core)
	return &Connection{URL: mountedUrl}, nil
}
