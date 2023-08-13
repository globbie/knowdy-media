# syntax=docker/dockerfile:1

FROM golang:1.21.0-alpine3.18 as builder

WORKDIR /build
ADD internal/ /build/internal/
COPY cmd/webgate/main.go .

COPY go.mod ./
RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build -o /build/knd-media

RUN addgroup -S knd && adduser -S knd -G knd

FROM scratch

COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
USER knd

COPY --from=builder /build/knd-media .

EXPOSE 8080
CMD ["/knd-media", "--listen-address=0.0.0.0:8080"]
