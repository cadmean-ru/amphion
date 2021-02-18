//+build windows linux darwin
//+build !ios
//+build !android

package atext

import (
	"fmt"
	"github.com/flopp/go-findfont"
	"github.com/golang/freetype/truetype"
	"io/ioutil"
)

func LoadFont(name string) (*Font, error) {
	fontPath, err := findfont.Find(fmt.Sprintf("%s.ttf", name))
	if err != nil {
		return nil, err
	}

	fontData, err := ioutil.ReadFile(fontPath)

	f, err := truetype.Parse(fontData)
	if err != nil {
		return nil, err
	}

	return &Font{
		name: name,
		ttf:  f,
	}, nil
}

func LoadDefaultFont() (*Font, error) {
	return LoadFont("Arial")
}