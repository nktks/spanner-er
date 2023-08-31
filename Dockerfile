FROM golang:1.14-alpine3.11 as builder

WORKDIR /go/src/github.com/nktks/spanner-er

COPY go.mod go.sum  ./
RUN apk add --no-cache git graphviz ttf-freefont &&\
    go mod download

COPY . .
RUN go build -i -o /bin/spanner-er  ./

ENTRYPOINT ["/bin/spanner-er"]
