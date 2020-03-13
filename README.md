[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/dimiro1/reply)
[![Build Status](https://travis-ci.org/dimiro1/reply.svg?branch=master)](https://travis-ci.org/dimiro1/reply)
[![Go Report Card](https://goreportcard.com/badge/github.com/dimiro1/reply)](https://goreportcard.com/report/github.com/dimiro1/reply)

Try browsing [the code on Sourcegraph](https://sourcegraph.com/github.com/dimiro1/reply)!

# reply

 Library to trim replies from plain text email. (Golang port of https://github.com/discourse/email_reply_trimmer)

# Usage

```go
package main

import (
    "fmt"

    "github.com/dimiro1/reply"
)

func main() {
    message := `
        This is before the embedded email.
        
        On Wed, Sep 25, 2013, at 03:57 PM, richard_clark wrote:
        
        Richard> This is the embedded email
        
        This is after the embedded email and will not show up because 99% of the times
        this is the signature...
    `
	fmt.Println(reply.FromText(message))
}
```

will output:

```text
This is before the embedded email.
```