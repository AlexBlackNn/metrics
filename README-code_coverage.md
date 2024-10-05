```bash
go test -v ./... 
```

```bash
go test -v -coverpkg=./... -coverprofile=profile.cov ./...
go tool cover -func profile.cov
```

```bash
go test -v ./... -coverprofile profile.out
go tool cover -func profile.out
go tool cover -o coverage.html -html=profile.out; sed -i 's/black/whitesmoke/g' coverage.html; sensible-browser coverage.html
```
