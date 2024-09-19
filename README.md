### Сервис сбора метрик

1. Профайл полученный с помощью  

```
go tool pprof -http=":9090" -seconds=30 http://localhost:8080/debug/pprof/heap
``` 

[profile_handlers.pb.gz](profiles%2Fprofile_handlers.pb.gz)



![Screenshot from 2024-09-19 16-47-26.png](..%2F..%2FPictures%2FScreenshots%2FScreenshot%20from%202024-09-19%2016-47-26.png)

Cогласно, отчету, наибольшее потребление памяти приходится на GzipCompressor и GzipDecompressor





** Literature **
RESTY USEFUL
https://www.alldevstack.com/ru/go-resty-tutorial/go-resty-quickstart.html

how to install golint
https://command-not-found.com/golint

easy-json commands
```
sudo apt install golang-easyjson
easyjson -all /home/alex/GolandProjects/metrics/internal/handlers/v3/metrics_handlers_response.go 
```

```
For manual local tests
```

```bash
curl -v -H "Content-Type: text/plain" -X POST  http://localhost:8080/update/gauge/param1/2
```

```bash
curl -v -H "Content-Type: text/plain" -X POST  http://localhost:8080/update/gauge1/param1/2
```


```bash
curl -v -H "Content-Type: text/plain" -X POST http://localhost:8080/update/counter/testCounter1/10
```


```bash
curl -v -H "Content-Type: text/plain" -X POST http://localhost:8080/update/gauge/testGauge/111
```

```bash
curl -v -H "Content-Type: text/plain" -X POST http://localhost:8080/update/gauge/Lookups/20.4
```

```bash
curl http://localhost:8080/value/gauge/Lookups
```

---------V2-----------
```bash
curl -v -H "Content-Type: text/plain" -X POST http://localhost:8080/update/gauge/Lookups/21.4
```
```bash
curl --header "Content-Type: application/json" --request POST --data '{"id":"Lookups","type":"gauge"}' http://localhost:8080/value/
```


```bash
curl -v -H "Content-Type: text/plain" -X POST http://localhost:8080/update/counter/testCounter1/10
```
```bash
curl --header "Content-Type: application/json" --request POST --data '{"id":"testCounter1","type":"counter"}' http://localhost:8080/value/
```



```bash
curl --header "Content-Type: application/json" --request POST --data '{"id":"testCounter1","type":"counter","delta":10}' http://localhost:8080/update/
```
```bash
curl --header "Content-Type: application/json" --request POST --data '{"id":"testCounter1","type":"counter"}' http://localhost:8080/value/
```


```bash
curl --header "Content-Type: application/json" --request POST --data '{"id":"testGauge","type":"gauge","value":465528.39165260154}' http://localhost:8080/update/
```
```bash
curl --header "Content-Type: application/json" --request POST --data '{"id":"testGauge","type":"gauge"}' http://localhost:8080/value/
```


----------------
```bash
curl -v -H "Content-Type: text/plain" -X POST  http://localhost:8080/update/gauge/param2/2
```
```bash
curl -v -H "Content-Type: text/plain" -X GET  http://localhost:8080/value/gauge/param2
```
```bash
curl --header "Content-Type: application/json" --request POST --data '{"id":"testCounter14","type":"counter","delta":10}' http://localhost:8080/update/
```
```bash
curl -v --header "Content-Type: application/json" --request POST --data '{"id":"testCounter14","type":"counter"}' http://localhost:8080/value/
```
```bash
curl -v --header "Accept-Encoding: gzip" --header "Content-Type: application/json" --request POST --data '{"id":"test_counter","type":"counter"}' http://localhost:8080/value/ --compressed
```

```bash
curl -v --request GET  http://localhost:8080/
```

```bash
curl -v --header "Accept-Encoding: gzip" --request GET  http://localhost:8080/ --compressed
```

```bash
curl -v -X GET  http://localhost:8080/ping
```

# Моки
Создаем руками папку для хранения мока
```bash
mkdir pkg/storage/mock
```
Запускаем создание мока 
```bash
mockgen -destination=pkg/storage/mockstorage/mock_storage.go -package=mockstorage github.com/AlexBlackNn/metrics/internal/services/metricsservice MetricsStorage,HealthChecker
```

# Миграции 
go run ./cmd/migrator/postgres  --migrations-path=./migrations

# SQL 
CREATE EXTENSION pg_stat_statements; 


```bash
curl -v --header "Content-Type: application/json" --request POST --data '{"id":"test_counter","type":"counter"}' http://localhost:8080/value/
```
```bash
curl -v --header "Content-Type: application/json" --request POST --data '{"id":"test_gauge","type":"gauge"}' http://localhost:8080/value/
```

```bash
curl --header "Content-Type: application/json" --request POST --data '[{"id":"testGaugeMult","type":"gauge","value":465528.39165260154},{"id":"testGauge1Mult","type":"gauge","value":123.39165260154} ]' http://localhost:8080/updates/
```




linters
https://golangci-lint.run/welcome/install/#binaries
golangci-lint run -v

staticcheck ./...


```bash
cd /home/alex/Dev/GolandYandex/metrics/app/agent/hash
go test -bench .
go test -bench . -benchmem 
```

```bash
goos: linux
goarch: amd64
pkg: github.com/AlexBlackNn/metrics/app/agent/hash
cpu: AMD Ryzen 5 5500U with Radeon Graphics         
BenchmarkMetricHash-12           2930110               401.1 ns/op           176 B/op          4 allocs/op
PASS
ok      github.com/AlexBlackNn/metrics/app/agent/hash   1.595s
```


## Подключение профилировщика к сервису

1. подключили профилировщик в routers через middleware
2. запускаем сервис 
3.  go tool pprof -http=":9090" -seconds=30 http://localhost:8080/debug/pprof/profile 

4. Делаем запросы 
```bash
curl --header "Content-Type: application/json" --request POST --data '{"id":"testCounter14","type":"counter","delta":10}' http://localhost:8080/update/
```
```bash
curl -v --header "Content-Type: application/json" --request POST --data '{"id":"testCounter14","type":"counter"}' http://localhost:8080/value/
```

5. Переходим по ссылке http://localhost:9090 и видим
   
       5.1 в меньшем масштабе
       ![Screenshot from 2024-09-17 21-16-30.png](cmd%2Fdocs%2FScreenshot%20from%202024-09-17%2021-16-30.png
        
        5.2 в большем масштабе
        ![Screenshot from 2024-09-17 21-16-49.png](cmd%2Fdocs%2FScreenshot%20from%202024-09-17%2021-16-49.png)

6. 
```bash
cd /home/alex/GolandProjects/metrics/app/agent/hash && \
go test -bench . -benchmem 
```

```bash
go test -bench . -benchmem -benchtime=1s -memprofile mem.out
```
```bash
go tool pprof mem.out
```

На базе Sync Pool

```bash

type gzipWriterPool struct {
	p sync.Pool
}

func (gp *gzipWriterPool) Get(w http.ResponseWriter, compressorLevel int) (*GzipWriter, error) {

	gzipWriter := gp.p.Get()
	if gzipWriter == nil {
		gzipWr, err := gzip.NewWriterLevel(w, compressorLevel)
		if err != nil {
			return nil, err
		}
		return &GzipWriter{ResWriter: w, Writer: gzipWr}, nil
	}
	return &GzipWriter{ResWriter: w, Writer: gzipWriter.(*gzip.Writer)}, nil
}

func (gp *gzipWriterPool) Put(gzipWriter *GzipWriter) error {
	// Reset the writer to its initial state
	err := gzipWriter.Writer.Flush()
	if err != nil {
		return err
	}
	gzipWriter.Writer.Reset(io.Discard)
	// Put the writer back into the pool
	gp.p.Put(gzipWriter.Writer)
	return nil
}

func (gp *gzipWriterPool) PutNoFlush(gzipWriter *GzipWriter) error {
	// Reset the writer to its initial state
	gzipWriter.Writer.Reset(io.Discard)
	// Put the writer back into the pool
	gp.p.Put(gzipWriter.Writer)
	return nil
}

var gzipWrPool = &gzipWriterPool{}

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

До оптимизации 

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

Оптимизировано:   163.53MB 24.70% 75.46%   163.53MB 24.70%  compress/gzip.NewWriterLevel за счет использования sync.Pool