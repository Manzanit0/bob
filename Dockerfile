FROM golang:1.14

WORKDIR /bob

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN GOOS=linux GOARCH=amd64 go build -o bob cmd/web/main.go

# We need to export these because bob uses them for compilation
ENV GOPATH=/go
ENV GOBIN=/go/bin

# Standard Heroku: use $PORT env
ENV PORT=$PORT

ENV GIN_MODE=release
ENTRYPOINT ["/bob/bob", "-html", "/bob/cmd/web/static"]
