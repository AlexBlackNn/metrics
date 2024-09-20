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
curl -v --header "Content-Type: application/json" --request POST --data '{"id":"testGaugeMult","type":"counter"}' http://localhost:8080/value/
```
```bash
curl -v --header "Content-Type: application/json" --request POST --data '{"id":"testGaugeMult","type":"gauge"}' http://localhost:8080/value/
```

```bash
curl --header "Content-Type: application/json" --request POST --data '[{"id":"testGaugeMult","type":"gauge","value":465528.39165260154},{"id":"testGauge1Mult","type":"gauge","value":123.39165260154} ]' http://localhost:8080/updates/
```

```bash
curl -v --header "Accept-Encoding: gzip" --request GET  http://localhost:8080/ --compressed
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

за счет использования sync.Pool