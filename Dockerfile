# syntax=docker/dockerfile:1

# Build the application from source
FROM golang:1.21 AS build-stage


WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN GOOS=linux CGO_ENABLED=1 go build -o /cdive

# Run the tests in the container
FROM build-stage AS run-test-stage

RUN go test -v ./...

# Deploy the application binary into a lean image
FROM build-stage AS build-release-stage

WORKDIR /build

COPY --from=build-stage /cdive ./
COPY views views
COPY config.yaml config.yaml

EXPOSE 8000

ENTRYPOINT ["./cdive", "start"]
