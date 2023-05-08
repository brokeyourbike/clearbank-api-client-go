# clearbank-api-client-go

[![Go Reference](https://pkg.go.dev/badge/github.com/brokeyourbike/clearbank-api-client-go.svg)](https://pkg.go.dev/github.com/brokeyourbike/clearbank-api-client-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/brokeyourbike/clearbank-api-client-go)](https://goreportcard.com/report/github.com/brokeyourbike/clearbank-api-client-go)
[![Maintainability](https://api.codeclimate.com/v1/badges/147e88944ef3303bba6d/maintainability)](https://codeclimate.com/github/brokeyourbike/clearbank-api-client-go/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/147e88944ef3303bba6d/test_coverage)](https://codeclimate.com/github/brokeyourbike/clearbank-api-client-go/test_coverage)

ClearBank API Client for Go

## Example

```go
package main

import (
    "context"

    "github.com/brokeyourbike/clearbank-api-client-go"
    "github.com/stretchr/testify/assert"
)

func main() {
    client := clearbank.NewClient("token", signer, clearbank.WithBaseURL("https://api.clear.bank"))

    ctx := context.Background()

    err := client.Test(clearbank.RequestIdContext(ctx, "123"), "hello")
    assert.NoError(t, err)
}
```