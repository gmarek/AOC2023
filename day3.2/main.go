package main

import (
  "bufio"
  "flag"
  "fmt"
  "os"
  "strconv"
  "strings"
  "unicode"
)

var prod = flag.Bool("real", false, "run on real input")

type point struct {
  x, y int
}

func processGear(engineMap [][]string, row, col int) int64 {
  checked := map[point]bool{}
  if row > 0 {
    if col > 0 {
      checked[point{row-1, col-1}] = false
    }
    checked[point{row-1, col}] = false
    if col < len(engineMap[row-1]) {
      checked[point{row-1, col+1}] = false
    }
  }
  if col > 0 {
    checked[point{row, col-1}] = false
  }
  if col < len(engineMap[row]) {
    checked[point{row, col+1}] = false
  }
  if row < len(engineMap) {
    if col > 0 {
      checked[point{row+1, col-1}] = false
    }
    checked[point{row+1, col}] = false
    if col < len(engineMap[row+1]) {
      checked[point{row+1, col+1}] = false
    }
  }

  nums := []int64{}
  for k, v := range checked {
    if v || !unicode.IsNumber(rune(engineMap[k.x][k.y][0])) {
      continue
    }
    buf := []string{engineMap[k.x][k.y]}
    next := k.y+1
    for next < len(engineMap[k.x]) {
      checked[point{k.x, next}] = true
      if unicode.IsNumber(rune(engineMap[k.x][next][0])) {
	buf = append(buf, engineMap[k.x][next])
	next++
      } else {
	break
      }
    }
    prev := k.y-1
    for prev >= 0 {
      checked[point{k.x, prev}] = true
      if unicode.IsNumber(rune(engineMap[k.x][prev][0])) {
	buf = append([]string{engineMap[k.x][prev]}, buf...)
	prev--
      } else {
	break
      }
    }
    if len(buf) > 0 {
      val, _ := strconv.Atoi(strings.Join(buf, ""))
      nums = append(nums, int64(val))
    }
  }

  if len(nums) == 2 {
    return nums[0] * nums[1]
  }
  return 0
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

  engineMap := [][]string{}
  for scan.Scan() {
    row := []string{}
    for _, c := range scan.Text() {
      row = append(row, string(c))
    }
    engineMap = append(engineMap, row)
  }

  sum := int64(0)
  for i := 0; i < len(engineMap); i++ {
    for j := 0; j < len(engineMap[i]); j++ {
      if engineMap[i][j] == "*" {
	sum += processGear(engineMap, i, j)
      }
    }
  }

  fmt.Printf("%v\n", sum)
}
