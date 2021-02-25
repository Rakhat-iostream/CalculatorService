FROM golang
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64
COPY ./greet/greet_client /app/greet/greet_client
COPY ./greet/greet.pb /app/greet/greet.pb
WORKDIR /app
RUN go build /app/greet/client && chmod +x /app/greet/client
ENTRYPOINT ["./client"]
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /app .
CMD ["./client"]
EXPOSE 50051