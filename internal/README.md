# internal

В данной директории и её поддиректориях будет содержаться имплементация вашего сервиса

go run ./cmd/server/main.go  --config="./cmd/server/config/demo.yaml"

```bash
curl -v -H "Content-Type: text/plain" -X POST  http://localhost:8076/update/gauge/param1/2
```

```bash
curl -v -H "Content-Type: text/plain" -X POST  http://localhost:8076/update/gauge1/param1/2
```

## Tests
```
cd cmd/server
go build -o server *.go
```

```bash
./metricstest -test.v -test.run=^TestIteration1$ -binary-path=/home/alex/Dev/GoYandex/metrics/cmd/server/server
```


```bash
curl -v -H "Content-Type: text/plain" -X POST http://localhost:8080/update/counter/testCounter/10
```


```bash
curl -v -H "Content-Type: text/plain" -X POST http://localhost:8080/update/gauge/testGauge/111
```

how to install golint
https://command-not-found.com/golint

```bash
golint ./...
```