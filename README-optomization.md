```bash
cd ../internal/middleware
```

Старая версия gzipCompressor

```go
package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"testing"
)

// DummyResponseWriter implements http.ResponseWriter but discards the output
type DummyResponseWriter struct {
	header http.Header
	code   int
	wrote  bool
}

func (w *DummyResponseWriter) Header() http.Header {
	if w.header == nil {
		w.header = make(http.Header)
	}
	return w.header
}

func (w *DummyResponseWriter) Write(b []byte) (int, error) {
	if !w.wrote {
		w.wrote = true
	}
	return len(b), nil // Discard the data
}

func (w *DummyResponseWriter) WriteHeader(code int) {
	if !w.wrote {
		w.wrote = true
	}
	w.code = code
}

func BenchmarkGzipCompressor(b *testing.B) {
	// Create a dummy handler to simulate a real handler
	dummyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate some data to be written
		_, err := io.WriteString(w, "This is some test data to be compressed.")
		if err != nil {
			b.Fatal(err)
		}
	})

	// Create a request for the benchmark
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		b.Fatal(err)
	}
	req.Header.Set("Accept-Encoding", "gzip")

	// Run the benchmark
	gzipCompressor := GzipCompressor(nil, gzip.DefaultCompression)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Create a DummyResponseWriter to discard output
		w := &DummyResponseWriter{}

		// Set the header before calling ServeHTTP
		w.Header().Set("Content-Encoding", "gzip")
		w.WriteHeader(200)

		gzipCompressor(dummyHandler).ServeHTTP(w, req)
	}
}
```

Запускаем тесты бенчмарка

```bash
go test -bench . -benchmem -benchtime=1s -memprofile mem.out
```
```bash
go tool pprof mem.out
```

```bash
go test -bench . -benchmem -benchtime=1s -memprofile mem.out
goos: linux
goarch: amd64
pkg: github.com/AlexBlackNn/metrics/internal/middleware/old_gzip
cpu: Intel(R) Core(TM) i7-5960X CPU @ 3.00GHz
BenchmarkGzipCompressor-16       1000000              1757 ns/op             712 B/op          7 allocs/op
PASS
ok      github.com/AlexBlackNn/metrics/internal/middleware/old_gzip     1.778s
```

```bash
go tool pprof mem.out
File: old_gzip.test
Type: alloc_space
Time: Sep 19, 2024 at 10:54pm (MSK)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top
Showing nodes accounting for 661.64MB, 99.92% of 662.14MB total
Dropped 4 nodes (cum <= 3.31MB)
Showing top 10 nodes out of 11
      flat  flat%   sum%        cum   cum%
  336.11MB 50.76% 50.76%   336.11MB 50.76%  net/textproto.MIMEHeader.Set (inline)
  163.53MB 24.70% 75.46%   163.53MB 24.70%  compress/gzip.NewWriterLevel
   53.50MB  8.08% 83.54%   254.03MB 38.37%  github.com/AlexBlackNn/metrics/internal/middleware/old_gzip.BenchmarkGzipCompressor.BenchmarkGzipCompressor.GzipCompressor.func2.func3
      50MB  7.55% 91.09%       50MB  7.55%  github.com/AlexBlackNn/metrics/internal/middleware/old_gzip.(*DummyResponseWriter).Header (inline)
      37MB  5.59% 96.68%       37MB  5.59%  io.WriteString
   21.50MB  3.25% 99.92%   661.64MB 99.92%  github.com/AlexBlackNn/metrics/internal/middleware/old_gzip.BenchmarkGzipCompressor
         0     0% 99.92%       37MB  5.59%  github.com/AlexBlackNn/metrics/internal/middleware/old_gzip.BenchmarkGzipCompressor.func1
         0     0% 99.92%   254.03MB 38.37%  net/http.HandlerFunc.ServeHTTP (partial-inline)
         0     0% 99.92%   336.11MB 50.76%  net/http.Header.Set (inline)
         0     0% 99.92%   661.64MB 99.92%  testing.(*B).launch
(pprof) 

```


сделали оптимизацию на базе Sync Pool
```bash
goos: linux
goarch: amd64
pkg: github.com/AlexBlackNn/metrics/internal/middleware
cpu: Intel(R) Core(TM) i7-5960X CPU @ 3.00GHz
BenchmarkGzipCompressor-16       1000000              1439 ns/op             536 B/op          6 allocs/op
PASS
ok      github.com/AlexBlackNn/metrics/internal/middleware      1.458s

```

```bash
 go tool pprof mem.out
File: middleware.test
Type: alloc_space
Time: Sep 19, 2024 at 10:53pm (MSK)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top
Showing nodes accounting for 471.11MB, 99.89% of 471.61MB total
Dropped 16 nodes (cum <= 2.36MB)
Showing top 10 nodes out of 11
      flat  flat%   sum%        cum   cum%
  317.60MB 67.34% 67.34%   317.60MB 67.34%  net/textproto.MIMEHeader.Set (inline)
      47MB  9.97% 77.31%       47MB  9.97%  github.com/AlexBlackNn/metrics/internal/middleware.(*DummyResponseWriter).Header (inline)
   46.50MB  9.86% 87.17%    46.50MB  9.86%  github.com/AlexBlackNn/metrics/internal/middleware.(*gzipWriterPool).Get
   35.50MB  7.53% 94.70%    35.50MB  7.53%  io.WriteString
   24.50MB  5.20% 99.89%   471.11MB 99.89%  github.com/AlexBlackNn/metrics/internal/middleware.BenchmarkGzipCompressor
         0     0% 99.89%       82MB 17.39%  github.com/AlexBlackNn/metrics/internal/middleware.BenchmarkGzipCompressor.BenchmarkGzipCompressor.GzipCompressor.func2.func3
         0     0% 99.89%    35.50MB  7.53%  github.com/AlexBlackNn/metrics/internal/middleware.BenchmarkGzipCompressor.func1
         0     0% 99.89%       82MB 17.39%  net/http.HandlerFunc.ServeHTTP (partial-inline)
         0     0% 99.89%   317.60MB 67.34%  net/http.Header.Set (inline)
         0     0% 99.89%   471.11MB 99.89%  testing.(*B).launch
(pprof)    

```

Оптимизировано:   163.53MB 24.70% 75.46%   163.53MB 24.70%  compress/gzip.NewWriterLevel за счет использования sync.Pool