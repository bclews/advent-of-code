# Day 11: Reactor

## Problem Summary

A factory reactor has communication issues with a new server rack due to problematic data paths through connected devices. Given a directed graph of devices and their output connections, Part One requires finding all distinct paths from device "you" to the main output "out" (answer: 788 paths). Part Two narrows the search to paths from "svr" to "out" that must visit both specific devices "dac" and "fft" in any order, identifying the problematic data routes (answer: 316,291,887,968,000 paths).

## Solution Approach

Part One uses a simple recursive depth-first search (DFS) to count all paths from start to target. Part Two extends this with boolean flags to track whether required nodes ("dac" and "fft") have been visited along each path, only counting paths that reach the target with both flags set. To handle the massive search space efficiently, the solution employs memoization by caching results for each unique state `(node, visited_dac, visited_fft)`, reducing time complexity from exponential to O(V) where V is the number of nodes.

More details can be found [here](https://adventofcode.com/2025/day/11).
