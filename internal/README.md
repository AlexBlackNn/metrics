# internal

В данной директории и её поддиректориях будет содержаться имплементация вашего сервиса

go run ./cmd/server/main.go  --config="./cmd/server/config/demo.yaml"

```bash
curl -v -H "Content-Type: text/plain" -X POST  http://localhost:8076/update/gauge/param1/2
```

```bash
curl -v -H "Content-Type: text/plain" -X POST  http://localhost:8076/update/gauge1/param1/2
```