FROM golang:1.18-alpine as build
WORKDIR /
COPY . .
RUN apk update && apk upgrade
RUN apk add --no-cache sqlite && apk add --no-cache build-base
RUN go mod download
RUN GOOS=linux go build -o elotus ./cmd/main.go

# This is for documentation purposes only.
# To actually open the port, runtime parameters
# must be supplied to the docker command.
EXPOSE 8080

FROM alpine:latest
WORKDIR /
COPY --from=build /elotus ./elotus
ENTRYPOINT ["./elotus"]