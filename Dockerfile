FROM golang:1.23.4
ENV TODO_PORT=7540
ENV TODO_DBFILE=./base/scheduler.db
ENV TODO_PASSWORD=1278
EXPOSE $TODO_PORT
WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./
COPY ./pkg ./pkg
COPY ./tests ./tests
COPY ./web ./web
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /go_final_project
CMD ["/go_final_project"]