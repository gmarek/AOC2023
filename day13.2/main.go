package main

import (
  "bufio"
  "flag"
  "fmt"
  "hash/crc32"
//  "math"
  "os"
//  "sort"
//  "strconv"
  "strings"
//  "unicode"
)

var prod = flag.Bool("real", false, "run on real input")
var debug = true

func min(a, b int) int {
  if a < b {
    return a
  }
  return b
}

func fingerprints(image [][]string) ([]uint32, []uint32) {
  rows := []uint32{}
  for _, r := range image {
    rows = append(rows, crc32.ChecksumIEEE([]byte(strings.Join(r, ""))))
  }

  cols := []uint32{}
  for x := 0; x < len(image[0]); x++ {
    col := []string{}
    for y := 0; y < len(image); y++ {
      col = append(col, image[y][x])
    }
    cols = append(cols, crc32.ChecksumIEEE([]byte(strings.Join(col, ""))))
  }

  return rows, cols
}

func isPalindrom(t []uint32) []int {
  ret := []int{}
  for i := 0; i < len(t) - 1; i++ {
    pal := true
    for j := 0; j < len(t); j++ {
      if i - j < 0 || i + 1 + j >= len(t) {
	break
      }
      if t[i-j] != t[i + 1 + j] {
	pal = false
	break
      }
    }
    if pal {
      ret = append(ret, i+1)
    }
  }
  return ret
}

func main() {
  flag.Parse()

  filename := "test"
  if *prod {
    filename = "input"
  }
  f, err := os.Open(filename)
  if err != nil {
    panic(err)
  }

  scan := bufio.NewScanner(bufio.NewReader(f))
  image := [][]string{}
  sum := 0
  count := 0
  for scan.Scan() {
    if scan.Text() != "" {
      row := []string{}
      for _, c := range scan.Text() {
	row = append(row, string(c))
      }
      image = append(image, row)
      continue
    }
    count++
    fmt.Printf("-----%v-----\n", count)
    rows, cols := fingerprints(image)

    vertical := isPalindrom(cols)
    horizontal := isPalindrom(rows)

    if debug {
      fmt.Printf("%v, %v\n", vertical, horizontal)
    }

    knownV := map[int]bool{}
    if len(vertical) > 0 {
      knownV[vertical[0]] = true
    }
    knownH := map[int]bool{}
    if len(horizontal) > 0 {
      knownH[horizontal[0]] = true
    }
    for y := 0; y < len(image); y++ {
      for x := 0; x < len(image[0]); x++ {
	c := image[y][x]
	if c == "#" {
	  image[y][x] = "."
	} else {
	  image[y][x] = "#"
	}
	newR, newC := fingerprints(image)
	newV := isPalindrom(newC)
	for _, v := range newV {
	  if !knownV[v] {
	    sum += v
	    knownV[v] = true
	    if debug {
	      fmt.Printf("vert: %v -> %v\n", vertical, newV)
	    }
	  }
	}
	newH := isPalindrom(newR)
	for _, v := range newH {
	  if !knownH[v] {
	    sum += 100*v
	    knownH[v] = true
	    if debug {
	      fmt.Printf("hor: %v -> %v\n", horizontal, newH)
	    }
	  }
	}
	image[y][x] = c
      }
    }
    image = [][]string{}
  }

  fmt.Printf("%v\n", sum)
}
