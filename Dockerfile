FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o greenfra ./src/main.go

FROM alpine:latest

WORKDIR /root/

RUN apk add --no-cache curl unzip

RUN curl -LO "https://releases.hashicorp.com/terraform/1.5.7/terraform_1.5.7_linux_amd64.zip" && \
    unzip terraform_1.5.7_linux_amd64.zip && \
    mv terraform /usr/local/bin/ && \
    chmod +x /usr/local/bin/terraform && \
    rm terraform_1.5.7_linux_amd64.zip

COPY --from=builder /app/greenfra .

RUN chmod +x ./greenfra

ENV PATH="/root:${PATH}"

CMD ["/bin/sh"]
