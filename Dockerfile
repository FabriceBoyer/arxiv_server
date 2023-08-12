ARG GO_VERSION=1.21

FROM golang:${GO_VERSION}-alpine as builder

WORKDIR /app
COPY ./ ./

RUN go mod download
RUN go build -v -o /arxiv_server

#################################################

FROM scratch

COPY --from=builder /arxiv_server /
COPY ./static /static
COPY ./.env.example /.env

CMD [ "/arxiv_server" ]


