FROM --platform=linux/amd64 golang:alpine as build

# Install OS updates
RUN apk update && apk upgrade

# Set up a build area
WORKDIR /build

# # Copy entire repo into container
COPY . .

WORKDIR /build

RUN go mod tidy

RUN go build

WORKDIR /staging

RUN cp /build/redirectServer .
RUN cp /build/.env .
RUN cp -r /build/static .

# # ================================
# # Run image
# # ================================

FROM --platform=linux/amd64 alpine:latest

# # Make sure all system packages are up to date, and install only essential packages.
RUN apk update && apk upgrade

RUN apk add --no-cache shadow

# # Create a vapor user and group with /app as its home directory
RUN useradd --user-group --create-home --system --skel /dev/null --home-dir /app gin

# Switch to the new home directory
WORKDIR /app

COPY --from=build --chown=gin:gin /staging /app

# # Provide configuration needed by the built-in crash reporter and some sensible default behaviors.
ENV GOCACHE=/tmp/go-cache

ENV GIN_MODE=release

# # Ensure all further commands run as the app user
USER gin:gin

# # Let Docker bind to port 8080
EXPOSE 8080

# # Start the service when the image is run, default to listening on 8080 in production environment

ENTRYPOINT ["./redirectServer"]

CMD ["run", "--env", "production", "--hostname", "0.0.0.0", "--port", "8080"]