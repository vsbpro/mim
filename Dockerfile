FROM busybox:latest

RUN mkdir /app

COPY ./mim.exe /app/mim.exe

CMD ["/app/mim.exe"]
