// Package rpc project rpc.go
package rpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// BtcRPC is request info.
type BtcRPC struct {
	URL  string // bitcoin full node endpoint url
	User string // rpcuser
	Pass string // rpcpassword
	View bool   // If true, the log is displayed.
}

// BtcRPCRequest is request parameters.
type BtcRPCRequest struct {
	// bitcoin rpc request format
	Jsonrpc string        `json:"jsonrpc,"`
	ID      string        `json:"id,"`
	Method  string        `json:"method,"`
	Params  []interface{} `json:"params,"`
}

// Response is response details.
type Response struct {
	Result interface{} `json:"result,"`
	Error  interface{} `json:"error,"`
	ID     string      `json:"id,"`
}

// Error is error details.
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// UnmarshalError converts the error response into an RPCError type.
func (res *Response) UnmarshalError() (*Error, error) {
	rerr := &Error{}
	if res.Error == nil {
		return nil, fmt.Errorf("Rpesponse Error is nil")
	}
	data, ok := res.Error.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("RpcResponse Error is not map[string]interface{}")
	}
	bs, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(bs, rerr)
	if err != nil {
		return nil, err
	}
	return rerr, nil
}

// UnmarshalResult converts the response into an result type.
func (res *Response) UnmarshalResult(result interface{}) error {
	if res.Result == nil {
		return fmt.Errorf("RpcResponse Result is nil")
	}
	var bs []byte
	var err error
	m, ok := res.Result.(map[string]interface{})
	if !ok {
		arr, ok := res.Result.([]interface{})
		if !ok {
			return fmt.Errorf("RpcResponse Result is neither map[string]interface{} nor []interface{}")
		}
		bs, err = json.Marshal(arr)
	} else {
		bs, err = json.Marshal(m)
	}
	if err != nil {
		return err
	}
	err = json.Unmarshal(bs, result)
	return err
}

// NewBtcRPC returns new BtcRPC.
func NewBtcRPC(url, user, pass string) *BtcRPC {
	return &BtcRPC{url, user, pass, false}
}

// Request requests server.
func (rpc *BtcRPC) Request(method string, params ...interface{}) (*Response, error) {
	res := &Response{}
	if len(params) == 0 {
		params = []interface{}{}
	}
	id := fmt.Sprintf("%d", time.Now().Unix())
	req := &BtcRPCRequest{"1.0", id, method, params}
	bs, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	rpc.log("%s\n", bs)
	client := &http.Client{}
	hreq, err := http.NewRequest("POST", rpc.URL, bytes.NewBuffer(bs))
	if err != nil {
		return nil, err
	}
	hreq.SetBasicAuth(rpc.User, rpc.Pass)
	hres, err := client.Do(hreq)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = hres.Body.Close()
		if err != nil {
			log.Printf("close error : %+v", err)
		}
	}()
	body, err := ioutil.ReadAll(hres.Body)
	if err != nil {
		return nil, err
	}
	rpc.log("%d, %s\n", hres.StatusCode, body)
	err = json.Unmarshal(body, res)
	if err != nil || hres.StatusCode != http.StatusOK || res.ID != id {
		return nil, fmt.Errorf("status:%v, error:%v, body:%s reqid:%v, resid:%v", hres.Status, err, body, id, res.ID)
	}
	return res, nil
}

func (rpc *BtcRPC) log(format string, v ...interface{}) {
	if rpc.View {
		log.Printf(format, v...)
	}
}

func main() {
	fmt.Println("Hello World")
}
