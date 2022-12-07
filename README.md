## Simple CLI version of proxx game 


### Run the game
```go run main.go --size 10 --black-holes 5 ```

#### CLI params
```black-holes``` count of black holes (default 10)

```size``` board size (default 8)

```emulate``` emulate the game. PC will click random cells itself


### About the project
I wasn't trying to build a perfect production app. I intended to create a working example in the shortest possible time.
That's why I wrote tests only for the most critical methods of the game's logic.
The current UI is just to represent the work of business logic (game.go) and should be replaced with some TUI library.