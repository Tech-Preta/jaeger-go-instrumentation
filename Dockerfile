# Etapa de build
FROM cgr.dev/chainguard/go AS builder

WORKDIR /app

COPY src/go.mod ./
COPY src/go.sum ./
COPY src/main.go ./

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Etapa final: imagem mínima
FROM cgr.dev/chainguard/static

WORKDIR /app

COPY --from=builder /app/main .

ENV PORT=8080

EXPOSE 8080

CMD ["./main"]
