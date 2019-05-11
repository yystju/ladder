# openssl commands

```shell
openssl genrsa -out example-key.pem 2048
openssl req -subj '/CN=*/' -new -x509 -key example-key.pem -out example-cert.pem -days 1095
```