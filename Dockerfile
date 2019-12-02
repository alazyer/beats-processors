ARG BEATS_VERSION=6.7.2
ARG GO_VERSION=1.10.8
ARG GO_PLATFORM=linux-amd64
ARG BEAT_NAME=filebeat

FROM docker.elastic.co/beats/filebeat:${BEATS_VERSION}

USER root
RUN mkdir -p /usr/share/filebeat/plugins/
COPY ./processors-modify-fields.so /usr/share/filebeat/plugins/
RUN chown -R filebeat:filebeat /usr/share/filebeat/plugins/
USER filebeat

CMD ["filebeat", "-e", "--plugin", "/usr/share/filebeat/plugins/processors-modify-fields.so", "-v"]