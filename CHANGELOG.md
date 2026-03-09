#CHANGELOG

## Disjoint algorithms implementation 
**File:** `algo.go`
  - **What:** Added `Disjoints` — added Disjoints functions.
  - **Why:** The ants used to only get their way by gettign the fastest path they add now they check each others itinerary and not create jams.

## Ant path handling
**File:** `fourmis.go`
  - **What:** Fixed `SimulateAntsonPath` — changed to SimulateAntsonPaths.
  - **Why:** The original SimulateAntsonPath dinamically calculated the Ants route but the clashes happened between ants when they could both move into a room.


## RoomOccupancy
**File: ** `fourmis.go` and `structs.go`

- **What:** Added room Occupancy slice
- **Why:** permitted the ants to check if they can wait for a room 

