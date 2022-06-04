FROM golang
MAINTAINER Stepan K. <xamust@gmail.com>
WORKDIR /server/
VOLUME ["/opt/server"]
COPY . ./

# install psql
#RUN apt-get update
#RUN apt-get -y install postgresql-client
# make wait-for-postgres.sh executable
#RUN chmod +x wait-for-postgres.sh

RUN make build
EXPOSE 8080
EXPOSE 8081
EXPOSE 8082
EXPOSE 5432
CMD [ "build/server" ]