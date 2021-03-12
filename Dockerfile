FROM scratch

ADD ca-certificates.crt /etc/ssl/certs/
ADD bin/analytics /
ADD db /db

EXPOSE 8080
EXPOSE 8081

ENTRYPOINT ["/analytics"]
