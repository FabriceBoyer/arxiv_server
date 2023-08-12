ARG GO_VERSION=1.21

FROM golang:${GO_VERSION}-alpine as builder

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY *.go ./
COPY ./arxiv ./arxiv
RUN go build -v -o /arxiv_server

#################################################

# FROM scratch
FROM gcr.io/distroless/static AS final

COPY --from=builder /arxiv_server /
COPY ./static /static
COPY ./.env.example /.env

CMD [ "/arxiv_server" ]


