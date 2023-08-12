FROM golang:1.21-bullseye

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . ./
RUN go build -v -o /arxiv_server
CMD [ "/arxiv_server" ]

# CMD ["go", "run", "main.go"]
