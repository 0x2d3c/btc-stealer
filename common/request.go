package common

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func httpClient() *http.Client {
	var client http.Client
	if config.Proxy.Enable {
		proxy, err := url.Parse(config.Proxy.Address)
		if err != nil {
			panic("proxy err")
		}
		client.Transport = &http.Transport{Proxy: http.ProxyURL(proxy)}
	}
	return &client
}

func HttpGetRequest(uri string, dst interface{}) {
	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		fmt.Println("request create error", err)
		return
	}

	res, err := httpClient().Do(req)
	if err != nil {
		fmt.Println("sending request:", err)
		return
	}
	defer res.Body.Close()

	bs, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("reading request:", err)
		return
	}

	if err = json.Unmarshal(bs, &dst); err != nil {
		fmt.Println("unmarshal resp", err)
		return
	}
}
