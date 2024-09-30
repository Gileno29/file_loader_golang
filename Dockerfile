FROM golang:1.22.5-alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

#DESABILITA OS COPILADORES DO C QUE NÃO ESTÁ PRESENTE NA IMAGEM FINAL
RUN CGO_ENABLED=0 GOOS=linux go build -o fileloader

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app .

EXPOSE 8080

CMD ["./fileloader"]