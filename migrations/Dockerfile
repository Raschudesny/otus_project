#Сервис для накатывания миграций гусем во время старта banner-rotation сервиса через dokcer-compose
FROM golang:1.16.6 as build
ENV MIGRATIONS_DIR /migrations
RUN go get -u github.com/pressly/goose/v3/cmd/goose
WORKDIR ${MIGRATIONS_DIR}
COPY . ${MIGRATIONS_DIR}
ADD https://raw.githubusercontent.com/vishnubob/wait-for-it/master/wait-for-it.sh ./wait
RUN chmod +x ./wait
CMD ./wait db:5432 --timeout=60 && goose -dir /migrations postgres "host=net-db user=danny password=danny dbname=rotation sslmode=disable" up
