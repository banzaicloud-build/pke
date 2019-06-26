FROM scratch
EXPOSE 8080
ENTRYPOINT ["/pke"]
COPY ./bin/ /