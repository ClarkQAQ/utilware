package codec

import (
	"encoding/json"
	"net/url"
	"strings"
)

func JsonMarshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func JsonUnmarshal(b []byte, v interface{}) error {
	return json.Unmarshal(b, &v)
}

func JsonUnmarshalData(b string) (v interface{}) {
	json.Unmarshal([]byte(b), &v)
	return
}

func JsonMarshalUnmarshal(v interface{}) (interface{}, error) {
	b, e := json.Marshal(v)
	if e != nil {
		return nil, e
	}
	var data interface{}
	e = json.Unmarshal(b, data)
	return data, e
}

func GetQueryJson(values url.Values) (string, error) {
	data := make(map[string]interface{})
	for k, v := range values {
		data[k] = strings.Join(v, "-")
	}
	b, e := JsonMarshal(data)
	return string(b), e
}
