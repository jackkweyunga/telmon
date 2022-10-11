## Build
FROM golang:1.16-buster AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download && \
    go mod vendor

RUN touch telmon.log

COPY . ./

RUN go build -o /telmon

## Deploy
FROM gcr.io/distroless/base-debian10:debug

WORKDIR /

COPY --from=build /telmon /telmon
COPY --from=build /app/telmon.log /telmon.log

EXPOSE 8080

USER nonroot:nonroot

HEALTHCHECK --interval=5s --timeout=5s CMD ["/healthcheck","http://localhost:8080/ping"]

ENTRYPOINT ["/telmon"]

# docker run -d -v /root/.telmon/.telmon-config.yaml:/.telmon-config.yaml -v /root/.telmon/telmon.log:/telmon.log --network host --name telmon  ghcr.io/jackkweyunga/telmon:web

# docker run -d -v ./.telmon-config.yaml:/.telmon-config.yaml -v ./telmon.log:/telmon.log --network host --name telmon  ghcr.io/jackkweyunga/telmon:web
