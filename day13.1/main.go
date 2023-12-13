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
var debug = false

func min(a, b int) int {
  if a < b {
    return a
  }
  return b
}

func isPalindrom(t []uint32) int {
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
      return i+1
    }
  }
  return  -1
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
    if debug {
      fmt.Printf("%v\n%v\n", rows, cols)
    }

    vertical := isPalindrom(cols)
    horizontal := isPalindrom(rows)
    if horizontal >= 0 {
      sum += 100 * horizontal
    }
    if vertical >= 0 {
      sum += vertical
    }

    image = [][]string{}
  }

  fmt.Printf("%v\n", count)
  fmt.Printf("%v\n", sum)
}
