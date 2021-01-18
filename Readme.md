# civ clone thing

prototype for a civilization clone. kinda idle game inspired, think kittens game for example

i want games to be short, like playable within a lunch time and for now it's just going to be a PvE single player affair.

## technical stuff

- written in golang
- targetting all contemporary desktop platforms

## setting

you start as a single neolithic village consisting of a few villagers

## win conditions

- launch a supership to colonise distant worlds

## mechanics

- in the "move phase" you place or destroy buildings, make decisions about the economy of the next year, decide research, allocate resources, roles etc
- every turn is one year
- game loop should evaluate negative factors first to minimise cheesing
- on the world map, each tile is a region or city
- each region or city consists of another 8x8 tile grid
	- movement within a region is free
	- before trains, movement between regions takes a turn
	- construction can only happen on neighbouring tiles
	- a settlement must be established in a neighbouring tile in order to move there
	- once trains have been researched, movement between regions is free

### ages

#### neolithic
- research modifier is basically zero in the neolithic age
- there is a chance of a monolith spawning which can be harvested for research but in this age research cannot be allocated

#### protohistory
- available once farming has been discovered
- writing available for research

#### literary age
- available once writing has been discovered. allows for research allocation
- monoliths no longer spawn

#### roman age
- available once ???

#### classical age
- available once ???

#### age of steam
- available once ???

#### atomic age
- available once ???
- there are further ages but game _could_ be "won" at this point
- global annihilation risk is now factored into 

#### computer age
- available once atom bomb and semiconductors have been researched
- researching the atom bomb adds a thermonuclear annihilation modifier to the death chance for all citizens

#### modern age
- available once ???
- representative of present day

#### transhuman age
- available once ???

### construction
- construction takes one year minimum (or can be sped up by assigning more villagers). concurrent projects will affect the construction time

### citizens
- every year a citizen has a chance of dying as a function of 
 	- their age
	- medical science progression
	- amount of doctors 
	- food supply (must be positive) 
	- environmental factors. 
- citizens deaths to be evaluated at turn end and will be evaluated before constructions, etc. not sure of the order of operations yet, but would like to avoid cheesing

#### children
- children cannot be assigned roles until they are 10 (increasing with each epoch)
- births will be calculated based on population demographics, food supply and environmental factors, similar to deaths. this could get very dark. could get quite creative with the death flavour text 
