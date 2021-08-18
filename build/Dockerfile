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

RUN CGO_ENABLED=0 go build \
        -ldflags "-X 'main.release=develop' -X 'main.buildDate=$(date -u +%Y-%m-%dT%H:%M:%S)' -X 'main.gitHash=$(git log --format="%h" -n 1)'" \
        -v -o ${BIN_FILE} ./cmd/

# На выходе тонкий образ
FROM alpine:3.9

LABEL ORGANIZATION="My organization"
LABEL SERVICE="banner-rotation"
LABEL MAINTAINERS="rifkinsm@gmail.com"

ENV BIN_FILE "/opt/rotation/banner-rotation-app"
COPY --from=build ${BIN_FILE} ${BIN_FILE}

ARG CONFIG_PATH
ENV APP_CONFIG /etc/banner-rotation/config.yaml
COPY ${CONFIG_PATH} ${APP_CONFIG}

CMD ${BIN_FILE} --config=${APP_CONFIG}

