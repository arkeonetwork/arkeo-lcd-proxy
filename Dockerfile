FROM golang:1.22

WORKDIR /app
COPY go.mod ./
COPY main.go ./
RUN go build -o lcd-proxy .

EXPOSE 1318
ENV LISTEN=:1318
CMD ["./lcd-proxy"]
