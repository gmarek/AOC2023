package main

import (
  "bufio"
  "flag"
  "fmt"
//  "math"
  "os"
  "sort"
//  "strconv"
//  "strings"
//  "unicode"
)

var prod = flag.Bool("real", false, "run on real input")
var debug = true

const e = 0
const s = 1
const w = 2
const n = 3

type dig struct {
  dir int
  l int
}

func max(a, b int) int {
  if a > b {
    return a
  }
  return b
}

func min(a, b int) int {
  if a < b {
    return a
  }
  return b
}

type point struct {
  x, y int
}

func abs(a int) int {
  if a < 0 {
    return -a
  }
  return a
}

func cutRectangle(s *[]point) int {
  fmt.Printf("in %v: %v\n", len(*s), *s)
  if len(*s) == 4 {
    ret := abs((*s)[0].x - (*s)[1].x) * abs((*s)[0].y - (*s)[3].y)
    *s = []point{}
    return ret
  }
  minXInd := -1
  minX := int(1e10)
  for i, p := range *s {
    if p.x < minX {
      minX = p.x
      minXInd = i
    }
  }
  secondMinInd := (minXInd+1) % len(*s)
  next1, next2 := 0, 0
  if (*s)[secondMinInd].x == minX {
    next1 = minXInd - 1
    if next1 < 0 {
      next1 = len(*s) - 1
    }
    next2 = (secondMinInd + 1) % len(*s)
  } else {
    secondMinInd = minXInd - 1
    if secondMinInd < 0 {
      secondMinInd = len(*s) - 1
    }
    next1 = (minXInd + 1) % len(*s)
    next2 = secondMinInd - 1
    if next2 < 0 {
      next2 = len(*s) - 1
    }
  }
  if debug {
    fmt.Printf("%v %v %v %v\n", minXInd, secondMinInd, next1, next2)
    fmt.Printf("%v %v %v %v\n", (*s)[minXInd], (*s)[secondMinInd], (*s)[next1], (*s)[next2])
  }

  r1 := abs((*s)[minXInd].y - (*s)[secondMinInd].y)
  r2 := abs(minX - min((*s)[next1].x, (*s)[next2].x))
  fmt.Printf("%v * %v = %v\n", r1, r2, r1 * r2)
  ret := r1 * r2
  indices := []int{minXInd, secondMinInd, next1, next2}
  sort.Ints(indices)
  if (*s)[next1].x == (*s)[next2].x {
    for i := 3; i >= 0; i-- {
      fmt.Printf("removing: %v\n", (*s)[i])
      remove(s, i)
    }
  } else {
    if indices[0] == 0 && indices[3] == len(*s) - 1 {
      //needs cleanup
      for indices[0] == 0 {
	*s = append((*s)[1:], (*s)[0])
	indices = append(indices[1:], indices[0])
	for i := range indices {
	  indices[i]--
	  if indices[i] < 0 {
	    indices[i] = len(*s) - 1
	  }
	}
      }
    }
    newPoint := point{}
    newPoint.x = min((*s)[indices[0]].x, (*s)[indices[3]].x)
    if (*s)[indices[0]].x < (*s)[indices[3]].x {
      newPoint.y = (*s)[indices[2]].y
    } else {
      newPoint.y = (*s)[indices[1]].y
    }

    fmt.Printf("adding: %v\n", newPoint)
    fmt.Printf("deleting: %v %v\n", (*s)[indices[2]], (*s)[indices[1]])
    remove(s, indices[2])
    remove(s, indices[1])

    suffix := (*s)[indices[3]-2:]
    *s = append((*s)[:indices[0]], newPoint)
    *s = append(*s, suffix...)
  }
  return ret
}

func remove(s *[]point, i int) {
  fmt.Printf("remove: %v from %v\n", i, *s)
  if i == 0 {
    *s = (*s)[1:]
  } else if i == len(*s) - 1 {
    *s = (*s)[:len(*s)-1]
  } else {
    *s = append((*s)[:i], (*s)[i+1:]...)
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

  digs := []dig{}
  for scan.Scan() {
    d := dig{}
    sTrash := ""
    iTrash := 0
    dString := ""
    fmt.Sscanf(scan.Text(), "%s %d %s", &sTrash, &iTrash, &dString)
    lString := dString[2:len(dString) - 2]
    fmt.Sscanf(lString, "%x", &d.l)
    dirStr := dString[len(dString)-2:]
    fmt.Sscanf(dirStr, "%d", &d.dir)
    digs = append(digs, d)
  }

  shape := []point{{0, 0}}
  for _, d := range digs {
    next := point{shape[len(shape) - 1].x, shape[len(shape) - 1].y}
    switch d.dir {
    case n: {
      next.y -= d.l
    }
    case e: {
      next.x += d.l
    }
    case s: {
      next.y += d.l
    }
    case w: {
      next.x -= d.l
    }
    }
    shape = append(shape, next)
  }

  shape = shape[:len(shape)-1]

  fmt.Printf("%v\n", shape)

  sum := 0
  for len(shape) > 0 {
    sum += cutRectangle(&shape)
  }

  fmt.Printf("%v\n", sum)
}
