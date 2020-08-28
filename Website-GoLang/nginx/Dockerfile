FROM alpine:latest

RUN apk update
RUN apk add nginx

COPY nginx.conf /etc/nginx/nginx.conf
#COPY nginx/www /var/www/nginx

RUN mkdir -p /run/nginx
CMD ["/usr/sbin/nginx"]

EXPOSE 80
