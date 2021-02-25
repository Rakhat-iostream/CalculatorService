FROM golang

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64
COPY ./greet/greet_server /app/greet/greet_server
COPY ./greet/greet.pb /app/greet/greet.pb
WORKDIR /app
RUN go build /app/greet/server && chmod +x /app/greet/server
ENTRYPOINT ["./server"]
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /app .
CMD ["./server"]
EXPOSE 50051