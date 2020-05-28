FROM alpine
LABEL maintainer="gideonhacer@gmail.com"
RUN apk update && \
   apk add ca-certificates && \
   update-ca-certificates && \
   rm -rf /var/cache/apk/* && \
   apk add libc6-compat
EXPOSE 80 443
WORKDIR /app
COPY dist dist
COPY static static
COPY api api
COPY gateway .
ENTRYPOINT [ "/app/gateway", "--env" ]