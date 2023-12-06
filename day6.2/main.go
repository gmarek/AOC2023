package main

import (
  "bufio"
  "flag"
  "fmt"
  "math"
  "os"
  "strconv"
  "strings"
//  "unicode"
)

var prod = flag.Bool("real", false, "run on real input")

type race struct {
  time, distance float64
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

  scan.Scan()
  timeLine := scan.Text()
  scan.Scan()
  distanceLine := scan.Text()

  timeLine = strings.TrimSpace(strings.TrimPrefix(timeLine, "Time:"))
  distanceLine = strings.TrimSpace(strings.TrimPrefix(distanceLine, "Distance:"))

  times := []string{}
  for len(timeLine) > 0 {
    s := ""
    fmt.Sscanf(timeLine, "%s", &s)
    times = append(times, s)
    timeLine = strings.TrimSpace(timeLine[len(s):])
  }

  distances := []string{}
  for len(distanceLine) > 0 {
    s := ""
    fmt.Sscanf(distanceLine, "%s", &s)
    distances = append(distances, s)
    distanceLine = strings.TrimSpace(distanceLine[len(s):])
  }

  timeString := strings.Join(times, "")
  distanceString := strings.Join(distances, "")

  t, err := strconv.Atoi(timeString)
  if err != nil {
    panic(err)
  }
  d, err := strconv.Atoi(distanceString)
  if err != nil {
    panic(err)
  }

  r := race{float64(t), float64(d)}

  fmt.Printf("%v\n", r)
  x1 := (r.time - (math.Sqrt(r.time*r.time - 4*r.distance)))/2
  rem := math.Ceil(x1) - x1
  x1 = math.Ceil(x1)
  if rem < 1e-20 {
    x1++
  }
  x2 := (r.time + (math.Sqrt(r.time*r.time - 4*r.distance)))/2
  rem = x2 - math.Floor(x2)
  x2 = math.Floor(x2)
  if rem < 1e-20 {
    x2--
  }
  fmt.Printf("%v\n", int(x2 - x1 + 1))
}
