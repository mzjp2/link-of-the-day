FROM golang

COPY . /app
WORKDIR /app

RUN go get && go build -o cmd/link .

ENV PORT 8080

CMD ["./cmd/link"]
