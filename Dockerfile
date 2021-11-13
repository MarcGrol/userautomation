FROM golang:1.16

WORKDIR /src
COPY go.* .

RUN go mod download

COPY . /src
RUN go build -o userautomation

ENTRYPOINT ["userautomation"]
