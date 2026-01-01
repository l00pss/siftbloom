package siftbloom

import (
	"errors"

	"github.com/l00pss/helpme/result"
)

var (
	IvalidArraySize = errors.New("Invalid array size")
)

type BitArray struct {
	data []byte
	size int
}

func NewBitArray(size int) result.Result[*BitArray] {
	if size <= 0 {
		return result.Err[*BitArray](IvalidArraySize)
	}
	return result.Ok(
		&BitArray{
			data: make([]byte, (size+7)/8),
			size: size,
		},
	)
}

func (b *BitArray) GetSize() int {
	return b.size
}

func (b *BitArray) Set(pos int, value bool) {
	byteIndex := pos / 8
	bitIndex := uint(pos % 8)

	if value {
		b.data[byteIndex] |= (1 << bitIndex)
	} else {
		b.data[byteIndex] &^= (1 << bitIndex)
	}
}

func (b *BitArray) Get(pos int) bool {
	byteIndex := pos / 8
	bitIndex := uint(pos % 8)
	return (b.data[byteIndex] & (1 << bitIndex)) != 0
}
