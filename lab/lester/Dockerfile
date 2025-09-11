FROM golang:1.24.2

WORKDIR /app

COPY . .

RUN go mod tidy
RUN go build -o lester .

EXPOSE 50051

CMD ["./lester"]