FROM debian:bookworm-slim

RUN apt-get update && apt-get -y --no-install-recommends install \
    supervisor \
    && rm -rf /var/lib/apt/lists/*

COPY supervisord.conf /etc/supervisor/supervisord.conf
COPY ../build/ /opt/sup-agent/

CMD ["/usr/bin/supervisord", "-c", "/etc/supervisor/supervisord.conf"]