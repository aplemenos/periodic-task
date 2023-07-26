FROM golang:1.20 AS builder

RUN mkdir /periodictask
ADD . /periodictask
WORKDIR /periodictask

RUN CGO_ENABLED=0 GOOS=linux go build -o periodictask cmd/periodic-task/main.go

FROM alpine:latest AS production
COPY --from=builder /periodictask .

RUN apk --no-cache add curl
RUN apk add --no-cache tzdata

CMD ["./periodictask"]