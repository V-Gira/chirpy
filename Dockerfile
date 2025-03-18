FROM debian:stable-slim
# COPY source destination
COPY chirpy /bin/chirpy
COPY /.env .
ENV PORT=8080
CMD ["/bin/chirpy"]