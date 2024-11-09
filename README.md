### Сервис сбора метрик

# Оптимизация памяти
[README-optimization.md](README-optimization.md)

# Документация
[README-code_documentation.md](README-code_documentation.md)

# Компиляция сервера 
`go build -ldflags "-s -w"` — скомпилирует исполняемый файл меньшего размера, так как в него не будет включена таблица символов и отладочная информация.

```bash
go build -ldflags '-s -w -X main.buildVersion=1.0.0 -X main.buildDate=2023-01-23 -X main.buildCommit=0c2fs'  -o cmd/agent/agent cmd/agent/main.go
go build -ldflags '-s -w -X main.buildVersion=1.0.0 -X main.buildDate=2023-01-23 -X main.buildCommit=0c2fs' -o cmd/server/server cmd/server/main.go
```
