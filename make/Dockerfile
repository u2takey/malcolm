FROM centurylink/ca-certs
EXPOSE 7700

COPY conf/localtime /etc/localtime
COPY conf/timezone /etc/timezone

COPY release/linux/amd64/malcolm /malcolm

ENTRYPOINT ["/malcolm"]
CMD ["server"]

