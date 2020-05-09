# Bob (The builder)

Build a repo in the cloud, download the binary.

## Getting started (CLI)

```
go run cmd/cli/main.go -repository=https://github.com/Manzanit0/golarm -out="/Users/manzanit0/Desktop"
```

## Getting started (web)

To run the server

```
go run cmd/web/main.go
```

To make the request:

```
curl localhost:8080/build -X POST -d '{"url":"https://github.com/Manzanit0/golarm", "entry_point": "."}' --output golarm
```

That should download the executable in the current directory. You might have to give it executable permissions though.
