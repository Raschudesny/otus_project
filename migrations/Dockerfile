#Накатываем миграции гусем и ждем пока база и редис встанут
FROM golang:1.16.6 as build
ENV MIGRATIONS_DIR /migrations
RUN go get -u github.com/pressly/goose/v3/cmd/goose
WORKDIR .
COPY migrations ${MIGRATIONS_DIR}
ADD https://github.com/ufoscout/docker-compose-wait/releases/download/2.9.0/wait /wait
RUN chmod +x /wait
CMD /wait && goose -dir ${MIGRATIONS_DIR} postgres "host=net-db user=danny password=danny dbname=rotation sslmode=disable" up
