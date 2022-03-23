FROM golang
MAINTAINER Stepan K. <xamust@gmail.com>
WORKDIR /server/
VOLUME ["/opt/server"]
COPY . ./
RUN make build
EXPOSE 8080
EXPOSE 8081
EXPOSE 8082
EXPOSE 5432
CMD [ "build/server" ]