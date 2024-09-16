FROM golang:1.22.0

ENV TODO_PORT=7540
ENV TODO_PASSWORD=12345
ENV TODO_DBFILE=./scheduler.db

WORKDIR /todolist

COPY go.mod go.sum ./
RUN go mod tidy
RUN go mod download

COPY . .

EXPOSE 7540

RUN go build -o /todolist_app ./cmd/server

CMD ["/todolist_app"]
