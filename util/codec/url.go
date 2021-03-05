package codec

import "net/url"

func EsURL(data string) string {
	return url.QueryEscape(data)
}

func UnURL(data string) (string, error) {
	return url.QueryUnescape(data)
}
