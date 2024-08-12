FROM golang:1.22.2 AS builder


WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . .

RUN go mod tidy

RUN go build -o main .

FROM debian:bookworm-slim


WORKDIR /app


COPY --from=builder /app/main .

COPY templates /app/templates
COPY banners /app/banners

EXPOSE 8080

# Apply metadata
LABEL org.opencontainers.image.title="Go Web Server" \
      org.opencontainers.image.description="Ascii Art Web" \
      org.opencontainers.image.version="1.0" \
      org.opencontainers.image.author="MeFerdi, jowala"

CMD ["./main"]