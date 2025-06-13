FROM golang:1.24 AS build-stage

WORKDIR /app

COPY go.mod go.sum Makefile ./
RUN echo "" > .env
RUN make prepare

COPY . .

RUN make build

FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /app/spot-instrument-service /spot-instrument-service

USER nonroot:nonroot
ENTRYPOINT ["/spot-instrument-service"]