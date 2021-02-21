//+build windows linux darwin
//+build !ios
//+build !android

package atext

import (
	"fmt"
	"github.com/flopp/go-findfont"
	"io/ioutil"
)

//LoadFont tries to find font file with the specified name in well-known locations and parse it.
//Supported only on PC.
func LoadFont(name string) (*Font, error) {
	fontPath, err := findfont.Find(fmt.Sprintf("%s.ttf", name))
	if err != nil {
		return nil, err
	}

	fontData, err := ioutil.ReadFile(fontPath)
	if err != nil {
		return nil, err
	}

	return ParseFont(fontData)
}

func LoadDefaultFont() (*Font, error) {
	return LoadFont("Arial")
}