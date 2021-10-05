package solr

type Document struct {
	Fields map[string]interface{}
}

type DocumentCollection struct {
	Collection []Document
	NumFacets  int
	NumFound   int
	Start      int
}

func (d *DocumentCollection) Get(i int) *Document {
	return &d.Collection[i]
}

func (d *DocumentCollection) Len() int {
	return len(d.Collection)
}

func (document Document) Field(field string) interface{} {
	r, _ := document.Fields[field]
	return r
}

func (document Document) Doc() map[string]interface{} {
	return document.Fields
}
