astgen-cpp: defs.go main.go
	go fmt
	go build

defs.go: defs.go.tmpl
	templator $< > $@
	go fmt $@

test: astgen-cpp
	mkdir -p testdir
	echo '#include "../base.h"' > testdir/ast.h
	./astgen-cpp --header ../../../zefram/spec/ast.def >> testdir/ast.h
	echo '#include "ast.h"' > testdir/ast.cpp
	./astgen-cpp --impl ../../../zefram/spec/ast.def >> testdir/ast.cpp
	clang++ --std=c++11 -c testdir/ast.cpp -o testdir/ast.o; rm -rf testdir