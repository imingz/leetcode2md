package request

import (
	"io"
	"net/http"
	"strings"
)

func request(payload *strings.Reader) ([]byte, error) {
	// 请求
	url := "https://leetcode.cn/graphql/"
	method := "POST"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	return io.ReadAll(res.Body)
}
