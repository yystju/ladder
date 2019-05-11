# README

1. Run socks service (on server)

```
./socks5 -l :11080
```

2. Run server service (on server)

```
./server -d localhost:11080 -l ec2-13-112-82-144.ap-northeast-1.compute.amazonaws.com:2080
```

3. Run client service (on client)

```
./client.exe -s ec2:2080
```
