FROM docker.elastic.co/beats/filebeat:6.2.4
MAINTAINER Michael Dai <jdai1@xiaohongshu.com>

USER root
COPY filebeat /usr/share/filebeat

RUN chmod 0750 /usr/share/filebeat/filebeat && \
    chown root:filebeat /usr/share/filebeat/filebeat
USER filebeat
