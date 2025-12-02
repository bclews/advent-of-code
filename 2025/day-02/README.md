# Day 2: Gift Shop

## Part One

Inside the North Pole base, you pass through the gift shop where a clerk asks for help fixing a database problem caused by a young Elf who entered many invalid product IDs. Only a few ID ranges (your puzzle input) still need checking. Each range lists a start and end number, and an ID is considered invalid if it consists of a sequence of digits repeated twice (like 55, 6464, or 123123), with no leading zeroes allowed. Your task is to identify all such duplicated-pattern IDs found within the specified ranges. In the provided example, several ranges contain one or more invalid IDs—such as 11 and 22 in 11–22, 99 in 95–115, and 1010 in 998–1012—while others contain none, and the total sum of all invalid IDs found is 1,227,775,554.

## Part Two

The clerk realises more invalid IDs remain, suggesting the young Elf used additional repetitive patterns. Now an ID is invalid if it consists of any sequence of digits repeated two or more times, such as 12341234, 123123123, 1212121212, or 1111111. Checking the same ranges again reveals extra invalid IDs—like 111 in 95–115, 999 in 998–1012, 565656 in 565653–565659, and several others—while some ranges remain unchanged. When all newly identified invalid IDs are added together with those found previously, the total for the example reaches 4,174,379,265.

More details can be found [here](https://adventofcode.com/2025/day/2).
