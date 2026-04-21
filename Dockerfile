FROM golang:1.25 AS build

WORKDIR /usr/src/app

COPY . .

RUN go mod tidy && \
	CGO_ENABLED=0 GOOS=linux go build -o typist . 
FROM alpine:3.19
WORKDIR /app
COPY --from=build /usr/src/app/typist /usr/local/bin/typist
RUN adduser -D appuser
USER appuser
ENTRYPOINT ["typist", "--web"]

