FROM golang:alpine as builder

WORKDIR /app

# Get Reflex for live reload in dev
ENV GO111MODULE=on
RUN go get github.com/cespare/reflex

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o ./main ./app/main.go

FROM alpine:latest
WORKDIR /root/

#Copy executable from builder
COPY --from=builder /app/main .

EXPOSE 9090
CMD ["./main"]
