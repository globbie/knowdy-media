FROM golang:1.16.2-alpine3.13 as builder

WORKDIR /build
ADD . .

RUN ls -l

RUN CGO_ENABLED=0 GOOS=linux go build -a -o knd-media cmd/main.go

RUN addgroup -S knd && adduser -S knd -G knd

FROM scratch

COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
USER knd

COPY --from=builder /build/knd-media .

EXPOSE 8080
CMD ["/knd-media", "--listen-address=0.0.0.0:8069"]
