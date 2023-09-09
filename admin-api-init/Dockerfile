FROM registry.suse.com/bci/golang:1.21 as builder
WORKDIR /go/src/builder
COPY . .
RUN go build -ldflags="-s -w" -o ./adventcalendar-api-init

FROM registry.suse.com/bci/bci-micro:15.5
USER root
RUN echo "user:x:10000:10000:user:/home/user:/bin/bash" >> /etc/passwd && mkdir /home/user
USER user
ENV GIN_MODE=release
COPY --from=builder /go/src/builder/adventcalendar-api-init /home/user/app/
WORKDIR /home/user/app
EXPOSE 8080
CMD ["./adventcalendar-api-init"]

