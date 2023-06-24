FROM golang:1.18-alpine

WORKDIR /candigator-api

COPY . .

RUN apk update && apk add git

RUN apk add poppler-utils

RUN go mod download

RUN go build -o ./bin/candidate-tracker-go

EXPOSE 9000

CMD [ "./bin/candidate-tracker-go" ]
