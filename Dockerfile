FROM debian:stable-slim
# COPY source destination
COPY chirpy /bin/chirpy
COPY /.env .
CMD ["/bin/chirpy"]