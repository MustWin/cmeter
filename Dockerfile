FROM frolvlad/alpine-glibc

ENV CMETER_CONFIG_PATH /etc/cmeter.default.yml

COPY ./dist /bin/cmeter
COPY ./config.default.yml /etc/cmeter.default.yml

ENTRYPOINT ["/bin/cmeter"]
