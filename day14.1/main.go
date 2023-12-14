package main

import (
  "bufio"
  "flag"
  "fmt"
//  "math"
  "os"
//  "sort"
//  "strconv"
//  "strings"
//  "unicode"
)

var prod = flag.Bool("real", false, "run on real input")
var debug = false

type field int

const empty field = 0
const static field = 1
const stone field = 2

func printDish(d [][]field) {
  if !debug {
    return
  }
  for y := 0; y < len(d); y++ {
    for x := 0; x < len(d[0]); x++ {
      if d[y][x] == empty {
	fmt.Print(".")
      }
      if d[y][x] == static {
	fmt.Print("#")
      }
      if d[y][x] == stone {
	fmt.Print("O")
      }
    }
    fmt.Println()
  }
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

  dish := [][]field{}
  for scan.Scan() {
    row := []field{}
    for _, c := range scan.Text() {
      if string(c) == "." {
	row = append(row, empty)
      } else if string(c) == "#" {
	row = append(row, static)
      } else if string(c) == "O" {
	row = append(row, stone)
      }
    }
    dish = append(dish, row)
  }

  for y := 0; y < len(dish); y++ {
    for x := 0; x < len(dish[0]); x++ {
      if dish[y][x] != stone {
	continue
      }
      for z := y-1; z >= 0; z-- {
	if dish[z][x] == empty {
	  dish[z][x] = stone
	  dish[z+1][x] = empty
	} else {
	  break
	}
      }
    }
  }

  sum := 0
  for y := 0; y < len(dish); y++ {
    for x := 0; x < len(dish[0]); x++ {
      if dish[y][x] == stone {
	sum += len(dish) - y
      }
    }
  }

  fmt.Printf("%v\n", sum)
  printDish(dish)
}
