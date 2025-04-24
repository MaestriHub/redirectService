FROM --platform=linux/amd64 golang:alpine as build

RUN apk update && apk upgrade

WORKDIR /build

COPY . .

RUN touch .env

WORKDIR /build

RUN go mod tidy

WORKDIR /build/cmd/api

RUN go build

WORKDIR /staging

RUN cp /build/cmd/api/api .
RUN cp /build/.env .
RUN cp -r /build/static .

FROM --platform=linux/amd64 alpine:latest

RUN apk update && apk upgrade

RUN apk add --no-cache shadow

RUN useradd --user-group --create-home --system --skel /dev/null --home-dir /app gin

WORKDIR /app

COPY --from=build --chown=gin:gin /staging /app

ENV GOCACHE=/tmp/go-cache

ENV GIN_MODE=release

USER gin:gin

EXPOSE 8080

ENTRYPOINT ["./api"]

CMD ["run", "--env", "production", "--hostname", "0.0.0.0", "--port", "8080"]