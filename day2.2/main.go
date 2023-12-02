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

  sum := int64(0)
  for scan.Scan() {
    game := strings.Split(strings.Split(scan.Text(), ":")[1], ";")
    minRed := int64(-1)
    minBlue := int64(-1)
    minGreen := int64(-1)
    for _, g := range game {
      balls := strings.Split(g, ",")
      for _, b := range balls {
	b = strings.TrimSpace(b)
	color := ""
	num := int64(0)
	fmt.Sscanf(b, "%d %s", &num, &color)
	switch (color) {
	case "red": if num > minRed {minRed = num}
	case "green": if num > minGreen {minGreen = num}
	case "blue": if num > minBlue {minBlue = num}
	}
      }
    }
    sum += minRed * minBlue * minGreen
  }
  fmt.Printf("%d\n", sum)
}
