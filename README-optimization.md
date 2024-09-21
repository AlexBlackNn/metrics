# Оптимизация проекта по используемой памяти

## Находим проблемные места в проекте по потребляемой памяти

1. Добавили middleware для сбора метрик для pprof https://github.com/go-chi/chi/blob/master/middleware/profiler.go 
2. Запускаем сервис   

```bash
   go run cmd/server/main.go --d postgres://app:app123@localhost:5432/metric_db?sslmode=disable
```

3. Запускаем запись метрик 

```
go tool pprof -http=":9090" -seconds=30 http://localhost:8080/debug/pprof/heap
``` 
4. Делаем запроса на ручки (лучше запускать Нагрузочное тестирование с типовой нагрузкой, но так как ее нет, руками)

```bash
curl --header "Content-Type: application/json" --request POST --data '{"id":"test_counter","type":"counter","delta":10}' http://localhost:8080/update/
```
```bash
curl -v --header "Content-Type: application/json" --request POST --data '{"id":"test_counter","type":"counter"}' http://localhost:8080/value/
```

```bash
curl --header "Content-Type: application/json" --request POST --data '[{"id":"testGa2ugeMult","type":"gauge","value":1.39165260154},{"id":"testGauge1Mult","type":"gauge","value":123.39165260154} ]' http://localhost:8080/updates/
```

```bash
curl -v --header "Accept-Encoding: gzip" --header "Content-Type: application/json" --request POST --data '{"id":"test_counter","type":"counter"}' http://localhost:8080/value/ 
```

4. Сохранем профиль 
[profile.pb.gz](profiler%2Fprofile.pb.gz)

5. смотрим профиль 
```bash
go tool pprof ./profiler/profile.pb.gz
```

6. 10 наиболее затратных по памяти операций 

```bash
(pprof) top
Showing nodes accounting for 2870.33kB, 72.49% of 3959.67kB total
Showing top 10 nodes out of 37
      flat  flat%   sum%        cum   cum%
 1805.17kB 45.59% 45.59%  2902.86kB 73.31%  compress/flate.NewWriter (inline)
 1097.69kB 27.72% 73.31%  1097.69kB 27.72%  compress/flate.(*compressor).initDeflate (inline)
 -544.67kB 13.76% 59.56%  -544.67kB 13.76%  net.open
  512.14kB 12.93% 72.49%   512.14kB 12.93%  github.com/go-playground/validator/v10.New.func1
         0     0% 72.49%  1097.69kB 27.72%  compress/flate.(*compressor).init
         0     0% 72.49%  2902.86kB 73.31%  compress/gzip.(*Writer).Write
         0     0% 72.49%     3415kB 86.24%  github.com/AlexBlackNn/metrics/cmd/server/router.NewChiRouter.GzipCompressor.func6.1
         0     0% 72.49%     3415kB 86.24%  github.com/AlexBlackNn/metrics/cmd/server/router.NewChiRouter.GzipDecompressor.func5.1
         0     0% 72.49%     3415kB 86.24%  github.com/AlexBlackNn/metrics/cmd/server/router.NewChiRouter.HashChecker.func4.1
         0     0% 72.49%     3415kB 86.24%  github.com/AlexBlackNn/metrics/cmd/server/router.NewChiRouter.Logger.func3.1
(pprof) 
```

7. по кумулятивному использованию памяти:
```bash
   (pprof) top10 -cum
   Showing nodes accounting for 0, 0% of 3.87MB total
   Showing top 10 nodes out of 37
   flat  flat%   sum%        cum   cum%
   0     0%     0%     3.33MB 86.24%  github.com/AlexBlackNn/metrics/cmd/server/router.NewChiRouter.GzipCompressor.func6.1
   0     0%     0%     3.33MB 86.24%  github.com/AlexBlackNn/metrics/cmd/server/router.NewChiRouter.GzipDecompressor.func5.1
   0     0%     0%     3.33MB 86.24%  github.com/AlexBlackNn/metrics/cmd/server/router.NewChiRouter.HashChecker.func4.1
   0     0%     0%     3.33MB 86.24%  github.com/AlexBlackNn/metrics/cmd/server/router.NewChiRouter.Logger.func3.1
   0     0%     0%     3.33MB 86.24%  github.com/AlexBlackNn/metrics/internal/handlers/v2.(*MetricHandlers).GetOneMetric
   0     0%     0%     3.33MB 86.24%  github.com/go-chi/chi/v5.(*Mux).Mount.func1
   0     0%     0%     3.33MB 86.24%  github.com/go-chi/chi/v5.(*Mux).ServeHTTP
   0     0%     0%     3.33MB 86.24%  github.com/go-chi/chi/v5.(*Mux).routeHTTP
   0     0%     0%     3.33MB 86.24%  github.com/go-chi/chi/v5/middleware.Recoverer.func1
   0     0%     0%     3.33MB 86.24%  github.com/go-chi/chi/v5/middleware.RequestID.func1
```

## Исправляем проблемные места в проекте. 

###  Оптимизация gzipCompressor middleware по используемой памяти.
Старая версия gzipCompressor

```go
package middleware

import (
	"compress/gzip"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"sync"
)

type GzipWriter struct {
	ResWriter       http.ResponseWriter
	Writer          *gzip.Writer
	GzipWriterMutex sync.Mutex
	GzipFlag        bool
}

func (w *GzipWriter) Header() http.Header {
	return w.ResWriter.Header()
}

func (w *GzipWriter) WriteHeader(statusCode int) {
	if !strings.Contains(w.ResWriter.Header().Get("Content-Type"), "application/json") &&
		!strings.Contains(w.ResWriter.Header().Get("Content-Type"), "text/html") {
		w.ResWriter.WriteHeader(statusCode)
		w.GzipFlag = false
		return
	}
	w.GzipFlag = true
	w.ResWriter.Header().Set("Content-Encoding", "gzip")
	w.ResWriter.WriteHeader(statusCode)
}

func (w *GzipWriter) Write(b []byte) (int, error) {
	if !strings.Contains(w.ResWriter.Header().Get("Content-Type"), "application/json") &&
		!strings.Contains(w.ResWriter.Header().Get("Content-Type"), "text/html") {
		return w.ResWriter.Write(b)
	}
	return w.Writer.Write(b)
}

func (w *GzipWriter) Close() error {
	return w.Writer.Close()
}

func GzipCompressor(log *slog.Logger, compressorLevel int) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		log := log.With(
			slog.String("component", "middleware/gzip"),
		)
		log.Info("gzip compressor enabled")
		fn := func(w http.ResponseWriter, r *http.Request) {

			if !strings.Contains(strings.Join(r.Header.Values("Accept-Encoding"), " "), "gzip") {
				// If gzip is not supported then return uncompressed page.
				next.ServeHTTP(w, r)
				return
			}

			log.Info("gzip is supported")

			gzipWr, err := gzip.NewWriterLevel(w, compressorLevel)
			if err != nil {
				log.Error("failed to compress gzip")
				_, err := io.WriteString(w, err.Error())
				if err != nil {
					log.Error("failed to inform user")
					return
				}
				return
			}

			gz := &GzipWriter{ResWriter: w, Writer: gzipWr}
			next.ServeHTTP(gz, r)
			if gz.GzipFlag {
				err := gzipWr.Close()
				if err != nil {
					log.Error("failed to close gzip")
					_, err := io.WriteString(w, err.Error())
					if err != nil {
						log.Error("failed to inform user")
						return
					}
					return
				}
			}
		}
		return http.HandlerFunc(fn)
	}
}
```

Запускаем тесты бенчмарка


```bash
cd ./internal/middleware
go test -bench . -benchmem -benchtime=1s -memprofile mem.out
```

```
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
```
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

Видим, что gzip.NewWriterLevel часто создается
```
163.53MB 24.70% 75.46%   163.53MB 24.70%  compress/gzip.NewWriterLevel
```

сделали оптимизацию на базе Sync Pool https://victoriametrics.com/blog/tsdb-performance-techniques-sync-pool/ 

```bash
goos: linux
goarch: amd64
pkg: github.com/AlexBlackNn/metrics/internal/middleware
cpu: Intel(R) Core(TM) i7-5960X CPU @ 3.00GHz
BenchmarkGzipCompressor-16       1000000              1439 ns/op             536 B/op          6 allocs/op
PASS
ok      github.com/AlexBlackNn/metrics/internal/middleware      1.458s

```

1 аллокация памяти ушла 

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

по кумулятивному использованию памяти:
```bash
(pprof) top10 -cum
Showing nodes accounting for 1468.87MB, 91.38% of 1607.38MB total
Dropped 19 nodes (cum <= 8.04MB)
Showing top 10 nodes out of 11
      flat  flat%   sum%        cum   cum%
      64MB  3.98%  3.98%  1606.88MB   100%  github.com/AlexBlackNn/metrics/internal/middleware.BenchmarkGzipCompressor
         0     0%  3.98%  1606.88MB   100%  testing.(*B).launch
         0     0%  3.98%  1606.88MB   100%  testing.(*B).runN
         0     0%  3.98%  1108.36MB 68.95%  net/http.Header.Set (inline)
 1108.36MB 68.95% 72.94%  1108.36MB 68.95%  net/textproto.MIMEHeader.Set (inline)
         0     0% 72.94%   280.51MB 17.45%  github.com/AlexBlackNn/metrics/internal/middleware.BenchmarkGzipCompressor.BenchmarkGzipCompressor.GzipCompressor.func2.func3
         0     0% 72.94%   280.51MB 17.45%  net/http.HandlerFunc.ServeHTTP (partial-inline)
  154.01MB  9.58% 82.52%   154.01MB  9.58%  github.com/AlexBlackNn/metrics/internal/middleware.(*DummyResponseWriter).Header (inline)
  142.51MB  8.87% 91.38%   142.51MB  8.87%  github.com/AlexBlackNn/metrics/internal/middleware.(*gzipWriterPool).Get
         0     0% 91.38%   137.51MB  8.55%  github.com/AlexBlackNn/metrics/internal/middleware.BenchmarkGzipCompressor.func1
(pprof) 
```

в топе находятся бенчмарк Go. 

Смотрим дельту: 
```bash
 go tool pprof -top -diff_base=profiles/base_gzip_compressor_mem.out profiles/result_gzip_compressor_mem.out
```

```
 go tool pprof -top -diff_base=profiles/base_gzip_compressor_mem.out profiles/result_gzip_compressor_mem.out
File: middleware.test
Type: alloc_space
Time: Sep 20, 2024 at 11:38am (MSK)
Showing nodes accounting for -372.54MB, 18.86% of 1974.94MB total
Dropped 17 nodes (cum <= 9.87MB)
      flat  flat%   sum%        cum   cum%
 -487.58MB 24.69% 24.69%  -487.58MB 24.69%  compress/gzip.NewWriterLevel
  143.01MB  7.24% 17.45%   143.51MB  7.27%  github.com/AlexBlackNn/metrics/internal/middleware.(*gzipWriterPool).Get
 -137.01MB  6.94% 24.38%  -482.58MB 24.44%  github.com/AlexBlackNn/metrics/internal/middleware.BenchmarkGzipCompressor.BenchmarkGzipCompressor.GzipCompressor.func2.func3
  119.04MB  6.03% 18.36%   119.04MB  6.03%  net/textproto.MIMEHeader.Set (inline)
     -18MB  0.91% 19.27%      -18MB  0.91%  github.com/AlexBlackNn/metrics/internal/middleware.(*DummyResponseWriter).Header (inline)
       8MB  0.41% 18.86%  -373.54MB 18.91%  github.com/AlexBlackNn/metrics/internal/middleware.BenchmarkGzipCompressor
         0     0% 18.86%  -482.58MB 24.44%  net/http.HandlerFunc.ServeHTTP (partial-inline)
         0     0% 18.86%   119.04MB  6.03%  net/http.Header.Set (inline)
         0     0% 18.86%  -373.54MB 18.91%  testing.(*B).launch
         0     0% 18.86%  -373.54MB 18.91%  testing.(*B).runN

```


Добавили 143.01MB за счет приведения типов, но сократили на 487.58MB выделение памяти на gzip.NewWriterLevel за счет использования sync Pool
```
-487.58MB 24.69% 24.69%  -487.58MB 24.69%  compress/gzip.NewWriterLevel
143.01MB  7.24% 17.45%   143.51MB  7.27%  github.com/AlexBlackNn/metrics/internal/middleware.(*gzipWriterPool).Get 
```