include .env
# -ldflags 编译选项，-s -w 去掉调试信息，可以减小构建后文件体积
GOBUILD=go build -ldflags '-w -s'
BINARY=gtoo_$(Version)
BINDIR=bin

all: linux macos win64

linux:
	GOARCH=amd64 GOOS=linux $(GOBUILD) -o $(BINDIR)/$(BINARY)_$@

macos:
	GOARCH=amd64 GOOS=darwin $(GOBUILD) -o $(BINDIR)/$(BINARY)_$@

win64:
	GOARCH=amd64 GOOS=windows $(GOBUILD) -o $(BINDIR)/$(BINARY)_$@.exe

clean:
	@if [ -d $(BINDIR) ] ; then rm $(BINDIR)/* ; fi

help:
	@echo "make - 格式化 Go 代码, 并编译生成二进制文件"
	@echo "make [linux|macos|win64] - 编译相应平台的 Go 代码, 生成二进制文件"
	@echo "make clean - 移除二进制文件和 vim swap files"
