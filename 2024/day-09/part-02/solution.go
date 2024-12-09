package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
)

func parseFile(r io.Reader) ([]string, error) {
	scanner := bufio.NewScanner(r)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

type chunk struct {
	Index, Offset, Size int
}

func parse(line string) ([]chunk, []chunk) {
	files, spaces := []chunk{}, []chunk{}
	totalSize := 0

	for i, char := range line {
		size := parseSize(char)
		if isFile(i) {
			files = append(files, createChunk(i/2, totalSize, size))
		} else {
			spaces = append(spaces, createChunk(0, totalSize, size))
		}
		totalSize += size
	}

	return files, spaces
}

// pars.Size converts a rune to its integer.Size value.
func parseSize(char rune) int {
	return int(char - '0')
}

// isFile determines if the current index corresponds to a file.
func isFile(index int) bool {
	return index%2 == 0
}

// createChunk initializes a chunk with the given parameters.
func createChunk(index, offset, size int) chunk {
	return chunk{
		Index:  index,
		Offset: offset,
		Size:   size,
	}
}

func checksum(files []chunk) int {
	checksum := 0
	for _, file := range files {
		// sum(n, n+1, ..., n+k) = sum(1, ..., n+k) - sum(1, ..., n-1)
		fileEnd := file.Offset + file.Size - 1
		checksum += (fileEnd*(fileEnd+1) - (file.Offset-1)*file.Offset) / 2 * file.Index
	}
	return checksum
}

func solve(input []string) string {
	// Parse input into files and spaces
	files, spaces := parse(input[0])

	// Adjust file offsets based on available space
	adjustFileOffsets(files, spaces)

	// Compute and return the checksum as a string
	return strconv.Itoa(computeChecksum(files))
}

// adjustFileOffsets moves files into available spaces as needed.
func adjustFileOffsets(files, spaces []chunk) {
	for i := len(files) - 1; i > 0; i-- {
		for j := range spaces {
			if canMoveFile(files[i], spaces[j]) {
				moveFile(&files[i], &spaces[j])
			}
		}
	}
}

// canMoveFile checks if a file can fit in a space.
func canMoveFile(file, space chunk) bool {
	return file.Offset >= space.Offset && file.Size <= space.Size
}

// moveFile updates file and space properties after moving a file.
func moveFile(file *chunk, space *chunk) {
	file.Offset = space.Offset
	space.Size -= file.Size
	space.Offset += file.Size
}

// computeChecksum calculates the checksum of the files.
func computeChecksum(files []chunk) int {
	// Placeholder for your checksum calculation logic
	return checksum(files) // Assuming `checksum` is defined elsewhere
}

func main() {
	// Open the file
	file, err := os.Open("input.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// parse the file
	input, err := parseFile(file)

	solution := solve(input)
	fmt.Println(solution)
}
