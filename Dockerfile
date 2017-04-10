FROM ubuntu:16.04

WORKDIR /opt/cloud_mail_ru
ADD cmri ./
CMD ["./cmri"]