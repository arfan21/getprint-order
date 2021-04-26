FROM golang:1.16.3
WORKDIR /app
COPY . .
RUN go get -u github.com/cosmtrek/air
CMD ["air"]