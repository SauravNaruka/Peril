# Peril

## Client Commands
### spawn
The spawn command allows a player to add a new unit to the map under their control. 

Possible unit types are:
- infantry
- cavalry
- artillery

Possible locations are:
- americas
- europe
- africa
- asia
- antarctica
- australia

Example usage:
```
spawn europe infantry
```

After spawning a unit, you should see its ID printed to the console.

### move
The move command allows a player to move their units to a new location. It accepts two arguments: the destination, and the ID of the unit.

Example usage:
```
move europe 1
```

### status
Print the current status of the player's game state.

### help
Print a list of available commands.

### spam
Spamming not allowed yet!

### quit
Quit exit the REPL.
