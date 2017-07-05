FROM ubuntu:16.04
EXPOSE 7700

COPY conf/localtime /etc/localtime
COPY conf/timezone /etc/timezone

COPY release/linux/amd64/malcolm /ke/malcolm

ENTRYPOINT ["/ke/malcolm"]
CMD ["server"]

