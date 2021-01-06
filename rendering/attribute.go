package rendering

import (
	"github.com/cadmean-ru/amphion/common/a"
)

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

// Deprecated: no longer needed
type Attribute struct {
	Code  byte
	Value []byte
}

func (at Attribute) GetLength() int {
	return len(at.Value) + 1
}

func (at Attribute) EncodeToByteArray() []byte {
	data := make([]byte, at.GetLength())
	data[0] = at.Code
	_ = a.CopyByteArray(at.Value, data, 1, len(at.Value))
	return data
}

func NewAttribute(code byte, value a.ByteArrayEncodable) Attribute {
	return Attribute{
		Code:  code,
		Value: value.EncodeToByteArray(),
	}
}
