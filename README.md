# HTTP Disruption POC

## About

## Getting Started

### Dependencies

- [Go v1.16](https://golang.org/doc/install)
- [mkcert](https://kifarunix.com/how-to-create-self-signed-ssl-certificate-with-mkcert-on-ubuntu-18-04/)

### Environment

### Creating the Certificates

```bash
mkcert -install
mkcert localhost 127.0.0.1 google.com *.google.com projectreclass.org *.projecreclass.org
```

### Creating IP Table Rules

```bash
iptables -t nat -A OUTPUT -p tcp -m tcp --dport 5000 -j DNAT --to-destination 127.0.0.1:9090
iptables -t nat -A OUTPUT -p tcp -m tcp --dport 4000 -j DNAT --to-destination 127.0.0.1:8080
```

## Running

```bash
go run main.go
```
