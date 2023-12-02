package main

import (
  "bufio"
  "flag"
  "fmt"
  "os"
//  "strconv"
  "strings"
//  "unicode"
)

var prod = flag.Bool("real", false, "run on real input")

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

  maxBlue := 14
  maxRed := 12
  maxGreen := 13

  sum := 0
  i := 0
  for scan.Scan() {
    i += 1
    game := strings.Split(strings.Split(scan.Text(), ":")[1], ";")
    ok := true
    for _, g := range game {
      balls := strings.Split(g, ",")
      for _, b := range balls {
	b = strings.TrimSpace(b)
	color := ""
	num := 0
	fmt.Sscanf(b, "%d %s", &num, &color)
	switch (color) {
	case "red": if num > maxRed {ok = false}
	case "green": if num > maxGreen {ok = false}
	case "blue": if num > maxBlue {ok = false}
	}
      }
    }
    if ok {
      sum += i
    }
  }
  fmt.Printf("%d\n", sum)
}
