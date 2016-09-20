FROM frolvlad/alpine-glibc

COPY ./dist /bin/cmeter

WORKDIR /

CMD /bin/cmeter agent
