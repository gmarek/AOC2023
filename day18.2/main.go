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

const left = -1
const right = 1

type dig struct {
  dir int
  l int
}

func reflect(s *[]point, start, end int) {
  for start < end {
    tmp := (*s)[start]
    (*s)[start] = (*s)[end]
    (*s)[end] = tmp
    start++
    end--
  }
}

func makeFirst(s *[]point, i int) {
  if i == 0 {
    return
  }
  reflect(s, 0, len(*s) - 1)
  reflect(s, 0, len(*s) - i - 1)
  reflect(s, len(*s) - i, len(*s) - 1)
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

func nextInd(i int, s *[]point) int {
  return (i+1) % len(*s)
}

func prevInd(i int, s *[]point) int {
  if i == 0 {
    return len(*s)-1
  }
  return i-1
}

func neigh(i int, s *[]point) []int {
  return []int{prevInd(i, s), nextInd(i, s)}
}

func cutRectangle(s []point) (int, [][]point) {
  if debug {
    fmt.Println("-----------------")
    fmt.Printf("in %v: %v\n", len(s), s)
  }
  minXInd := map[int]bool{}
  minX := int(1e10)
  nearMinX := int(1e10)
  for i, p := range *s {
    if p.x < minX {
      if minX != int(1e10) {
	nearMinX = minX
      }
      minX = p.x
      minXInd = map[int]bool{i: true}
    } else if p.x == minX {
      minXInd[i] = true
    }
    if p.x < nearMinX && p.x > minX {
      nearMinX = p.x
    }
  }
  if debug {
    fmt.Printf("min: %v, second: %v\n", minX, nearMinX)
  }

  if len(*s) == 4 {
    r1 := abs(nearMinX - minX)
    r2 := 0
    for i := range *s {
      val := abs((*s)[i].y - (*s)[(i+1)%len(*s)].y)
      if val > r2 {
	r2 = val
      }
    }
    if debug {
      fmt.Printf("last: %v * %v = %v\n", r1, r2, r1*r2)
    }
    ret := r1 * r2
    *s = []point{}
    return ret
  }
  remMap := map[int]bool{}
  if debug {
    fmt.Printf("%v\n", minXInd)
  }
  ret := 0
  for i := range minXInd {
    ns := neigh(i, s)
    if debug {
      fmt.Printf("========\n")
      fmt.Printf("(%v) %v -> (%v) %v (%v) %v\n", i, (*s)[i], ns[0], (*s)[ns[0]], ns[1], (*s)[ns[1]])
    }
    if minXInd[ns[0]] {
      r1 := nearMinX - minX
      r2 := abs((*s)[i].y - (*s)[ns[0]].y)
      fmt.Printf("0) %v * %v = %v\n", r1, r2, r1*r2)
      if debug {
	fmt.Printf("%v %v\n", (*s)[ns[0]], (*s)[i])
      }
      ret += r1 * r2
    }
    if minXInd[ns[1]] {
      r1 := nearMinX - minX
      r2 := abs((*s)[i].y - (*s)[ns[1]].y)
      fmt.Printf("1) %v * %v = %v\n", r1, r2, r1*r2)
      if debug {
	fmt.Printf("%v %v\n", (*s)[ns[1]], (*s)[i])
      }
      ret += r1 * r2
    }
  }

  for i := range minXInd {
    (*s)[i].x = nearMinX
  }
  if debug {
    fmt.Printf("post squash: %v\n", *s)
  }

  for i := range minXInd {
    ns := neigh(i, s)
    if (*s)[i].x == (*s)[ns[0]].x && (*s)[i].y == (*s)[ns[0]].y {
      if debug {
	fmt.Printf("removing duplicate %v: %v\n", i, (*s)[i])
      }
      remMap[i] = true
    }
    if (*s)[i].x == (*s)[ns[1]].x && (*s)[i].y == (*s)[ns[1]].y {
      if debug {
	fmt.Printf("removing duplicate %v: %v\n", i, (*s)[i])
      }
      remMap[i] = true
    }
  }

  remList := []int{}
  for k := range remMap {
    remList = append(remList, k)
  }
  sort.Ints(remList)

  for i := len(remList) - 1; i >= 0; i-- {
    remove(s, remList[i])
  }

  if debug {
    fmt.Printf("no duplicates: %v\n", *s)
  }

  slice := map[int]bool{}
  for i, v := range *s {
    if v.x == nearMinX {
      slice[i] = true
    }
  }

  remMap = map[int]bool{}
  remList = []int{}
  for k := range slice {
    ns := neigh(k, s)
    same := true
    for _, n := range ns {
      if (*s)[k].x != (*s)[n].x {
	same = false
	break
      }
    }
    if same {
      remMap[k] = true
    }
  }
  for k := range remMap {
    remList = append(remList, k)
  }
  sort.Ints(remList)

  for i := len(remList) - 1; i >= 0; i-- {
    remove(s, remList[i])
  }

  if debug {
    fmt.Printf("return: %v\n", ret/2)
  }
  return ret / 2
}

func remove(s *[]point, i int) {
  fmt.Printf("remove: %v (%v) from %v\n", i, (*s)[i], *s)
  if i == 0 {
    *s = (*s)[1:]
  } else if i == len(*s) - 1 {
    *s = (*s)[:len(*s)-1]
  } else {
    *s = append((*s)[:i], (*s)[i+1:]...)
  }
}

func insertTurn(turns *[]int, t int, global int) int {
  if len(*turns) < 2 {
    *turns = append(*turns, t)
    return 1
  }
  (*turns)[0] = (*turns)[1]
  (*turns)[1] = t
  if (*turns)[1] != (*turns)[0] {
    return 0
  }
  if (*turns)[1] == global {
    return 1
  }
  return -1
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
    dString := ""
    fmt.Sscanf(scan.Text(), "%s %d", &dString, &d.l)
    switch dString {
    case "R": d.dir = e
    case "D": d.dir = s
    case "L": d.dir = w
    case "U": d.dir = n
    }
    /*sTrash := ""
    iTrash := 0
    fmt.Sscanf(scan.Text(), "%s %d %s", &sTrash, &iTrash, &dString)
    lString := dString[2:len(dString) - 2]
    fmt.Sscanf(lString, "%x", &d.l)
    dirStr := dString[len(dString)-2:]
    fmt.Sscanf(dirStr, "%d", &d.dir)
    */
    digs = append(digs, d)
  }

  shape := []point{{0, 0}}
  turns := []int{}
  prev := -1
  globalTurn := -1
  for _, d := range digs {
    if len(turns) == 1 {
      globalTurn = turns[0]
    }
    offset := 0
    if prev == -1 {
      prev = d.dir
    } else {
      switch d.dir {
      case n: {
	if prev == e {
	  offset = insertTurn(&turns, left, globalTurn)
	  shape[len(shape)-1].x += offset
	} else {
	  offset = insertTurn(&turns, right, globalTurn)
	  shape[len(shape)-1].x -= offset
	}
      }
      case e: {
	if prev == n {
	  offset = insertTurn(&turns, right, globalTurn)
	  shape[len(shape)-1].y -= offset
	} else {
	  offset = insertTurn(&turns, left, globalTurn)
	  shape[len(shape)-1].y += offset
	}
      }
      case s: {
	if prev == e {
	  offset = insertTurn(&turns, right, globalTurn)
	  shape[len(shape)-1].x += offset
	} else {
	  offset = insertTurn(&turns, left, globalTurn)
	  shape[len(shape)-1].x -= offset
	}
      }
      case w: {
	if prev == n {
	  offset = insertTurn(&turns, left, globalTurn)
	  shape[len(shape)-1].y -= offset
	} else {
	  offset = insertTurn(&turns, right, globalTurn)
	  shape[len(shape)-1].y += offset
	}
      }
      }
      prev = d.dir
    }
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
    fmt.Printf("insert %v\n", next)
    shape = append(shape, next)
  }

  shape = shape[:len(shape)-1]
  sum := 0
  for len(shape) > 0 {
    sum += cutRectangle(&shape)
  }

  fmt.Printf("%v\n", sum)
}
