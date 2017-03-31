FROM scratch

ADD ./cloud-cli /
ENTRYPOINT ["/cloud-cli"]