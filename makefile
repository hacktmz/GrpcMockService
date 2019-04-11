ifndef APPNAME
	APPNAME = mockservice
	EXENAME = ${APPNAME}
	PACKAGE = ${APPNAME}
else
	EXENAME = ${APPNAME}-${SCHEMA}
	PACKAGE = ${APPNAME}-${SCHEMA}-${ENV}
endif

ifeq ($(SCHEMA),default)
	EXENAME = ${APPNAME}
endif
ifndef VERSION
	VERSION = 1.0.0
endif
ifndef RELEASE
	RELEASE = $(shell date +'%Y%m%d%H%M')
endif

DEPVERSION  = 
FULLNAME= "${PACKAGE}-${VERSION}"

export PATH := $(PATH):$(GOPATH)/bin

.PHONY: all distclean clean dist rpm

target := GrpcMockService helloclient   parserclient

all: ${target} 

GrpcMockService: server.go
	go build -v 
	
helloclient:  
	go build -o ./helloclient ./test/helloclient/hello_client.go
	
parserclient:
	go build -o ./parserclient ./test/parserclient/parser_client.go

clean:
	rm -rf buildversion.go pkg bin ${target} ${target}*.gz ${target}*.rpm  src  *.log 

distclean: clean

dist: ${target}
	test -d ${FULLNAME} || mkdir ${FULLNAME}
	test ! -d .git || git log --graph --oneline --pretty=format:"%ad: %an<%ae> %n     %s" --date=rfc >ChangeLog
	cp -a ${EXENAME} pbs packages/scripts/${EXENAME}.init ${FULLNAME}
	cp ChangeLog README ${FULLNAME}
	tar zcvf ${FULLNAME}.tar.gz ${FULLNAME}
	rm -r ${FULLNAME}


