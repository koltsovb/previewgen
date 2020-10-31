FROM golang:latest 
RUN mkdir /app 
RUN mkdir /tmp/pg
ADD . /app/ 
WORKDIR /app 

ENV PORT 8080
EXPOSE $PORT

RUN go build -o main . 
CMD ["/app/main", "-p=8080", "-d=/tmp/pg"]