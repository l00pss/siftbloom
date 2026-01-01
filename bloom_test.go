package siftbloom

import (
	"testing"
)

func TestNewSiftBloom(t *testing.T) {
	bloomResult := NewSiftBloom(1000, 5)
	if bloomResult.IsErr() {
		t.Fatal("Failed to create SiftBloom:", bloomResult.UnwrapErr())
	}

	bf := bloomResult.Unwrap()
	if bf == nil {
		t.Fatal("SiftBloom is nil")
	}

	if bf.hashFactor != 5 {
		t.Errorf("Expected hashFactor to be 5, got %d", bf.hashFactor)
	}
}

func TestAddAndContains(t *testing.T) {
	bloomResult := NewSiftBloom(1000, 3)
	if bloomResult.IsErr() {
		t.Fatal("Failed to create SiftBloom:", bloomResult.UnwrapErr())
	}

	bf := bloomResult.Unwrap()

	// Add elements
	bf.Add("hello")
	bf.Add("world")
	bf.Add(123)
	bf.Add(45.67)

	// Test contains
	if !bf.Contains("hello") {
		t.Error("Expected 'hello' to be in the bloom filter")
	}

	if !bf.Contains("world") {
		t.Error("Expected 'world' to be in the bloom filter")
	}

	if !bf.Contains(123) {
		t.Error("Expected 123 to be in the bloom filter")
	}

	if !bf.Contains(45.67) {
		t.Error("Expected 45.67 to be in the bloom filter")
	}

	// Test non-existent element
	if bf.Contains("nonexistent") {
		t.Log("'nonexistent' shows as present (false positive - this is expected)")
	}
}

func TestContainsNonExistentElements(t *testing.T) {
	bloomResult := NewSiftBloom(10000, 5)
	if bloomResult.IsErr() {
		t.Fatal("Failed to create SiftBloom:", bloomResult.UnwrapErr())
	}

	bf := bloomResult.Unwrap()

	// Add a few elements
	bf.Add("apple")
	bf.Add("banana")

	// Test many non-existent elements
	nonExistentElements := []any{
		"orange", "grape", "pear", "kiwi", "mango",
		999, 888, 777, 666, 555,
		"test123", "random456", "unknown789",
	}

	falsePositives := 0
	for _, elem := range nonExistentElements {
		if bf.Contains(elem) {
			falsePositives++
		}
	}

	// With good hash functions and reasonable size, false positives should be low
	falsePositiveRate := float64(falsePositives) / float64(len(nonExistentElements))
	t.Logf("False positive rate: %.2f%% (%d/%d)", falsePositiveRate*100, falsePositives, len(nonExistentElements))
}

func TestClear(t *testing.T) {
	bloomResult := NewSiftBloom(1000, 3)
	if bloomResult.IsErr() {
		t.Fatal("Failed to create SiftBloom:", bloomResult.UnwrapErr())
	}

	bf := bloomResult.Unwrap()

	// Add elements
	bf.Add("test1")
	bf.Add("test2")

	// Verify elements exist
	if !bf.Contains("test1") {
		t.Error("Expected 'test1' to be in the bloom filter before clear")
	}

	// Clear the filter
	bf.Clear()

	// Elements should not exist after clear
	if bf.Contains("test1") {
		t.Error("'test1' should not be in the bloom filter after clear")
	}

	if bf.Contains("test2") {
		t.Error("'test2' should not be in the bloom filter after clear")
	}
}

func TestConcurrency(t *testing.T) {
	bloomResult := NewSiftBloom(10000, 5)
	if bloomResult.IsErr() {
		t.Fatal("Failed to create SiftBloom:", bloomResult.UnwrapErr())
	}

	bf := bloomResult.Unwrap()

	// Test concurrent writes and reads
	done := make(chan bool, 2)

	// Writer goroutine
	go func() {
		for i := 0; i < 1000; i++ {
			bf.Add(i)
		}
		done <- true
	}()

	// Reader goroutine
	go func() {
		for i := 0; i < 1000; i++ {
			bf.Contains(i)
		}
		done <- true
	}()

	// Wait for both goroutines
	<-done
	<-done

	// Verify some elements exist
	if !bf.Contains(500) {
		t.Error("Expected 500 to be in the bloom filter")
	}
}

func TestDifferentTypes(t *testing.T) {
	bloomResult := NewSiftBloom(1000, 4)
	if bloomResult.IsErr() {
		t.Fatal("Failed to create SiftBloom:", bloomResult.UnwrapErr())
	}

	bf := bloomResult.Unwrap()

	// Test different data types
	testCases := []any{
		"string",
		123,
		45.67,
		true,
		[]byte("bytes"),
		struct{ Name string }{"test"},
	}

	// Add all test cases
	for _, tc := range testCases {
		bf.Add(tc)
	}

	// Verify all test cases exist
	for _, tc := range testCases {
		if !bf.Contains(tc) {
			t.Errorf("Expected %v (%T) to be in the bloom filter", tc, tc)
		}
	}
}

func BenchmarkAdd(b *testing.B) {
	bloomResult := NewSiftBloom(100000, 5)
	if bloomResult.IsErr() {
		b.Fatal("Failed to create SiftBloom:", bloomResult.UnwrapErr())
	}

	bf := bloomResult.Unwrap()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bf.Add(i)
	}
}

func BenchmarkContains(b *testing.B) {
	bloomResult := NewSiftBloom(100000, 5)
	if bloomResult.IsErr() {
		b.Fatal("Failed to create SiftBloom:", bloomResult.UnwrapErr())
	}

	bf := bloomResult.Unwrap()

	// Pre-populate with some data
	for i := 0; i < 10000; i++ {
		bf.Add(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bf.Contains(i % 20000) // Mix of existing and non-existing
	}
}
