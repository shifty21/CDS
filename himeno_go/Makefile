CFLAGS=build
CC=go
RM=rm -f
EXEC=himeno

all: $(EXEC)

$(EXEC):
	$(CC) $(CFLAGS) -tags netgo -v -o ./bin/$(EXEC)  $(EXEC).go

run:
	./$(EXEC)

clean:
	$(RM) $(EXEC).o $(EXEC)
