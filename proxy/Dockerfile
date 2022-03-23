FROM golang
MAINTAINER Stepan K. <xamust@gmail.com>
WORKDIR /proxy/
VOLUME [ "/opt/proxy" ]
COPY . ./
#RUN go build -o build/proxy -v ./cmd/app
RUN make build
EXPOSE 9090
EXPOSE 8080
EXPOSE 8081
EXPOSE 8082
CMD [ "build/proxy" ]
