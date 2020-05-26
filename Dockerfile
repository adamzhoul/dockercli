FROM registry.vipkid.com.cn/base_env/build_ratel_golang:v1.13.5 as go-builder
WORKDIR /code
COPY go.mod go.sum /code/
RUN go version \
    && go env -w GOPROXY=https://goproxy.cn,direct \
    && go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o microctl .

FROM alpine
RUN apk add --no-cache tzdata
ENV TZ=Asia/Shanghai
COPY --from=go-builder /code/microctl /app/
COPY --from=go-builder /code/configs /app/
WORKDIR /app
ENTRYPOINT ["/app/microctl"]
