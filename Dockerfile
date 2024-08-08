FROM golang:1.22
LABEL authors="dalinarkholin"
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

ENV CONNECTION_STRING="mongodb+srv://pojebiemnie:pojebiemnie@simpledb.swlqbjl.mongodb.net/?retryWrites=true&w=majority&appName=simpleDB"

EXPOSE 8080

RUN go build -o ./out/server .
CMD ./out/server
