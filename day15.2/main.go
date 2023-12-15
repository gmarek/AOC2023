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

func hash(s string) int {
  val := 0
  for _, c := range s {
    val += int(c)
    val *= 17
    val = val % 256
  }
  return val
}

type lens struct {
  label string
  focal int
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

  tokens := []string{}
  for scan.Scan() {
    tokens = strings.Split(scan.Text(), ",")
  }

  boxes := [256][]lens{}
  for _, t := range tokens {
    if t[len(t)-1] == '-' {
      label := t[:len(t)-1]
      id := hash(label)
      for i := range boxes[id] {
	if boxes[id][i].label == label {
	  boxes[id] = append(boxes[id][:i], boxes[id][i+1:]...)
	  break
	}
      }
      continue
    }
    lensId := strings.Split(t, "=")
    label := lensId[0]
    focal, _ := strconv.Atoi(lensId[1])
    found := false
    id := hash(label)
    for i := range boxes[id] {
      if boxes[id][i].label == label {
	boxes[id][i].focal = focal
	found = true
	break
      }
    }
    if !found {
      boxes[id] = append(boxes[id], lens{label, focal})
    }
  }

  if debug {
    for i, b := range boxes {
      if len(b) > 0 {
	fmt.Printf("%v: %v\n", i, b)
      }
    }
  }

  sum := 0
  for i, b := range boxes {
    for j, v := range b {
      sum += (i+1)*(j+1)*v.focal
    }
  }

  fmt.Printf("%v\n", sum)
      
}
