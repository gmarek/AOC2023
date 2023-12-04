package main

import (
  "bufio"
  "flag"
  "fmt"
  "os"
  "strconv"
  "strings"
//  "unicode"
)

var prod = flag.Bool("real", false, "run on real input")

func contains(set []int, val int) bool {
  for _, v := range set {
    if v == val {
      return true
    }
  }
  return false
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

  sum := int64(0)
  queue := []int64{}
  for scan.Scan() {
    add := int64(0)
    if len(queue) > 0 {
      add = queue[0]
      queue = queue[1:]
    }
    sum += add + 1

    game := strings.Split(strings.Split(scan.Text(), ":")[1], "|")
    winning, has := []int{}, []int{}
    for _, v := range strings.Split(game[0], " ") {
      if v != "" {
	val, _ := strconv.Atoi(v)
	winning = append(winning, val)
      }
    }
    for _, v := range strings.Split(game[1], " ") {
      if v != "" {
	val, _ := strconv.Atoi(v)
	has = append(has, val)
      }
    }

    hits := 0
    for _, v := range has {
      if contains(winning, v) {
	hits++
      }
    }

    for i:=0; i < hits; i++ {
      if len(queue) > i {
	queue[i] += add + 1
      } else {
	queue = append(queue, 1+add)
      }
    }
  }

  fmt.Printf("%v\n", sum)
}
