### v0.8.1

 - Explore replaces object with same root label.

### v0.8.0

- add SetDefault to make a Service globally accessing to add new structs.

### v0.7.0

- add theme toggle button (dark,light,auto) (thx to @flaticols) 

### v0.6.2

- fix panic for pointer to slice of embedded type

### v0.6.1

- rename Follow to ExplorePath

### v0.6.0

- add Follow to service for navigation of structs
- add ExploreOption for control of where to put a struct

### v0.5.0

- add "c" (clear) button for not-root removals
- add "⇈" for placing object on the row above
- fixed exploring empty slice and map

### v0.4.2

- fix fixed sized integer keys in maps

### v0.4.0

- (fixed) sorted integer keys in slice/array
- interval fields for large slice/array
- display byte (uint8) values with character
- better cell placement strategy (favor horizontal over vertical)

### < v0.4.0

- see git log