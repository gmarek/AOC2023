package main

import (
  "bufio"
  "flag"
  "fmt"
//  "math"
  "os"
//  "sort"
//  "strconv"
  "strings"
//  "unicode"
)

var prod = flag.Bool("real", false, "run on real input")
var debug = true

type signal struct {
  high bool
  source, target string
}

func (s *signal) printSigal() {
  fmt.Printf("%v -> %v: %v\n", s.source, s.target, s.high)
}

const high = true
const low = false

type module interface {
  process(signal) []signal
  addInput(string)
  printModule()
  outs() []string
}

type flipFlop struct {
  name string
  state bool
  out []string
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
  return
}

func (f *flipFlop) printModule() {
  fmt.Printf("%v\n", *f)
}

func (f *flipFlop) outs() []string {
  return f.out
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

func (b* button) outs() []string {
  return []string{"broadcast"}
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
      mod = &flipFlop{name, false, outs}
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
  for k, v := range modules {
    for _, o := range v.outs() {
      if _, ok := modules[o]; ok {
	modules[o].addInput(k)
      }
    }
  }

  /*
  if debug {
    for _, m := range modules {
      m.printModule()
    }
  }*/

  sumHigh := 0
  sumLow := 0
  for i := 0; i < 1000; i++ {
    queue := []signal{{low, "fake_start", "button"}}
    for len(queue) > 0 {
      pop := queue[0]
      if debug {
	//pop.printSigal()
      }
      if pop.high {
	sumHigh++
      } else {
	if pop.source != "fake_start" {
	  sumLow++
	}
      }
      queue = queue[1:]
      signals := modules[pop.target].process(pop)
      for _, s := range signals {
	if s.target == "output" || s.target == "rx"{
	  if s.high {
	    sumHigh++
	  } else {
	    sumLow++
	  }
	} else {
	  queue = append(queue, s)
	}
      }
    }
  }

  fmt.Printf("high: %v, low: %v -> %v\n", sumHigh, sumLow, sumHigh * sumLow)
}
