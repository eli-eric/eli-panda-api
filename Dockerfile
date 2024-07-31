# syntax=docker/dockerfile:1

FROM golang:1.22-alpine as build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN go install github.com/swaggo/swag/cmd/swag@latest

RUN swag init -g swagger_prod.go 
RUN cp -r ./docs/swagger.yaml ./open-api-specification/panda-api.yaml

RUN go build -v -o panda-api -ldflags "-s -w"

## Deploy
FROM alpine:latest

RUN apk add --no-cache tzdata
ENV TZ=Europe/Prague

WORKDIR /root/

COPY --from=build /app/panda-api ./
COPY --from=build /app/db ./db
COPY --from=build /app/open-api-specification ./open-api-specification

EXPOSE 50000

CMD [ "./panda-api" ]