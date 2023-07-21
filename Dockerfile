FROM golang:1.20 as build

ADD . /app
WORKDIR /app

ARG VERSION
ARG COMMIT

RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-X main.Version=$VERSION -X main.GitCommit=$COMMIT"
RUN adduser --uid 1000 --disabled-password klum-user

FROM scratch
COPY --from=build /etc/passwd /etc/passwd
COPY --from=build /app/klum /klum
USER klum-user
CMD ["/klum"]
