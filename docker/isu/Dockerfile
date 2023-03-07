FROM golang:1.20 as build

WORKDIR /app
RUN go get -u github.com/uroborosq/isu-itmo/internal/isu
RUN go build -o service

FROM alpine:3

COPY --from=build /app/service /service

ENTRYPOINT [ "/xaosbot" ]

