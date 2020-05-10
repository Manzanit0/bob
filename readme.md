# Bob (The builder)

Build a repo in the cloud, download the binary.

## Getting started

To use the CLI:

```
go run cmd/cli/main.go -repository=https://github.com/Manzanit0/golarm -out="/Users/manzanit0/Desktop"
```

To run the web server:

```
go run cmd/web/main.go
```

Then you can go to `locahost:8080` for the UI or use curl to access the actual endpoint:

```
curl localhost:8080/build -X POST -d '{"url":"https://github.com/Manzanit0/golarm", "entry_point": "."}' --output golarm
```

That should download the executable in the current directory. You might have to
give it executable permissions though. `chmod +x <executable>` should do the
trick.

## Deploying the application

The application is hosted on Heroku. According to the [instructions](https://devcenter.heroku.com/articles/container-registry-and-runtime))
to deploy, run the following commands:

```bash
heroku container:login
heroku container:push web
heroku container:release web
heroku open
```

To run the container locally:

```
docker build -t bob .
docker run -p 8080:8080 -e PORT=8080 bob
```
