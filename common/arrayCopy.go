package common

import (
	"errors"
)

func CopyByteArray(src []byte, target []byte, start, length int) error {
	if start > len(target) {
		return errors.New("start index outside of the target array")
	}

	if length > len(src) {
		return errors.New("length is greater than the length of source array")
	}

	if start + length > len(target) {
		return errors.New("target array is to small")
	}

	j := start
	for i := 0; i < length; i++ {
		target[j] = src[i]
		j++
	}

	return nil
}
