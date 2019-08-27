FROM golang:latest 
RUN mkdir /app 
ADD . /app/ 
WORKDIR /app 
EXPOSE 8086
RUN go get github.com/paypal/gatt
RUN go get github.com/influxdata/influxdb1-client/v2
RUN go build -o main . 
CMD ["/app/main"]