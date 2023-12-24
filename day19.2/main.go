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
var debug = true

const x = 0
const m = 1
const a = 2
const s = 3

const lt = 0
const gt = 1

type tool struct {
  v map[int]int
}

type op struct {
  prop int
  oper int
  threshold int
  next string
}

func (o *op) process(t tool) string {
  v := t.v[o.prop]
  if o.oper == lt {
    if v < o.threshold {
      return o.next
    }
  } else {
    if v > o.threshold {
      return o.next
    }
  }
  return ""
}

type filter struct {
  name string
  ops []op
  dflt string
}

func (f *filter) processTool(t tool) string {
  ret := ""
  for _, o := range f.ops {
    ret = o.process(t)
    if ret != "" {
      return ret
    }
  }
  return f.dflt
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

  scanTools := false
  tools := []tool{}
  filters := map[string]filter{}
  for scan.Scan() {
    if scan.Text() == "" {
      scanTools = true
      continue
    }
    if scanTools {
      t := tool{v: map[int]int{}}
      stripped := strings.TrimSuffix(strings.TrimPrefix(scan.Text(), "{"), "}")
      split := strings.Split(stripped, ",")
      for i := range split {
	v, _ := strconv.Atoi(split[i][2:])
	t.v[i] = v
      }
      tools = append(tools, t)
      continue
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
    filters[label] = f
  }

  sum := 0
  for _, t := range tools {
    f := "in"
    for f != "A" && f != "R" {
      tmp := filters[f]
      f = tmp.processTool(t)
    }
    if f == "A" {
      for _, v := range t.v {
	sum += v
      }
    }
  }

  fmt.Printf("%v\n", sum)
}
