FROM golang:1.20 as build

ADD . /app
WORKDIR /app

RUN CGO_ENABLED=0 GOOS=linux go build

FROM scratch
COPY --from=build /app/klum /klum
CMD ["/klum"]
