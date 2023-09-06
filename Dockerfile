FROM golang:1.21-alpine AS build

RUN adduser --uid 1000 --disabled-password klum-user

WORKDIR /app

COPY go.mod go.sum /app/
RUN go mod download

COPY . .

ARG VERSION
ARG COMMIT
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-X main.Version=$VERSION -X main.GitCommit=$COMMIT"

FROM scratch
COPY --from=build /etc/passwd /etc/passwd
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /app/klum /klum
USER klum-user
CMD ["/klum"]
