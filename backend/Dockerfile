FROM golang:1.22

WORKDIR /app
COPY main.go .
RUN go mod init backend
RUN go get github.com/prometheus/client_golang/prometheus
RUN go get github.com/prometheus/client_golang/prometheus/promhttp
RUN go build -o backend .

CMD ["./backend"]