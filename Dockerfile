FROM golang:1.15

ENV GO111MODULE=on
RUN go get -u github.com/jinykim0x80/42api_gateway
COPY ./ ./src
WORKDIR src
RUN go build -o 42api_gateway 

EXPOSE 443
EXPOSE 80 

CMD ["./42api_gateway"]
