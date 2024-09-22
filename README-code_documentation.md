##  Локальное отображения godoc-документации 
```bash
 godoc -http=:8085
```

##  Скачать документацию в формате HTML
```bash
wget -r -np -N -E -p -k http://localhost:8085/pkg/github.com/AlexBlackNn/metrics/
```

##  Вывод докумнтации для указанного модуля в консоль
```bash
go doc -all github.com/AlexBlackNn/metrics/pkg/storage/postgres
```

## Разрешать запускать примеры из документации
```bash
godoc -http=:8085 -play
```
http://localhost:8085/pkg/bytes/#example_Buffer

## По умолчанию godoc не отображает пакеты, расположенные в поддиректориях internal. 
Чтобы увидеть служебные пакеты, добавьте в браузере 
```
http://localhost:8080/pkg/?m=all
```

## swagger
```bash
swag init -g ./cmd/server/main.go -o ./cmd/server/docs
```

http://localhost:8080/swagger/index.htm/index.html