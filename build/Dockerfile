# Собираем в гошке
FROM golang:1.16.6 as build

ENV BIN_FILE /opt/rotation/banner-rotation-app
ENV CODE_DIR /go/src/

WORKDIR ${CODE_DIR}

# Кэшируем слои с модулями
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . ${CODE_DIR}

# Собираем статический бинарник Go (без зависимостей на Си API),
# иначе он не будет работать в alpine образе.
ARG LDFLAGS
RUN CGO_ENABLED=0 go build \
        -ldflags "$LDFLAGS" \
        -v -o ${BIN_FILE} ./cmd/

# На выходе тонкий образ
FROM alpine:3.9

LABEL ORGANIZATION="My organization"
LABEL SERVICE="banner-rotation"
LABEL MAINTAINERS="rifkinsm@gmail.com"

ENV BIN_FILE "/opt/rotation/banner-rotation-app"
COPY --from=build ${BIN_FILE} ${BIN_FILE}

ENV CONFIG_FILE /etc/banner-rotation/config.yaml
COPY configs/config.yaml ${CONFIG_FILE}

CMD ${BIN_FILE} --config=${CONFIG_FILE}
