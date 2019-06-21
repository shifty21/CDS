CC=gcc
OUT=himeno

all: himeno

amd: himeno.c
	$(CC) -o $(OUT) himeno.c -lpthread

debug: himeno.c
	$(CC) -o $(OUT).debug himeno.c -DDEBUG=1 -Wall

clean:
	rm -f $(OUT) $(OUT).debug
