FROM golang:1.19.4

ENV ADDRESS=0.0.0.0:9501

WORKDIR /build
COPY . .

RUN go get
RUN go build -o app
RUN mkdir /app
RUN cp app /app/
RUN cp views/dist /app/views/dist
RUN rm -rf /build

EXPOSE 9501

CMD ["/app/app"]
