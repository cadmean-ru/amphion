// +build js

package engine

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type webResourceReader struct {

}

func (m *webResourceReader) readResourceAsync(path string, callback ReadResourceCallback) {
	if !IsValidResourcePath(path) {
		callback(nil, errors.New("invalid resource path"))
		return
	}
	task := NewTaskBuilder().Run(func() (interface{}, error) {
		return http.Get(fmt.Sprintf("http://%s:%s/res%s", GetInstance().GetGlobalContext().GetHost(), GetInstance().GetGlobalContext().GetPort(), path))
	}).Than(func(res interface{}) {
		response := res.(*http.Response)
		defer response.Body.Close()
		body, err := ioutil.ReadAll(response.Body)
		callback(body, err)
	}).Err(func(err error) {
		fmt.Println(err)
		callback(nil, err)
	}).Build()

	GetInstance().RunTask(task)
}

func (m *webResourceReader) readResource(path string) ([]byte, error) {
	if !IsValidResourcePath(path) {
		return nil, errors.New("invalid resource path")
	}

	response, err := http.Get(fmt.Sprintf("http://%s/res%s", GetInstance().GetGlobalContext().GetHost(), path))
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	return body, err
}

func newResourceReader() resourceReader {
	return &webResourceReader{}
}