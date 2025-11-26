# WinboxGet

**WinboxGet** is a tiny Go-powered service that automatically scrapes the official MikroTik
download page and provides **direct, simple, version-aware download links** for
Winbox 4 (Windows/Mac/Linux) and Winbox 3 (Windows x64 and 32-bit).

It exists simply because people over the average internet-age and with higher than average understanding of BGP think 4 clicks are unreasonable to download a file.

It is hosted at https://winboxget.fly.dev

## Run / build

### Prerequisites
- Go 1.25.4 or later

### Running locally

```bash
go run main.go
```

### Building locally
```bash
go build -o winboxget
```

### Fly.io 

Check fly.io
Don't be lazy.