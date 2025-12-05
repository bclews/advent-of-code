# Day 5: Cafeteria

## Part One

After the forklifts break through the wall, you discover a cafeteria kitchen where the Elves are struggling with their new inventory management system—they can't determine which ingredients are fresh versus spoiled. The system database consists of two sections separated by a blank line: a list of fresh ingredient ID ranges (like 3-5 or 10-14, where ranges are inclusive and can overlap), followed by a list of available ingredient IDs to check. An ingredient ID is considered fresh if it falls within any of the given ranges. In the example with ranges 3-5, 10-14, 16-20, and 12-18 and available IDs 1, 5, 8, 11, 17, and 32, three IDs are fresh: 5 (in range 3-5), 11 (in range 10-14), and 17 (in both ranges 12-18 and 16-20). Your task is to process the database and count how many of the available ingredient IDs are fresh.

## Part Two

As the Elves start bringing spoiled inventory to the trash chute, they ask for help determining all possible fresh ingredient IDs so they can stop bothering you with future inventory checks. Now the second section of the database (the available ingredient IDs) is irrelevant—instead, you must count every unique ingredient ID that the fresh ranges collectively consider to be fresh, accounting for overlapping ranges. Using the same example ranges 3-5, 10-14, 16-20, and 12-18, these ranges cover the IDs 3, 4, 5, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, and 20, which totals 14 unique ingredient IDs. The key insight is to merge overlapping ranges first—sorting them and combining adjacent or overlapping sections—then calculate the total coverage arithmetically rather than expanding every range into individual IDs, ensuring the solution remains efficient even for very large ranges.

More details can be found [here](https://adventofcode.com/2025/day/5).
