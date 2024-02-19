FROM golang:1.21

WORKDIR /wokkibot

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./
COPY commands/*.go ./commands/
COPY config/*.go ./config/
COPY utils/*.go ./utils/

RUN go build -o /wokkibot

CMD [ "./wokkibot" ]