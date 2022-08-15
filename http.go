package gopkg

import (
	"encoding/json"
	"net/http"
	"time"
)

// GetJSON 请求url，解析json数据到target
func GetJSON(url string, timeout time.Duration, target interface{}) error {
	var client = &http.Client{Timeout: timeout * time.Second}
	r, err := client.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(target)
}
