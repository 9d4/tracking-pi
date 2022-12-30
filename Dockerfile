FROM golang:1.19.4

ENV ADDRESS=0.0.0.0:9501

WORKDIR /build
COPY . .

RUN go get
RUN go build -o app
RUN mkdir /app
RUN cp app /app/
RUN cp .env /app/

RUN apt install curl
RUN curl -fsSL https://deb.nodesource.com/setup_lts.x | bash -
RUN apt install nodejs -y
RUN cd views && npm install && npm run build
RUN cp -r /build/views /app/

WORKDIR /app
RUN rm -rf /build

EXPOSE 9501

CMD ["/app/app"]
