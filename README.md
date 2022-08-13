# go-deno, makes Deno be embedded easily

[Deno](https://github.com/denoland/deno) is a simple, modern and secure runtime for JavaScript and TypeScript that uses V8 and is built in Rust.

This package is intended to provide a wrapper to interact `Deno` with application written in golang.
With some helper functions, `go-deno` makes it simple to calle Deno from Golang, and `go-deno` can be
treated as an embeddable JavaScript/TypeScript.

### Install

The package is fully go-getable, So, just type

  `go get github.com/rosbit/go-deno`

to install.

### Usage

Suppose there's a Javascript file named `a.js` like this:

```javascript
function add(a, b) {
    return a+b
}
```

one can call the Javascript function `add()` in Go code like the following:

```go
package main

import (
  "github.com/rosbit/go-deno"
  "fmt"
)

var add func(int, int)int

func main() {
  ctx, err := deno.NewDeno("/path/to/deno-exe/deno", "a.js")
  if err != nil {
     fmt.Printf("%v\n", err)
     return
  }
  defer ctx.Quit()

  // method 1: bind JS function with a golang var
  if err := ctx.BindFunc("add", &add); err != nil {
     fmt.Printf("%v\n", err)
     return
  }
  res := add(1, 2)

  // method 2: call JS function using Call
  res, err := ctx.CallFunc("add", 1, 2)
  if err != nil {
     fmt.Printf("%v\n", err)
     return
  }

  fmt.Println("result is:", res)
}
```

### Status

The package is not fully tested, so be careful.

### Contribution

Pull requests are welcome! Also, if you want to discuss something send a pull request with proposal and changes.

__Convention:__ fork the repository and make changes on your fork in a feature branch.
