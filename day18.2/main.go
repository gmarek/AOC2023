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
  orientation int
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
  from int
  inTurn bool
}

func abs(a int) int {
  if a < 0 {
    return -a
  }
  return a
}

func nextInd(i int, s []point) int {
  return (i+1) % len(s)
}

func prevInd(i int, s []point) int {
  if i == 0 {
    return len(s)-1
  }
  return i-1
}

func neigh(i int, s []point) []int {
  return []int{prevInd(i, s), nextInd(i, s)}
}

func orientation(prev, cur int) int {
  turnMap := map[int]point{
    n: point{0, -1, 0, true},
    s: point{0, 1, 0, true},
    e: point{1, 0, 0, true},
    w: point{-1, 0, 0, true},
  }

  if turnMap[prev].x * turnMap[cur].y - turnMap[prev].y * turnMap[cur].x < 0 {
    return right
  }
  return left
}

func cutRectangle(s []point) (int, [][]point) {
  if debug {
    fmt.Println("-----------------")
    fmt.Printf("in %v: %v\n", len(s), s)
  }
  minXInd := map[int]bool{}
  minX := int(1e10)
  nearMinX := int(1e10)
  for i, p := range s {
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

  if len(s) == 4 {
    r1 := abs(nearMinX - minX)
    r2 := 0
    for i := range s {
      val := abs((s)[i].y - (s)[(i+1)%len(s)].y)
      if val > r2 {
	r2 = val
      }
    }
    if debug {
      fmt.Printf("last: %v * %v = %v\n", r1, r2, r1*r2)
    }
    ret := r1 * r2
    return ret, [][]point{}
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
      fmt.Printf("(%v) %v -> (%v) %v (%v) %v\n", i, s[i], ns[0], s[ns[0]], ns[1], s[ns[1]])
    }
    if minXInd[ns[0]] {
      r1 := nearMinX - minX
      r2 := abs(s[i].y - s[ns[0]].y)
      if debug {
	fmt.Printf("0) %v * %v = %v\n", r1, r2, r1*r2)
	fmt.Printf("%v %v\n", s[ns[0]], s[i])
      }
      ret += r1 * r2
    }
    if minXInd[ns[1]] {
      r1 := nearMinX - minX
      r2 := abs(s[i].y - s[ns[1]].y)
      if debug {
	fmt.Printf("1) %v * %v = %v\n", r1, r2, r1*r2)
	fmt.Printf("%v %v\n", s[ns[1]], s[i])
      }
      ret += r1 * r2
    }
  }

  for i := range minXInd {
    s[i].x = nearMinX
  }
  if debug {
    fmt.Printf("post squash: %v\n", s)
  }

  for i := range minXInd {
    ns := neigh(i, s)
    if s[i].x == s[ns[0]].x && s[i].y == s[ns[0]].y {
      if debug {
	fmt.Printf("removing duplicate %v: %v\n", i, s[i])
      }
      remMap[i] = true
    }
    if s[i].x == s[ns[1]].x && s[i].y == s[ns[1]].y {
      if debug {
	fmt.Printf("removing duplicate %v: %v\n", i, s[i])
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
    remove(&s, remList[i])
  }

  if debug {
    fmt.Printf("no duplicates: %v\n", s)
  }

  slice := map[int]bool{}
  for i, v := range s {
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
      if s[k].x != s[n].x {
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
    remove(&s, remList[i])
  }

  // detect split...
  retShapes := [][]point{}

  minList := []int{}
  for i, v := range s {
    if v.x == nearMinX {
      minList = append(minList, i)
    }
  }
  sort.Ints(minList)

  splitOn := []int{}
  for _, i := range minList {
    if !s[i].inTurn {
      splitOn = append(splitOn, i)
    }
  }

  if len(splitOn) == 0 {
    if len(s) > 0 {
      retShapes = append(retShapes, s)
    }
  } else {

    mapped := map[int]bool{}
    for _, i := range splitOn {
      connected := []point{s[i]}
      mapped[i] = true
      forward := false
      n := nextInd(i, s)
      if s[n].x != nearMinX {
	forward = true
      }
      if forward {
        for s[n].x != nearMinX {
	  connected = append(connected, s[n])
	  mapped[n] = true
	  n = nextInd(n, s)
	}
      } else {
	n = prevInd(i, s)
	for s[n].x != nearMinX {
	  connected = append(connected, s[n])
	  mapped[n] = true
	  n = prevInd(n, s)
        }
      }
      connected = append(connected, s[n])
      if debug {
	fmt.Printf("connected: %v\n", connected)
      }
      retShapes = append(retShapes, connected)
    }

    remS := []point{}
    for i := range s {
      if !mapped[i] {
	remS = append(remS, s[i])
      }
    }
    if debug {
      fmt.Printf("remained: %v\n", remS)
    }
  }

  // post split...

  if debug {
    fmt.Printf("return: %v\n", ret/2)
  }
  return ret / 2, retShapes
}

func remove(s *[]point, i int) {
  if debug {
    fmt.Printf("remove: %v (%v) from %v\n", i, (*s)[i], *s)
  }
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
    dString := ""

    fmt.Sscanf(scan.Text(), "%s %d", &dString, &d.l)
    switch dString {
    case "R": d.dir = e
    case "D": d.dir = s
    case "L": d.dir = w
    case "U": d.dir = n
    }
    /*
    sTrash := ""
    iTrash := 0
    fmt.Sscanf(scan.Text(), "%s %d %s", &sTrash, &iTrash, &dString)
    lString := dString[2:len(dString) - 2]
    fmt.Sscanf(lString, "%x", &d.l)
    dirStr := dString[len(dString)-2:]
    fmt.Sscanf(dirStr, "%d", &d.dir)
    */
    digs = append(digs, d)
  }

  countLeft := 0
  for i, d := range digs {
    digs[i].orientation = orientation(d.dir, digs[(i+1)%len(digs)].dir)
    if digs[i].orientation == left {
      countLeft++
    } else {
      countLeft--
    }
  }

  globalTurn := 0
  if countLeft == 4 {
    globalTurn = left
  } else if countLeft == -4 {
    globalTurn = right
  } else {
    panic(fmt.Sprintf("countLeft: %v\n", countLeft))
  }

  shape := []point{{0, 0, digs[len(digs)-1].dir, digs[len(digs)-1].orientation == globalTurn}}
  for i, d := range digs {
    pInd := i-1
    if pInd < 0 {
      pInd = len(digs)-1
    }
    offset := 0
    if digs[pInd].orientation == d.orientation {
      if d.orientation == globalTurn {
	offset = 1
      } else {
	offset = -1
      }
    }
    if debug {
      fmt.Printf("dig: %v, offset: %v\n", d, offset)
    }
    next := point{shape[len(shape) - 1].x, shape[len(shape) - 1].y, d.dir, d.orientation == globalTurn}
    switch d.dir {
    case n: {
      next.y -= d.l + offset
    }
    case e: {
      next.x += d.l + offset
    }
    case s: {
      next.y += d.l + offset
    }
    case w: {
      next.x -= d.l + offset
    }
    }
    shape = append(shape, next)
  }

  fmt.Printf("shape: %v\n", shape)
  shape = shape[:len(shape)-1]
  shapes := [][]point{shape}
  sum := 0
  for len(shapes) > 0 {
    newShapes := [][]point{}
    for _, s := range shapes {
      v, newShape := cutRectangle(s)
      sum += v
      for _, ns := range newShape {
	newShapes = append(newShapes, ns)
      }
    }
    shapes = newShapes
  }

  fmt.Printf("%v\n", sum)
}
