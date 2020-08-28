FROM scratch

ADD /bin/gserve /gserver
ADD src/gserve/templates/ /templates/
ENTRYPOINT ["/gserver"]

EXPOSE 9090
