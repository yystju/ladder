# README

1. Run socks service (on server)

```
./socks5 -l :11080
```

2. Run server service (on server)

```
./server -d localhost:11080 -l :2080
```

3. Run client service (on client)

```
./client.exe -s :2080
```
