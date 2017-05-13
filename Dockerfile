FROM alpine
EXPOSE 7000
COPY ./bin/goproxy-linux-amd64 /app/goproxy
CMD /app/goproxy remote
