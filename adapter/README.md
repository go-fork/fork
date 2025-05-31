# Fork HTTP Adapter Interface

This directory contains the adapter interface definition for the Fork HTTP framework.

## Overview

The adapter interface provides a standard contract for HTTP server implementations. All adapters must implement the `Adapter` interface defined in `adapter.go`.

## Available Adapters

The unified adapter (which previously resided here) has been moved to its own module:

- **Unified Adapter**: `github.com/go-fork/adapter/unified` - Supports HTTP/1.1, HTTP/2, and HTTP/3 (QUIC)

Other specialized adapters are available in the adapter directory:

- **FastHTTP Adapter**: `github.com/go-fork/adapter/fasthttp` - High-performance HTTP using FastHTTP
- **Standard HTTP Adapter**: `github.com/go-fork/adapter/http` - Standard Go net/http based adapter  
- **HTTP/2 Adapter**: `github.com/go-fork/adapter/http2` - Dedicated HTTP/2 support
- **QUIC Adapter**: `go.fork.vn/adapter/quic` - Dedicated HTTP/3 via QUIC

## Adapter Interface

```go
type Adapter interface {
    // Start the server with the given handler
    Serve(handler http.Handler) error
    
    // Stop the server gracefully
    Stop() error
    
    // Additional adapter-specific methods...
}
```

## Usage

To use an adapter with Fork:

```go
import (
    "go.fork.vn/fork"
    "github.com/go-fork/adapter/unified" // or any other adapter
)

func main() {
    adapter := unified.NewUnifiedAdapter()
    app := fork.New(fork.WithAdapter(adapter))
    
    app.Get("/", func(ctx fork.Context) {
        ctx.String(200, "Hello World!")
    })
    
    app.Start()
}
```

For detailed documentation on specific adapters, please refer to their respective module documentation.
