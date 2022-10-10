FROM golang:1.19-alpine3.16 as build
WORKDIR /go/src/build
COPY . .
RUN mkdir -p dist \
    && go mod vendor \
    && go build -o dist/compose-metrics .

FROM alpine:3.16.2

EXPOSE 10000
COPY --from=build /go/src/build/dist/ /usr/local/bin/

ENTRYPOINT ["/usr/local/bin/compose-scaler"]