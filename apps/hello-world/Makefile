# hello-world Makefile
CC ?= clang
CFLAGS ?= -Wall -Wextra -static

hello-world: hello-world.c
	$(CC) $(CFLAGS) -o $@ $<

clean:
	rm -f hello-world

.PHONY: clean
