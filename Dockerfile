FROM golang:1.22 AS build

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o bot main.go

FROM alpine

COPY --from=build /build/bot /bin/bot

CMD ["/bin/bot"]