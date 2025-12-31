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

import "github.com/l00pss/siftbloom"

func main() {
    // Create a new Bloom filter
    bf := siftbloom.New(1000, 0.01) // 1000 expected items, 1% false positive rate

    // Add items
    bf.Add([]byte("hello"))
    bf.Add([]byte("world"))

    // Check membership
    bf.Contains([]byte("hello")) // true
    bf.Contains([]byte("foo"))   // false (probably)
}
```

## License

MIT
