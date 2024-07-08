### Сервис сбора метрик

** Literature **
RESTY USEFUL
https://www.alldevstack.com/ru/go-resty-tutorial/go-resty-quickstart.html

how to install golint
https://command-not-found.com/golint

easy-json commands
```
sudo apt install golang-easyjson
easyjson -all /home/alex/Dev/GolandYandex/metrics/internal/handlers/v2/metrics_handlers.go 
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
curl --header "Content-Type: application/json" --request POST --data '{"id":"testGauge","type":"gauge","value":465529.39165260154}' http://localhost:8080/update/
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
curl -v --header "Accept-Encoding: gzip" --header "Content-Type: application/json" --request POST --data '{"id":"testCounter14","type":"counter"}' http://localhost:8080/value/ --compressed
```

```bash
curl -v --request GET  http://localhost:8080/
```

```bash
curl -v --header "Accept-Encoding: gzip" --request GET  http://localhost:8080/ --compressed
```