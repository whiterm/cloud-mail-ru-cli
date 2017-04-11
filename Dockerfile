FROM ubuntu:16.04

RUN apt-get update && apt-get install -y ca-certificates

ADD cmri /usr/local/bin/
CMD ["cmri"]