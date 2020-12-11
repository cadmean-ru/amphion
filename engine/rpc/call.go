package rpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/cadmean-ru/amphion/engine"
	"io/ioutil"
	"net/http"
)

type FunctionCall struct {
	Arguments []interface{} `json:"args"`
	Auth      string        `json:"auth"`
}

type CallBuilder struct {
	fName           string
	call            FunctionCall
	successCallback func(res interface{})
	errorCallback   func(err Error)
}

func (b *CallBuilder) Args(args ...interface{}) *CallBuilder {
	b.call.Arguments = args
	return b
}

func (b *CallBuilder) Auth(auth string) *CallBuilder {
	b.call.Auth = auth
	return b
}

func (b *CallBuilder) Then(callback func(res interface{})) *CallBuilder {
	b.successCallback = callback
	return b
}

func (b *CallBuilder) Err(callback func(err Error)) *CallBuilder {
	b.errorCallback = callback
	return b
}

func (b *CallBuilder) Call() {
	engine.GetInstance().RunTask(engine.NewTaskBuilder().Run(func() (interface{}, error) {
		url := fmt.Sprintf("%s/api/rpc/%s", instance.url, b.fName)

		jsonStr, err := json.Marshal(b.call)
		if err != nil {
			return nil, newError(ErrorEncode, "failed to encode call")
		}

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
		if err != nil {
			return nil, newError(ErrorEncode, "failed to encode call")
		}

		req.Header.Set("Cadmean-RPC-Version", "2.1")
		req.Header.Set("Content-Type", "application/json")

		resp, err := instance.httpClient.Do(req)
		if err != nil {
			return nil, newError(ErrorCouldNotSendCall, "http error")
		}

		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return nil, newError(ErrorNotSuccessfulStatusCode, "not successful status code")
		}

		body, _ := ioutil.ReadAll(resp.Body)

		output := FunctionOutput{}
		err = json.Unmarshal(body, &output)
		if err != nil {
			return nil, newError(ErrorDecode, "failed to decode response")
		}

		return output, nil
	}).Than(func(res interface{}) {
		output := res.(FunctionOutput)
		if output.Error != 0 {
			b.errorCallback(newError(output.Error, fmt.Sprintf("function error")))
		} else {
			b.successCallback(output.Result)
		}
	}).Err(func(err error) {
		b.errorCallback(err.(Error))
	}).Build())
}

func newCallBuilder(fName string) *CallBuilder {
	return &CallBuilder{
		fName: fName,
		call: FunctionCall{
			Arguments: make([]interface{}, 0),
			Auth:      "",
		},
	}
}