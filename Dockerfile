FROM golang:1.22
LABEL authors="dalinarkholin"
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

ENV PORT=8888
ENV CONNECTION_STRING="mongodb+srv://pojebiemnie:pojebiemnie@simpledb.swlqbjl.mongodb.net/?retryWrites=true&w=majority&appName=simpleDB"


RUN go build -o ./out/server .
CMD ./out/server
