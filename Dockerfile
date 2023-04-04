FROM arm32v7/golang:1.20.1-alpine3.17 as gobuild

WORKDIR /app
RUN apk update
RUN apk upgrade
RUN apk add --no-cache ca-certificates
RUN update-ca-certificates
COPY . .
RUN CGO_ENABLED=0 go build -o app .

FROM scratch
WORKDIR /app
COPY --from=gobuild /app/app .
COPY --from=gobuild /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

CMD ["./app"]