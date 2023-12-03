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

func isSpecial(s string) bool {
  if unicode.IsNumber(rune(s[0])) || s == "." {
    return false
  }
  return true
}

func hasSpecialNeighbor(engineMap [][]string, row, col int) bool {
  if row > 0 {
    if col > 0 {
      if isSpecial(engineMap[row-1][col-1]) {
	return true 
      }
    }
    if isSpecial(engineMap[row-1][col]) {
      return true
    }
    if col < len(engineMap[row-1]) - 1 {
      if isSpecial(engineMap[row-1][col+1]) {
	return true
      }
    }
  }
  if col > 0 {
    if isSpecial(engineMap[row][col-1]) {
      return true 
    }
  }
  if isSpecial(engineMap[row][col]) {
    return true
  }
  if col < len(engineMap[row]) - 1 {
    if isSpecial(engineMap[row][col+1]) {
      return true
    }
  }
  if row < len(engineMap) - 1 {
    if col > 0 {
      if isSpecial(engineMap[row+1][col-1]) {
	return true
      }
    }
    if isSpecial(engineMap[row+1][col]) {
      return true
    }
    if col < len(engineMap[row+1]) - 1 {
      if isSpecial(engineMap[row+1][col+1]) {
	return true
      }
    }
  }
  return false
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

  i, j := 0, 0
  partNums := []string{}
  nonPartNums := []string{}
  for i < len(engineMap) {
    buf := []string{}
    part := false
    for j < len(engineMap[i]) {
      if unicode.IsNumber(rune(engineMap[i][j][0])) {
	buf = append(buf, engineMap[i][j])
	part = part || hasSpecialNeighbor(engineMap, i, j)
	j++
	continue
      }
      if len(buf) != 0 {
	if part {
	  partNums = append(partNums, strings.Join(buf, ""))
	} else {
	  nonPartNums = append(nonPartNums, strings.Join(buf, ""))
	}
	buf = []string{}
	part = false
      }
      j++
    }
    if len(buf) != 0 {
      if part {
	partNums = append(partNums, strings.Join(buf, ""))
      } else {
	nonPartNums = append(nonPartNums, strings.Join(buf, ""))
      }
    }
    i++
    j = 0
  }

  fmt.Printf("%v\n%v\n", partNums, nonPartNums)

  sum := 0
  for _, v := range partNums {
    val, _ := strconv.Atoi(v)
    sum += val
  }

  fmt.Printf("%v\n", sum)
}
