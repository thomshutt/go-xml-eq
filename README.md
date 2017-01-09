# go-xml-eq


### About

go-xml-eq is a Go library for checking the equality of two XML strings.

It aims to give human-readable output (rather than a diff) when mismatches happen.

### Usage

```go
package main

import (
	"fmt"

	"github.com/thomshutt/go-xml-eq/xml"
)

var xml1 = `
	<dogs>
		<dog name="rover"></dog>
		<dog name="lassie"></dog>
	</dogs>
`

var xml2 = `
	<dogs>
		<dog name="rover"></dog>
		<dog name="frank"></dog>
	</dogs>
`

func main() {
	err := xml.Equal(xml1, xml2)
	fmt.Println(err.Error())
}
```

```
$ go run main.go 
Attribute mismatch - "lassie" != "frank"
```

### TODO

* Return the path to where a mismatch is located
* Namespaces
