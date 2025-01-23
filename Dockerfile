FROM golang:1.23.3 AS build

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

COPY . .

ARG GIT_COMMIT
RUN CGO_ENABLED=0 go build -ldflags="-X 'main.version=${GIT_COMMIT}' -w -s" -o bot main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

COPY --from=build /build/bot /bin/bot

RUN apk --no-cache add ca-certificates curl python3 ffmpeg curl

RUN curl -L https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp -o /usr/local/bin/yt-dlp \
    && chmod a+rx /usr/local/bin/yt-dlp

RUN chmod +x /bin/bot

CMD ["/bin/bot"]