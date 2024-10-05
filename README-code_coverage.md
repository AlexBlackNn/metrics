```bash
go test -v ./... 
```

# coverage without tests and mocks files
```bash
go test -v -coverpkg=./... -coverprofile=profile.cov.tmp ./...
cat profile.cov.tmp  | grep -v "test_" > profile.cov.tmp1
cat profile.cov.tmp1  | grep -v "mock_" > profile.cov.tmp2
cat profile.cov.tmp2  | grep -v "_easyjson.go" > profile.cov
rm profile.cov.tmp1 && rm profile.cov.tmp && rm profile.cov.tmp2
go tool cover -func profile.cov

```

# show files with zero coverage 
```bash
go tool cover -func profile.cov | awk '{if ($NF == "0.0%") print $0}'
```


```bash
go test -v ./... -coverprofile profile.out
go tool cover -func profile.out
go tool cover -o coverage.html -html=profile.out; sed -i 's/black/whitesmoke/g' coverage.html; sensible-browser coverage.html
```
