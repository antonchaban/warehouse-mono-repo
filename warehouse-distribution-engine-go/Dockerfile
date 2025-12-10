FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git make

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /distributor ./cmd/distributor

FROM gcr.io/distroless/static-debian12

WORKDIR /

COPY --from=builder /distributor /distributor

USER nonroot:nonroot

ENTRYPOINT ["/distributor"]