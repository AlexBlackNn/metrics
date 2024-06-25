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
./metricstest -test.v -test.run=^TestIteration1$ -binary-path=./cmd/server/server
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
curl --header "Content-Type: application/json" --request POST --data '{"id":"testGauge","type":"gauge","value":465529.39165260154}' http://localhost:8080/update/
```
```bash
curl --header "Content-Type: application/json" --request POST --data '{"id":"testGauge","type":"gauge"}' http://localhost:8080/value/
```


----------------
```bash
curl --header "Content-Type: application/json" --request POST --data '{"id":"testCounter1","type":"counter","delta":10}' http://localhost:8080/update/
```
```bash
curl -v --header "Content-Type: application/json" --request POST --data '{"id":"testCounter1","type":"counter"}' http://localhost:8080/value/
```
```bash
curl -v --header "Accept-Encoding: gzip" --header "Content-Type: application/json" --request POST --data '{"id":"testCounter1","type":"counter"}' http://localhost:8080/value/ --compressed
```

```bash
curl -v --request GET  http://localhost:8080/
```

```bash
curl -v --header "Accept-Encoding: gzip" --request GET  http://localhost:8080/ --compressed
```




how to install golint
https://command-not-found.com/golint

```bash
golint ./...
```

iter 2
```bash
go build -o agent *.go
```



sudo apt install golang-easyjson
easyjson -all /home/alex/Dev/GolandYandex/metrics/internal/handlers/v2/metrics_handlers.go 


