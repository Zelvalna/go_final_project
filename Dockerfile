FROM golang:1.22.0

ENV TODO_PORT=7540
ENV TODO_PASSWORD=12345
ENV TODO_DBFILE=/todolist/scheduler.db

WORKDIR /todolist

COPY . .

RUN go mod download

RUN go build -o /todolist_app ./cmd/server

CMD ["/todolist_app"]
