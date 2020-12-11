package rendering

import "github.com/cadmean-ru/amphion/common"

const (
	AttributeTransform    = 0
	AttributeFillColor    = 1
	AttributeStrokeColor  = 2
	AttributeStrokeWeight = 3
	AttributeText         = 4
	AttributeFontSize     = 5
	AttributeResIndex     = 6
	AttributeCornerRadius = 7
	AttributePoint        = 8
)

type Attribute struct {
	Code  byte
	Value []byte
}

func (a Attribute) GetLength() int {
	return len(a.Value) + 1
}

func (a Attribute) EncodeToByteArray() []byte {
	data := make([]byte, a.GetLength())
	data[0] = a.Code
	_ = common.CopyByteArray(a.Value, data, 1, len(a.Value))
	return data
}

func NewAttribute(code byte, value common.ByteArrayEncodable) Attribute {
	return Attribute{
		Code:  code,
		Value: value.EncodeToByteArray(),
	}
}
