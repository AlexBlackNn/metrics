### Сервис сбора метрик

# Оптимизация памяти
[README-optimization.md](README-optimization.md)

# Документация
[README-code_documentation.md](README-code_documentation.md)

# Компиляция сервера 
`go build -ldflags "-s -w"` — скомпилирует исполняемый файл меньшего размера, так как в него не будет включена таблица символов и отладочная информация.

```bash
go build -ldflags '-s -w' -o cmd/server/server cmd/server/main.go
```
