package main

import (
  "bufio"
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

const x = 0
const m = 1
const a = 2
const s = 3

const lt = 0
const gt = 1

func min(a, b int) int {
  if a < b {
    return a
  }
  return b
}

func max(a, b int) int {
  if a > b {
    return a
  }
  return b
}

type tool struct {
  v map[int]int
}

type op struct {
  prop int
  oper int
  threshold int
  next string
}

type interval struct {
  min, max int
}

type filter struct {
  ins [][]interval
  name string
  ops []op
  dflt string
}

func (f *filter) processInputs(fMap map[string]*filter) []string {
  if debug {
    fmt.Printf("processing: %v\n", *f)
  }
  outs := []string{}
  for _, in := range f.ins {
    if debug {
      fmt.Printf("process input: %v\n", in)
    }
    for _, o := range f.ops {
      if o.oper == lt {
	if o.threshold < in[o.prop].min {
	  continue
	}
	outs = append(outs, o.next)
	if o.threshold > in[o.prop].max {
	  newIns := append(fMap[o.next].ins, in)
	  (*fMap[o.next]).ins = newIns
	  break
	}
	newMap := []interval{in[0], in[1], in[2], in[3]}
	newMap[o.prop].max = o.threshold - 1
	fMap[o.next].ins = append(fMap[o.next].ins, newMap)
	in[o.prop].min = o.threshold
      } else {  // o.oper == gt
	if o.threshold > in[o.prop].max {
	  continue
	}
	outs = append(outs, o.next)
	if o.threshold < in[o.prop].min {
	  fMap[o.next].ins = append(fMap[o.next].ins, in)
	  break
	}
	newMap := []interval{in[0], in[1], in[2], in[3]}
	newMap[o.prop].min = o.threshold + 1
	if debug {
	  fmt.Printf("inserting into %v: %v\n", o.next, fMap[o.next])
	}
	fMap[o.next].ins = append(fMap[o.next].ins, newMap)
	in[o.prop].max = o.threshold
      }
    }
    outs = append(outs, f.dflt)
    fMap[f.dflt].ins = append(fMap[f.dflt].ins, in)
  }
  return outs
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

  filters := map[string]*filter{}
  for scan.Scan() {
    if scan.Text() == "" {
      break
    }
    split := strings.Split(scan.Text(), "{")
    label := split[0]
    split = strings.Split(strings.Split(split[1], "}")[0], ",")
    f := filter{name: label}
    for _, rule := range split {
      if !strings.Contains(rule, ":") {
	f.dflt = rule
	continue
      }
      secondSplit := strings.Split(rule, ":")
      prop := x
      sign := lt
      bufVal := []rune{}
      for i, c := range secondSplit[0] {
	if i == 0 {
	  switch c{
	  case 'x': prop = x
	  case 'm': prop = m
          case 'a': prop = a
	  case 's': prop = s
	  }
	  continue
	}
	if i == 1 {
	  if c == '>' {
	    sign = gt
	  }
	  continue
	}
	bufVal = append(bufVal, c)
      }
      v, _ := strconv.Atoi(string(bufVal))
      f.ops = append(f.ops, op{prop, sign, v, secondSplit[1]})
    }
    filters[label] = &f
  }
  aFilt := filter{}
  rFilt := filter{}
  filters["A"] = &aFilt
  filters["R"] = &rFilt

  queue := []string{"in", "A"}
  filters["in"].ins = [][]interval{{
    {min:1, max:4000},
    {min:1, max:4000},
    {min:1, max:4000},
    {min:1, max:4000},
  }}

  for len(queue) > 1 {
    pop := queue[0]
    queue = queue[1:]
    if pop == "R" {
      continue
    }
    if pop == "A" {
      found := false
      for _, v := range queue {
	if v == "A" {
	  found = true
	  break
	}
      }
      if !found {
	queue = append(queue, "A")
      }
      continue
    }

    filter := filters[pop]
    outs := filter.processInputs(filters)
    queue = append(queue, outs...)
    if debug {
      fmt.Printf("------\n%v\n", queue)
      fmt.Printf("%v\n+++++++\n", filters)
    }
  }

  if debug {
    fmt.Printf("A filter: %v\n", filters["A"])
  }

  sum := 0
  for _, in := range filters["A"].ins {
    part := 1
    for _, interval := range in {
      part *= (interval.max - interval.min + 1)
    }
    sum += part
  }

  fmt.Printf("%v\n", sum)

/*
  prevs := map[string][]string{}
  for _, f := range filters {
    for _, o := range f.ops {
      prevs[o.next] = append(prevs[o.next], f.name)
    }
    prevs[f.dflt] = append(prevs[f.dflt], f.name)
  }

  for k, v := range prevs {
    if len(v) > 1 {
      fmt.Printf("long prev for %v: %v\n\n", k, v)
    }
  }
  fmt.Printf("prevs: %v\n", prevs)
*/
}
