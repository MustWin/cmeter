FROM frolvlad/alpine-glibc

ENV CMETER_CONFIG_PATH /etc/cmeter.config.yml

COPY ./config.default.yml /etc/cmeter.config.yml
COPY ./dist /bin/cmeter

ENTRYPOINT ['/bin/cmeter']
