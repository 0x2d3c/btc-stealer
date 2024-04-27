package common

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync/atomic"
)

var ethReqTemp = `{
    "jsonrpc": "2.0",
    "method": "eth_getBalance",
    "params": [
        "%s",
        "latest"
    ],
    "id": 1
}`

type ethResp struct {
	Result string `json:"result"`
}

var ethCount int64

func HttpETH(addr string) string {
	req, err := http.NewRequest(http.MethodPost, config.ETHGW, strings.NewReader(fmt.Sprintf(ethReqTemp, addr)))
	if err != nil {
		return ""
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := (&http.Client{}).Do(req)
	if err != nil {
		return err.Error()
	}
	defer res.Body.Close()

	bs, err := io.ReadAll(res.Body)
	if err != nil {
		return ""
	}

	var resp ethResp
	if err = json.Unmarshal(bs, &resp); err != nil {
		return ""
	}

	atomic.AddInt64(&ethCount, 1)

	return resp.Result
}

func ETHCount() int64 {
	return ethCount
}

type btcResp struct {
	ChainStats struct {
		FundedTxoSum int64 `json:"funded_txo_sum"`
	} `json:"chain_stats"`
}

var btcCount int64

func HttpBTC(addr string) (bool, error) {
	res, err := http.Get(config.BTCGW + addr)
	if err != nil {
		return false, err
	}
	defer res.Body.Close()

	bs, err := io.ReadAll(res.Body)
	if err != nil {
		return false, nil
	}

	var resp btcResp
	if err = json.Unmarshal(bs, &resp); err != nil {
		return false, nil
	}

	atomic.AddInt64(&btcCount, 1)

	return resp.ChainStats.FundedTxoSum > 0, nil
}

func BTCCount() int64 {
	return btcCount
}
