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

const empty = "."
const up_down = "\\"
const down_up = "/"
const hsplit = "-"
const vsplit = "|"

const n = 0
const e = 1
const s = 2
const w = 3

type field struct {
  typ string
  light_in [4]bool
  light_out [4]bool
}

func max(a, b int) int {
  if a > b {
    return a
  }
  return b
}

func printMap(m [][]field, printLight bool) {
  if !debug {
    return
  }
  for _, r := range m {
    for _, v := range r {
      if v.typ != empty {
	fmt.Print(v.typ)
      } else {
	if v.light_in[0] || v.light_in[1] || v.light_in[2] || v.light_in[3] {
	  fmt.Print("#")
	} else {
	  fmt.Print(empty)
	}
      }
    }
    fmt.Println()
  }
  fmt.Println()
}

func printEnergized(m [][]field) {
  if !debug {
    return
  }
  for _, r := range m {
    for _, v := range r {
      if v.light_in[0] || v.light_in[1] || v.light_in[2] || v.light_in[3] {
	fmt.Print("#")
      } else {
	fmt.Print(empty)
      }
    }
    fmt.Println()
  }
  fmt.Println()
}

type ray struct {
  x, y, from int
}

func passRay(r ray, f field) []ray {
  ret := []ray{}
  switch r.from {
  case n:
    switch f.typ {
    case empty, vsplit: ret = append(ret, ray{r.x, r.y+1, n})
    case up_down: ret = append(ret, ray{r.x+1, r.y, w})
    case down_up: ret = append(ret, ray{r.x-1, r.y, e})
    case hsplit: ret = append(ret, ray{r.x-1, r.y, e}, ray{r.x+1, r.y, w})
    default: panic("aa")
    }
  case w:
    switch f.typ {
    case empty, hsplit: ret = append(ret, ray{r.x+1, r.y, w})
    case up_down: ret = append(ret, ray{r.x, r.y+1, n})
    case down_up: ret = append(ret, ray{r.x, r.y-1, s})
    case vsplit: ret = append(ret, ray{r.x, r.y-1, s}, ray{r.x, r.y+1, n})
    default: panic("bb")
    }
  case s:
    switch f.typ {
    case empty, vsplit: ret = append(ret, ray{r.x, r.y-1, s})
    case up_down: ret = append(ret, ray{r.x-1, r.y, e})
    case down_up: ret = append(ret, ray{r.x+1, r.y, w})
    case hsplit: ret = append(ret, ray{r.x+1, r.y, w}, ray{r.x-1, r.y, e})
    default: panic("cc")
    }
  case e:
    switch f.typ {
    case empty, hsplit: ret = append(ret, ray{r.x-1, r.y, e})
    case up_down: ret = append(ret, ray{r.x, r.y-1, s})
    case down_up: ret = append(ret, ray{r.x, r.y+1, n})
    case vsplit: ret = append(ret, ray{r.x, r.y+1, n}, ray{r.x, r.y-1, s})
    default: panic("dd")
    }
  }
  return ret
}

func getEnergy(original [][]field, in ray) int {
  m := make([][]field, len(original))
  for i := range original {
    for _, v := range original[i] {
      m[i] = append(m[i], v)
    }
  }

  wavefront := []ray{in}
  for len(wavefront) > 0 {
    pop := wavefront[0]
    wavefront = wavefront[1:]
    if m[pop.y][pop.x].light_in[pop.from] {
      continue
    }
    m[pop.y][pop.x].light_in[pop.from] = true
    nextDir := passRay(pop, m[pop.y][pop.x])
    for _, r := range nextDir {
      if r.x < 0 || r.x >= len(m[0]) || r.y < 0 || r.y >= len(m) {
	continue
      }
      if m[r.y][r.x].light_in[r.from] {
	continue
      }
      wavefront = append(wavefront, r)
    }
  }

  sum := 0
  for _, r := range m {
    for _, v := range r {
      if v.light_in[0] || v.light_in[1] || v.light_in[2] || v.light_in[3] {
	sum++
      }
    }
  }

  return sum
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

  mapa := [][]field{}
  for scan.Scan() {
    row := []field{}
    for _, c := range scan.Text() {
      row = append(row, field{typ: string(c)})
    }
    mapa = append(mapa, row)
  }

  maximum := 0
  for x := 0; x < len(mapa[0]); x++ {
    maximum = max(maximum, getEnergy(mapa, ray{x, 0, n}))
    maximum = max(maximum, getEnergy(mapa, ray{x, len(mapa) - 1, s}))
  }
  for y := 0; y < len(mapa); y++ {
    maximum = max(maximum, getEnergy(mapa, ray{0, y, w}))
    maximum = max(maximum, getEnergy(mapa, ray{len(mapa[0]) - 1, y, e}))
  }

  fmt.Printf("%v\n", maximum)
}
