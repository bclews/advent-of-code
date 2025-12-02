# Day 1: Secret Entrance

## Part One

On the first day, the Elves reveal that although they’ve finally adopted project management and avoided their usual Christmas crisis, they now face a new one: none of them have time left to decorate the North Pole. You arrive at a secret entrance to help, but the password has been changed and now lies inside a safe that opens by following a list of dial rotations. The dial, numbered 0–99, starts at 50 and moves left or right based on the instructions; because it wraps around, reaching 0 can happen from either direction. However, your security training reveals that the safe is a decoy—the real password is simply the number of times the dial ends up pointing at 0 after any rotation. The puzzle for Day 1 is to process the given sequence of rotations and count those zero landings to discover the true password.

## Part Two

In part two, you discover that the door still won’t open because an updated security method applies: instead of counting only the times the dial ends a rotation at 0, you must count every single click that lands on 0 while the dial is turning. This means intermediate positions matter—long rotations can pass over 0 many times. Using the same list of rotations as before, the example shows that the dial hits 0 three times at the ends of rotations and three more times during the movements, giving a total of six. The new task is therefore to process all rotations and count every click that results in the dial pointing at 0, including wrap-arounds, to determine the true password.

More details can be found [here](https://adventofcode.com/2025/day/1).
