Open a terminal and run the following command to generate a private key:
```bash
openssl genpkey -algorithm RSA -out private_key.pem -pkeyopt rsa_keygen_bits:2048
```

Once you have the private key, you can generate the corresponding public key with the following command:
```bash
openssl rsa -pubout -in private_key.pem -out public_key.pem
```

launch server
```bash
go run cmd/server/main.go -crypto-key /home/alex/Dev/GolandYandex/metrics/private_key.pem
```


launch agent
```bash
go run cmd/agent/main.go -crypto-key /home/alex/Dev/GolandYandex/metrics/public_key.pem
```