FROM golang:1.20 as build

ADD . /app
WORKDIR /app

ARG VERSION
ARG COMMIT

RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-X main.Version=$VERSION -X main.GitCommit=$COMMIT"

FROM scratch
COPY --from=build /app/klum /klum
CMD ["/klum"]
