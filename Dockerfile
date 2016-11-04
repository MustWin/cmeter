FROM frolvlad/alpine-glibc

COPY ./dist /bin/cmeter

WORKDIR /

ENTRYPOINT ['/bin/cmeter']
