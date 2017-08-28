FROM golang:1.8.3-alpine
RUN apk add --no-cache git gcc musl-dev
RUN mkdir -p /tmp/data
COPY data.zip /tmp/data/data.zip
WORKDIR /go/src/github.com/elBroom/highloadCup
ADD . .
RUN go build -a -o app_ .
EXPOSE 80
CMD ["./app_"]

# docker build -t elbroom/highloadcup .
# docker push elbroom/highloadcup