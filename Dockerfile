FROM frolvlad/alpine-glibc

#ENV GOPATH $HOME/go
#ENV PKG github.com/MustWin/cmeter
#ENV SRCPATH $GOPATH/src/$PKG

#RUN apk add --update alpine-sdk
COPY ./dist /bin/cmeter
WORKDIR /

#RUN make clean && make compile && cp $SRCPATH/cmeter /bin && rm -rf $SRCPATH

CMD /bin/cmeter agent
