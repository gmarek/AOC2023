package main

import (
  "bufio"
  "flag"
  "fmt"
//  "math"
  "os"
  "sort"
//  "strconv"
  "strings"
//  "unicode"
)

var prod = flag.Bool("real", false, "run on real input")
var debug = false

type signal struct {
  high bool
  source, target string
}

func (s *signal) printSigal() {
  fmt.Printf("%v -> %v: %v\n", s.source, s.target, s.high)
}

func gcd(a, b int) int {
  for b != 0 {
    t := b
    b = a % b
    a = t
  }
  return a
}

func lcm(a, b int) int {
  result := a * b / gcd(a, b)
  return result
}

const high = true
const low = false

type module interface {
  process(signal) []signal
  addInput(string)
  printModule()
  outs() []string
  ins() []string
}

type flipFlop struct {
  name string
  state bool
  out []string
  in map[string]bool
}

func (f *flipFlop) process(s signal) []signal {
  if s.high == high {
    return []signal{}
  }
  ret := []signal{}
  for _, o := range f.out {
    ret = append(ret, signal{!f.state, f.name, o})
  }
  f.state = !f.state
  return ret
}

func (f *flipFlop) addInput(s string) {
  f.in[s] = true
}

func (f *flipFlop) printModule() {
  fmt.Printf("%v\n", *f)
}

func (f *flipFlop) outs() []string {
  return f.out
}

func (f *flipFlop) ins() []string {
  ret := []string{}
  for k := range f.in {
    ret = append(ret, k)
  }
  return ret
}

type conjunction struct {
  name string
  inState map[string]bool
  out []string
}

func (c *conjunction) addInput(s string) {
  c.inState[s] = false
}

func (c *conjunction) printModule() {
  fmt.Printf("%v\n", *c)
}

func (c *conjunction) process(s signal) []signal {
  c.inState[s.source] = s.high
  all := true
  for _, h := range c.inState {
    if !h {
      all = false
      break
    }
  }
  ret := []signal{}
  for  _, o := range c.out {
    ret = append(ret, signal{!all, c.name, o})
  }
  return ret
}

func (c *conjunction) outs() []string {
  return c.out
}

func (c *conjunction) ins() []string {
  ret := []string{}
  for k := range c.inState {
    ret = append(ret, k)
  }
  return ret
}

type broadcast struct {
  name string
  out []string
}

func (b *broadcast) addInput(s string) {
  return
}

func (b *broadcast) printModule() {
  fmt.Printf("%v\n", *b)
}

func (b *broadcast) process(s signal) []signal {
  ret := []signal{}
  for _, o := range b.out {
    ret = append(ret, signal{s.high, b.name, o})
  }
  return ret
}

func (b *broadcast) outs() []string {
  return b.out
}

func (b *broadcast) ins() []string {
  return []string{"button"}
}

type button struct {
}

func (b *button) process(s signal) []signal {
  return []signal{{low, "button", "broadcast"}}
}

func (b *button) addInput(s string) {
  return
}

func (b *button) printModule() {
  fmt.Println("button")
}

func (b *button) outs() []string {
  return []string{"broadcast"}
}

func (b *button) ins() []string {
  return []string{}
}

type sink struct {
  in []string
}

func (s *sink) process(signal) []signal {
  return []signal{}
}

func (s *sink) addInput(name string) {
  s.in = append(s.in, name)
}

func (s *sink) outs() []string {
  return []string{}
}

func (s *sink) ins() []string {
  return s.in
}

func (s *sink) printModule() {
  fmt.Printf("sink: %v\n", s.in)
}

func expand(target string, modules map[string]module) []string {
  queue := []string{target}
  visited := map[string]bool{}
  for len(queue) > 0 {
    pop := queue[0]
    visited[pop] = true
    queue = queue[1:]
    ins := modules[pop].ins()
    for _, s := range ins {
      if !visited[s] {
	queue = append(queue, s)
      }
    }
  }
  ret := []string{}
  for k := range visited {
    ret = append(ret, k)
  }
  sort.Strings(ret)
  return ret
}

func main() {
  flag.Parse()

  filename := "test2"
  if *prod {
    filename = "input"
  }
  f, err := os.Open(filename)
  if err != nil {
    panic(err)
  }

  scan := bufio.NewScanner(bufio.NewReader(f))

  modules := map[string]module{}
  for scan.Scan() {
    tokens := strings.Split(scan.Text(), "->")
    typ := 0
    if !strings.HasPrefix(tokens[0], "broad") {
      if tokens[0][0] == '%' {
	typ = 1
      } else {
	typ = 2
      }
    }
    outs := strings.Split(tokens[1], ",")
    for i := range outs {
      outs[i] = strings.TrimSpace(outs[i])
    }
    name := "broadcast"
    var mod module
    switch typ {
    case 0: // broadcast
    {
      mod = &broadcast{name, outs}
    }
    case 1: // flipFlop
    {
      name = strings.TrimSpace(tokens[0][1:])
      mod = &flipFlop{name, false, outs, map[string]bool{}}
    }
    case 2:
    {
      name = strings.TrimSpace(tokens[0][1:])
      mod = &conjunction{name, map[string]bool{}, outs}
    }
    }

    modules[name] =  mod
  }

  modules["button"] = &button{}
  modules["rx"] = &sink{}
  for k, v := range modules {
    for _, o := range v.outs() {
      if _, ok := modules[o]; ok {
	modules[o].addInput(k)
      }
    }
  }

  for _, k := range expand("hr", modules) {
    fmt.Printf("%v, ", k)
  }
  fmt.Println()
  fmt.Printf("----\n")
  for _, k := range expand("nr", modules) {
    fmt.Printf("%v, ", k)
  }
  fmt.Println()
  fmt.Printf("----\n")
  for _, k := range expand("gl", modules) {
    fmt.Printf("%v, ", k)
  }
  fmt.Println()
  fmt.Printf("----\n")
  for _, k := range expand("gk", modules) {
    fmt.Printf("%v, ", k)
  }
  fmt.Println()
  fmt.Printf("----\n")

  targets := map[string]int{"hr":0, "nr":0, "gl":0, "gk":0}
  pressed := 0
  found := false
  for !found {
    queue := []signal{{low, "fake_start", "button"}}
    pressed++
    for len(queue) > 0 {
      pop := queue[0]
      queue = queue[1:]
      signals := modules[pop.target].process(pop)
      for _, s := range signals {
	if s.target == "rx" {
	  if !s.high {
	    found = true
	  }
	} else {
	  if v, ok := targets[s.source]; ok && v == 0 && !s.high {
	    targets[s.source] = pressed
	  }
	  queue = append(queue, s)
	}
      }
    }
    found = true
    for _, v := range targets {
      if v == 0 {
	found = false
      }
    }
  }

  fmt.Printf("%v\n", targets)
  val := 1
  for _, v := range targets {
    val = lcm(val, v)
  }
  fmt.Printf("%v\n", val)
}
