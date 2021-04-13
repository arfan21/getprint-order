# Please keep up to date with the new-version of Golang docker for builder
FROM golang:1.16.3 AS build
WORKDIR /go/src/github.com/arfan21/getprint-order/
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server .

FROM alpine:3.12
COPY --from=build /go/src/github.com/arfan21/getprint-order/server .
CMD ["./server"]