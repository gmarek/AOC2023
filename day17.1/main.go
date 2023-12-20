package main

import (
  "bufio"
  "container/heap"
  "flag"
  "fmt"
//  "math"
  "os"
//  "sort"
  "strconv"
  "strings"
//  "unicode"
)

var prod = flag.Bool("real", false, "run on real input")
var debug = false

const ver = 0
const hor = 1

type minHeap []*block

func (h minHeap) Len() int {
  return len(h)
}

func (h minHeap) Less(i, j int) bool {
  return h[i].dist < h[j].dist
}

func (h minHeap) Swap(i, j int) {
  h[i], h[j] = h[j], h[i]
}

func (h *minHeap) Push(x interface{}) {
  *h = append(*h, x.(*block))
}

func (h *minHeap) Pop() interface{} {
  old := *h
  n := len(old)
  x := old[n-1]
  *h = old[:n-1]
  return x
}

func getInteraval(a, b point) []point {
  ret := []point{}
  if a.x > b.x {
    for i := b.x; i < a.x; i++ {
      ret = append(ret, point{i, a.y})
    }
  }
  if a.x < b.x {
    for i := a.x + 1; i <= b.x; i++ {
      ret = append(ret, point{i, a.y})
    }
  }
  if a.y > b.y {
    for i := b.y; i < a.y; i++ {
      ret = append(ret, point{a.x, i})
    }
  }
  if a.y < b.y {
    for i := a.y + 1; i <= b.y; i++ {
      ret = append(ret, point{a.x, i})
    }
  }
  return ret
}

type block struct {
  x, y int
  visited bool
  dist int
  v int
  prev *block
  typ int
}

func printHeap(h []*block) {
  s := []string{}
  for _, v := range h {
    if (*v).dist != 1e10 || (*v).dist != 1e10 {
      s = append(s, fmt.Sprintf("(%v,%v)", (*v).x, (*v).y))
    }
  }
  fmt.Println(strings.Join(s, ","))
}

type point struct {
  x, y int
}

func neigh(m [][]*block, b *block, from int) []point {
  x := (*b).x
  y := (*b).y
  c := []point{{x-3, y}, {x-2, y}, {x-1, y},
		 {x+1, y}, {x+2, y}, {x+3, y},
		 {x, y-1}, {x, y-2}, {x, y-3},
		 {x, y+1}, {x, y+2}, {x, y+3}}
  ret := []point{}
  for _, v := range c {
    if v.x >= 0 && v.x < len(m[0]) && v.y >= 0 && v.y < len(m) {
      if from == -1 {
	ret = append(ret, v)
      } else if from == ver && v.x != x {
	ret = append(ret, v)
      } else if from == hor && v.y != y {
	ret = append(ret, v)
      }
    }
  }
  return ret
}

func printCityV(c [][]*block) {
  for y := range c {
    for x := range c[y] {
      fmt.Printf("%d", (*c[y][x]).v)
    }
    fmt.Println()
  }
}

func printCity(c [][]*block) {
  if !debug {
    return
  }
  for y := range c {
    for x := range c[y] {
      if (*c[y][x]).dist == 1e10 {
	fmt.Print(".. ")
      } else {
	fmt.Printf("%d ", c[y][x].dist)
      }
    }
    fmt.Println()
  }
}

func printPath(c [][]*block) {
  cur := c[len(c) - 1][len(c[0])-1]
  for cur != nil {
    fmt.Printf("(%v, %v)<-\n", cur.x, cur.y)
    cur = c[cur.y][cur.x].prev
  }
  fmt.Printf("(0, 0)\n")
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

  cityH := [][]*block{}
  cityV := [][]*block{}
  y := 0
  for scan.Scan() {
    rowH := []*block{}
    rowV := []*block{}
    x := 0
    for _, c := range scan.Text() {
      v, _ := strconv.Atoi(string(c))
      bH := block{x:x, y: y, v: v, dist: 1e10, typ: hor}
      bV := block{x:x, y: y, v: v, dist: 1e10, typ: ver}
      rowH = append(rowH, &bH)
      rowV = append(rowV, &bV)
      x++
    }
    cityH = append(cityH, rowH)
    cityV = append(cityV, rowV)
    y++
  }
  (*cityH[0][0]).visited = true
  (*cityH[0][0]).dist = 0
  (*cityV[0][0]).visited = true
  (*cityV[0][0]).dist = 0

  current := cityH[0][0]
  queue := minHeap{}
  for y := range cityH {
    for x := range cityH[y] {
      queue = append(queue, cityH[y][x])
      queue = append(queue, cityV[y][x])
    }
  }
  heap.Init(&queue)
  foundH := false
  foundV := false

  for !foundH || !foundV {
    current = heap.Pop(&queue).(*block)
    if current.x == len(cityH[0]) - 1 && current.y == len(cityH) - 1 {
      if current.typ == hor {
	foundH = true
      } else {
	foundV = true
      }
    }
    city := &cityH
    if current.typ == hor {
      city = &cityV
    }
    ns := neigh(*city, current, current.typ)
    if debug {
      fmt.Printf("(%v,%v)%v\n", (*current).x, (*current).y, current.typ)
    }
    (*current).visited = true
    for _, n := range ns {
      if (*(*city)[n.y][n.x]).visited {
	continue
      }
      posCost := (*current).dist
      path := getInteraval(point{(*current).x, (*current).y}, n)
      for _, v := range path {
	posCost += (*(*city)[v.y][v.x]).v
      }
      if (*(*city)[n.y][n.x]).dist > posCost {
	(*(*city)[n.y][n.x]).dist = posCost
	pathCur := point{(*current).x, (*current).y}
	for _, v := range path {
	  if v.x == n.x && v.y == n.y {
	    (*(*city)[v.y][v.x]).prev = current
	  } else {
	    (*(*city)[v.y][v.x]).prev = (*city)[pathCur.y][pathCur.x]
	  }
	  pathCur = v
	}
	for i := range queue {
	  if queue[i] == (*city)[n.y][n.x] {
	    heap.Fix(&queue, i)
	    break
	  }
	}
      }
      if debug {
	printCity(cityH)
	fmt.Println("----")
	printCity(cityV)
	fmt.Println("++++")
      }
    }
  }

  //city := &cityH
  min := cityH[len(cityH) - 1][len(cityH[0]) - 1].dist
  if min > cityV[len(cityV) - 1][len(cityV[0]) - 1].dist {
    min = cityV[len(cityV) - 1][len(cityV[0]) - 1].dist
    //city = &cityV
  }
  fmt.Printf("%v\n", min)
  //printPath(*city)
}
