## Battleship

This is an entry to the [Gophercon Coding Challenge](https://github.com/gojek-engineering/gophercon-2017).

### Building and running

```sh
$ https://github.com/prashanthpai/battleship
$ go build battleship.go
$ ./battleship input.txt output.txt
```

The result of the game run is printed to stdout and also to the output file.

Sample run:

```sh
[ppai@gd2-1 battleship]$ ./battleship input.txt output.txt
Player1
O O _ _ _ 
_ X _ _ _ 
B _ _ X _ 
_ _ _ _ B 
_ _ _ X _ 

Player2
_ X _ _ _ 
_ _ _ _ _ 
_ _ _ X _ 
B O _ _ B 
_ X _ O _ 

P1:3
P2:3
It is a draw
```
