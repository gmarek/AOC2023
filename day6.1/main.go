package main

import (
  "bufio"
  "flag"
  "fmt"
  "math"
  "os"
//  "strconv"
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

  times := []int{}
  for len(timeLine) > 0 {
    v := 0
    fmt.Sscanf(timeLine, "%d", &v)
    times = append(times, v)
    timeLine = strings.TrimSpace(timeLine[len(fmt.Sprintf("%d", v)):])
  }

  distances := []int{}
  for len(distanceLine) > 0 {
    v := 0
    fmt.Sscanf(distanceLine, "%d", &v)
    distances = append(distances, v)
    distanceLine = strings.TrimSpace(distanceLine[len(fmt.Sprintf("%d", v)):])
  }

  races := []race{}
  for i := 0; i < len(times); i++ {
    races = append(races, race{float64(times[i]), float64(distances[i])})
  }

  sum := 1
  for _, r := range(races) {
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
    fmt.Printf("%v\n", x2 - x1 + 1)
    sum *= int(x2-x1+1)
  }

  fmt.Printf("%v\n", sum)
}
