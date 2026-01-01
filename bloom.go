package siftbloom

import (
	"fmt"
	"hash/fnv"
	"reflect"
	"sync"

	"github.com/l00pss/helpme/result"
)

type SiftBloom struct {
	mu         sync.RWMutex
	bits       BitArray
	hashFactor int
}

func NewSiftBloom(size int, hashFactor int) result.Result[*SiftBloom] {
	bitsResult := NewBitArray(size)
	if bitsResult.IsErr() {
		return result.Err[*SiftBloom](bitsResult.UnwrapErr())
	}

	return result.Ok(&SiftBloom{
		bits:       *bitsResult.Unwrap(),
		mu:         sync.RWMutex{},
		hashFactor: hashFactor,
	})
}

func (s *SiftBloom) Add(element any) {
	s.mu.Lock()
	defer s.mu.Unlock()

	bytes := ToBytes(element)
	hashes := s.getMultipleHashes(bytes)

	for _, hashIndex := range hashes {
		s.bits.Set(hashIndex, true)
	}
}

func (s *SiftBloom) getMultipleHashes(data []byte) []int {
	h1 := s.hashFNV1(data)
	h2 := s.hashFNV2(data)

	hashes := make([]int, s.hashFactor)
	for i := 0; i < s.hashFactor; i++ {
		hashValue := (h1 + uint64(i)*h2) % uint64(s.bits.GetSize())
		hashes[i] = int(hashValue)
	}
	return hashes
}

func (s *SiftBloom) hashFNV1(data []byte) uint64 {
	h := fnv.New64a()
	h.Write(data)
	return h.Sum64()
}

func (s *SiftBloom) hashFNV2(data []byte) uint64 {
	h := fnv.New64()
	h.Write(data)
	return h.Sum64()
}

func ToBytes(element any) []byte {
	val := reflect.ValueOf(element)
	switch val.Kind() {
	case reflect.Struct:
		return fmt.Appendf(nil, "%+v", element)
	case reflect.Slice:
		return fmt.Appendf(nil, "%v", element)
	default:
		return fmt.Appendf(nil, "%v", element)
	}
}

func (s *SiftBloom) Contains(element any) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	bytes := ToBytes(element)
	hashes := s.getMultipleHashes(bytes)

	for _, hashIndex := range hashes {
		if !s.bits.Get(hashIndex) {
			return false
		}
	}
	return true
}

func (s *SiftBloom) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()

	res := NewBitArray(s.bits.GetSize())
	if res.IsOk() {
		s.bits = *res.Unwrap()
	}
}
