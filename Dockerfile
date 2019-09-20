FROM golang:1.12.0-alpine3.9
RUN mkdir /app
ADD . /app
ADD . /go/src/github.com/ajdinahmetovic/item-service
WORKDIR /app/server
RUN apk update && apk add git && go get github.com/dgrijalva/jwt-go && go get -d github.com/lib/pq && go get golang.org/x/crypto/bcrypt && go get google.golang.org/grpc && go get go.uber.org/zap
CMD go run .