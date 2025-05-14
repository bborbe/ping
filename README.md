# Better Ping

This service allow ping multiple targets and show lost packages.

## Usage

```bash
sudo go run main.go 8.8.8.8 8.8.4.4 193.101.111.10 192.168.177.1 192.168.180.1 192.168.178.5
```

```bash
go build -o ping .
sudo ./ping 8.8.8.8 8.8.4.4 193.101.111.10 192.168.177.1 192.168.180.1 192.168.178.5
```
