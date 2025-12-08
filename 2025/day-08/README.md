# Day 8: Playground

## Part One

The Elves are setting up Christmas decorations by connecting suspended junction boxes with strings of lights in 3D space. To save on lights, they want to connect the closest pairs of junction boxes first, using straight-line Euclidean distance. Given positions of junction boxes as X,Y,Z coordinates, connect the 1000 pairs that are closest together (some connections may be skipped if they would create redundant paths). After making these connections, find the sizes of all resulting circuits and multiply together the sizes of the three largest circuits.

## Part Two

Continue connecting junction boxes until all boxes form a single circuit, effectively building a complete minimum spanning tree. The task is to identify the final edge that unifies all junction boxes into one component, then multiply together the X coordinates of the two junction boxes involved in that last connection.

More details can be found [here](https://adventofcode.com/2025/day/8).
