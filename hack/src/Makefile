#!/usr/bin/env bash

CC = gcc

CFLAGS = 
LOCAL_CFLAGS = 
DEFS = 
LOCAL_DEFS = 

CCFLAGS = $(DEFS) $(LOCAL_DEFS) $(LOCAL_CFLAGS) $(CFLAGS)

SHOBJ_CC = gcc
SHOBJ_CFLAGS = -fno-common
SHOBJ_LD = ${CC}
SHOBJ_LDFLAGS = -dynamiclib -dynamic -undefined dynamic_lookup 
SHOBJ_XLDFLAGS = 
SHOBJ_LIBS = 
SHOBJ_STATUS = supportedc

.c.o:
	$(SHOBJ_CC) $(SHOBJ_CFLAGS) $(CCFLAGS) -c -o $@ $<

../libexec/devctlenv.dylib: realpath.o
	$(SHOBJ_LD) $(SHOBJ_LDFLAGS) $(SHOBJ_XLDFLAGS) -o $@ realpath.o $(SHOBJ_LIBS)

clean:
	rm -f *.o ../libexec/*.dylib

