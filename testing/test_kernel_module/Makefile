obj-m := test.o

.PHONY: all
all: clean test.ko.xz

build:
	make -C /lib/modules/$(shell uname -r)/build M=$(PWD) modules

clean:
	rm -f *module-common* *modules.order* *Module.symvers* *test.ko* *test.mod* *test.o*

test.ko.xz: build
	xz -f test.ko
