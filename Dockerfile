## Build
FROM golang:1.16-buster AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download && \
    go mod vendor

COPY . ./

RUN go build -o /telmon

## Deploy
FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /telmon /telmon

EXPOSE 8080

USER nonroot:nonroot

HEALTHCHECK --interval=5s --timeout=5s CMD ["/healthcheck","http://localhost:8080/ping"]

ENTRYPOINT ["/telmon"]
