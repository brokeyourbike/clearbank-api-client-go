# clearbank-api-client-go

[![Go Reference](https://pkg.go.dev/badge/github.com/brokeyourbike/clearbank-api-client-go.svg)](https://pkg.go.dev/github.com/brokeyourbike/clearbank-api-client-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/brokeyourbike/clearbank-api-client-go)](https://goreportcard.com/report/github.com/brokeyourbike/clearbank-api-client-go)
[![codecov](https://codecov.io/gh/brokeyourbike/clearbank-api-client-go/graph/badge.svg?token=mbn319zqNO)](https://codecov.io/gh/brokeyourbike/clearbank-api-client-go)

[ClearBank](https://clearbank.github.io/) API Client for Go

## Installation

```bash
go get github.com/brokeyourbike/clearbank-api-client-go
```

## Usage

```go
client := clearbank.NewClient("token", signer)

err := client.Test(clearbank.RequestIdContext(context.TODO(), "123"), "hello")
require.NoError(t, err)
```

## Authors
- [Ivan Stasiuk](https://github.com/brokeyourbike) | [Twitter](https://twitter.com/brokeyourbike) | [LinkedIn](https://www.linkedin.com/in/brokeyourbike) | [stasi.uk](https://stasi.uk)

## License
[BSD-3-Clause License](https://github.com/brokeyourbike/clearbank-api-client-go/blob/main/LICENSE)
