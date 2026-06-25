FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod main.go ./
RUN go build -o server .
COPY job/ ./job/
RUN go build -o job ./job/

FROM alpine:3.21
ARG GIT_COMMIT=unknown
ENV GIT_COMMIT=$GIT_COMMIT
RUN apk --no-cache add ca-certificates wget
WORKDIR /app
COPY --from=builder /app/server .
COPY --from=builder /app/job .
EXPOSE 8080
CMD ["./server"]
