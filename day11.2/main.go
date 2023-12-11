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
const scale = 1000000

type coord struct {
  x, y int
}

func inflate(g []coord, maxx, maxy int) ([]coord, int, int) {
  ret := []coord{}
  addX := []int{}
  addY := []int{}
  for x := 0; x < maxx; x++ {
    found := false
    for _, v := range g {
      if v.x == x {
	found = true
	break
      }
    }
    if !found {
      addX = append(addX, x)
    }
  }
  for y := 0; y < maxy; y++ {
    found := false
    for _, v := range g {
      if v.y == y {
	found = true
	break
      }
    }
    if !found {
      addY = append(addY, y)
    }
  }

  tickX := map[coord]int{}
  tickY := map[coord]int{}

  for _, v := range addX {
    for i := range g {
      if g[i].x > v {
	tickX[g[i]]++
      }
    }
  }
  for _, v := range addY {
    for i := range g {
      if g[i].y > v {
	tickY[g[i]]++
      }
    }
  }

  if debug {
    fmt.Printf("%v\n%v\n", tickX, tickY)
  }

  for i := range g {
    incX := tickX[g[i]]
    incY := tickY[g[i]]
    g[i].x += (scale * incX)
    if incX > 0 {
      g[i].x -= incX
    }
    g[i].y += (scale * incY)
    if incY > 0 {
      g[i].y -= incY
    }
  }

  if debug {
    fmt.Printf("%v\n%v\n", addX, addY)
  }

  for _, v := range g {
    ret = append(ret, v)
  }
  return ret, maxx+(scale*len(addX) - 1), maxy+(scale*len(addY) - 1)
}

func printGalaxies(g []coord, maxx, maxy int) {
  if !debug {
    return
  }
  tmp := map[coord]bool{}
  for _, v := range g {
    tmp[v] = true
  }

  for y := 0; y < maxy; y++ {
    for x := 0; x < maxx; x++ {
      if tmp[coord{x, y}] {
	fmt.Print("#")
      } else {
	fmt.Print(".")
      }
    }
    fmt.Println()
  }
}

func abs(x int) int {
  if x > 0 {
    return x
  }
  return -x
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

  galaxies := []coord{}
  y := 0
  maxX := 0
  for scan.Scan() {
    x := 0
    for _, c := range scan.Text() {
      maxX = len(scan.Text())
      if string(c) == "#" {
	galaxies = append(galaxies, coord{x, y})
      }
      x++
    }
    y++
  }

  printGalaxies(galaxies, maxX, y)
  galaxies, maxX, y = inflate(galaxies, maxX, y)
  printGalaxies(galaxies, maxX, y)

  sum := 0
  for i := 0; i < len(galaxies); i++ {
    for j := i + 1; j < len(galaxies); j++ {
      dist := abs(galaxies[i].x - galaxies[j].x) + abs(galaxies[i].y - galaxies[j].y)
      sum += dist
      if debug {
	fmt.Printf("%v-%v: %v\n", i+1, j+1, dist)
      }
    }
  }

  fmt.Printf("%v\n", sum)
}
