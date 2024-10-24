Open a terminal and run the following command to generate a private key:
```bash
openssl genpkey -algorithm RSA -out private_key.pem -pkeyopt rsa_keygen_bits:2048
```

Once you have the private key, you can generate the corresponding public key with the following command:
```bash
openssl rsa -pubout -in private_key.pem -out public_key.pem
```