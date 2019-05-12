# openssl commands

```shell
openssl genrsa -out example-key.pem 2048
openssl req -subj '/CN=*/' -new -x509 -key example-key.pem -out example-cert.pem -days 1095
```

NOTE: Set MSYS_NO_PATHCONV=1 in case you're using git bash console or msys2 on windows...
