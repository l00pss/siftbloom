# SiftBloom

A fast and memory-efficient Bloom filter implementation in Go.

<br>
<div align="center">
  <a href="https://www.buymeacoffee.com/l00pss" target="_blank"><img src="https://cdn.buymeacoffee.com/buttons/v2/default-yellow.png" alt="Buy Me A Coffee" style="height: 60px !important;width: 217px !important;" ></a>
</div>

## Installation

```bash
go get github.com/l00pss/siftbloom
```

## Usage

```go
package main

import (
    "fmt"
    "github.com/l00pss/siftbloom"
)

func main() {
    // Create a new Bloom filter
    bloomResult := siftbloom.NewSiftBloom(1000, 5) // 1000 size, 5 hash functions
    if bloomResult.IsErr() {
        panic(bloomResult.UnwrapErr())
    }
    
    bf := bloomResult.Unwrap()

    // Add items
    bf.Add("hello")
    bf.Add("world")
    bf.Add(123)

    // Check membership
    fmt.Println(bf.Contains("hello")) // true
    fmt.Println(bf.Contains("world")) // true
    fmt.Println(bf.Contains("foo"))   // false (probably)
    fmt.Println(bf.Contains(123))     // true
    
    // Clear filter
    bf.Clear()
}

```

## License

MIT
