FROM golang:1.22 AS build

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -ldflags="-X 'main.version=${VERSION}'-w -s" -o bot main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

COPY --from=build /build/bot /bin/bot

RUN chmod +x /bin/bot

CMD ["/bin/bot"]