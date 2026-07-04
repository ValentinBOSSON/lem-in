# Lem-in

A digital ant farm, written in Go. The program reads a colony map (chambers and tunnels), finds several shortest paths that share no chamber, and distributes the ants across them so the whole colony reaches the exit in as few turns as possible.

Built during my training at Zone01 Rouen.

## How it works

1. **Parsing**: the map file declares the number of ants, the chambers (`name x y`, with `##start` and `##end` markers) and the tunnels (`a-b`). Malformed input is rejected with `ERROR: invalid data format`.
2. **Pathfinding**: breadth-first searches under constraints find all shortest paths, then a vertex-disjoint set is selected so no two paths share a chamber. Priority queues come from Go's `container/heap`.
3. **Distribution**: ants are split across the selected paths to minimize the number of turns (a longer path can still be worth using when the short one is crowded).
4. **Simulation**: the run is simulated turn by turn and written to `deplacements.log`.

## Run

```
go run . test0.txt
```

The repository includes six example maps, `test0.txt` through `test5.txt`.

## Project structure

```
main.go             entry point, output
initialisation.go   map parsing
algo.go             pathfinding: disjoint paths, BFS, heaps
fourmis.go          ant creation and turn-by-turn simulation
structs.go          chambers, ants, graph structures
```
