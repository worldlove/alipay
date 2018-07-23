package alipay

import (
	"encoding/json"
	"net/url"
)

func copyUrlValues(src url.Values) url.Values {
	var dst = make(url.Values)
	for k, v := range src {
		dst[k] = v
	}
	return dst
}
func copyValuesToMap(src url.Values) map[string]string {
	var dst = make(map[string]string)
	for k := range src {
		dst[k], _ = url.QueryUnescape(src.Get(k))
	}
	return dst
}

func copyMap(src map[string]string) map[string]string {
	var dst = make(map[string]string)
	for k, v := range src {
		dst[k] = v
	}
	return dst
}
func copyMapToJSON(src map[string]string) []byte {
	str, _ := json.Marshal(src)
	return str
}

func copyMapToValues(src map[string]string) url.Values {
	var dst = make(url.Values)
	for k, v := range src {
		dst.Set(k, v)
	}
	return dst
}
