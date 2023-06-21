#
# Build the application package & install node.js dependencies
#
FROM golang:1.20-alpine as builder

WORKDIR /build

RUN adduser -S app -G users -u 99 -H -D && \
    apk add --no-cache ca-certificates

COPY go.mod .
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -v -o github_backup -ldflags "-s -w" cmd/github_backup/main.go && \
    ls -lah github_backup


#
# Prep the final stage with only the required dependencies
#
FROM scratch
WORKDIR /home/app

COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --from=builder --chown=app:users /build/github_backup ./github_backup
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

USER app

ENTRYPOINT ["./github_backup"]
