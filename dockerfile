FROM alpine:latest
EXPOSE 5050
COPY . /
CMD ["/excel"]
