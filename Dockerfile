FROM scratch

ADD build/bin/linux/amd64/prometheus-exporter /prometheus-exporter
ADD ./configs/prometheus-exporter.yml /config/application.yml

CMD ["/prometheus-exporter"]
