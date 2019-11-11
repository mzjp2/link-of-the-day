FROM golang

COPY . /app
WORKDIR /app

RUN go get && go build -o bin/link-of-the-day .

ENV PORT 8080

CMD ["./bin/link-of-the-day"]
