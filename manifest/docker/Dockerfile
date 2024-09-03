###############################################################################
#                                BUILD
###############################################################################
FROM golang:1.21.8 as builder

ARG BRANCH

WORKDIR /root
RUN git clone -b $BRANCH https://github.com/junqirao/api-gateway.git
WORKDIR /root/api-gateway
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags '-w -s' -o server

###############################################################################
#                                INSTALLATION
###############################################################################
FROM alphine:3.19.1

COPY --from=builder /root/api-gateway/server                             /app/server
COPY --from=builder /root/api-gateway/manifest/docker/entrypoint.sh      /entrypoint.sh

###############################################################################
#                                   START
###############################################################################
WORKDIR /app
RUN chmod +x /entrypoint.sh && chmod 775 /entrypoint.sh

EXPOSE 8000
ENTRYPOINT /entrypoint.sh
