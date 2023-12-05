package main

import (
  "bufio"
  "flag"
  "fmt"
  "os"
  "sort"
  "strconv"
  "strings"
  //"unicode"
)

var prod = flag.Bool("real", false, "run on real input")
var debug = false

type seedsVec struct {
  start, l int64
}

type mapping struct {
  src, offset, l int64
}

func min(a, b int64) int64 {
  if a < b {
    return a
  }
  return b
}

func max(a, b int64) int64 {
  if a >= b {
    return a
  }
  return b
}

func parseAndInsertMapping(m *[]mapping, s, prefix string) {
  if strings.HasPrefix(s, prefix) {
    return
  }
  nums := strings.Split(s, " ")
  dst, _ := strconv.Atoi(nums[0])
  src, _ := strconv.Atoi(nums[1])
  l, _ := strconv.Atoi(nums[2])
  *m = append(*m, mapping{int64(src), int64(dst-src), int64(l)})
}

func processMapping(seeds []seedsVec, mpg []mapping) []seedsVec {
  ret := []seedsVec{}
  for _, vec := range seeds {
    for _, m := range mpg {
      if debug {
	fmt.Printf("before: %v %v\n", vec, m)
      }
      if vec.start > m.src && vec.start < m.src + m.l {
	if debug {
	  fmt.Printf("first: %v %v\n", vec, m)
	}
	end := min(vec.start + vec.l - 1, m.src + m.l - 1)
	ret = append(ret, seedsVec{vec.start + m.offset, end - vec.start + 1})
	if end == vec.start + vec.l - 1 {
	  vec.start = -1
	  break
	}
	vec.l -= end - vec.start + 1
	vec.start = end + 1
      }
      if vec.start <= m.src && vec.start + vec.l >= m.src {
	if debug {
	  fmt.Printf("second: %v %v\n", vec, m)
	}
	if vec.start < m.src {
	  ret = append(ret, seedsVec{vec.start, m.src - vec.start + 1})
	}
	end := min(vec.start + vec.l - 1, m.src + m.l - 1)
	ret = append(ret, seedsVec{vec.start + m.offset, end - vec.start + 1})
	if end == vec.start + vec.l - 1 {
	  vec.start = -1
	  break
	}
	vec.l -= end - vec.start + 1
	vec.start = end + 1
      }
    }
    if vec.start != -1 {
      ret = append(ret, vec)
    }
  }
  if debug {
    fmt.Printf("%v\n", ret)
  }

  return ret
}

func sortMapping(m *[]mapping) {
  sort.Slice(*m, func(i, j int) bool {
    return (*m)[i].src < (*m)[j].src || ((*m)[i].src == (*m)[j].src && (*m)[i].offset < (*m)[j].offset)
  })
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

  seeds := []seedsVec{}
  seedToSoil := []mapping{}
  soilToFertilizer := []mapping{}
  fertilizerToWater := []mapping{}
  waterToLight := []mapping{}
  lightToTemperature := []mapping{}
  temperatureToHumidity := []mapping{}
  humidityToLocation := []mapping{}
  phase := 1
  for scan.Scan() {
    if scan.Text() == "" {
      phase++
      continue
    }
    if phase == 1 {
      s := strings.Split(scan.Text(), ":")[1]
      numS := strings.Split(s, " ")
      numSS := []string{}
      for _, k := range numS {
	if k != "" {
	  numSS = append(numSS, k)
	}
      }
      for i := 0; i < len(numSS); i+=2 {
	if numSS[i] == "" {
	  panic("empty string")
	}
	start, _ := strconv.Atoi(numSS[i])
	l, _ := strconv.Atoi(numSS[i+1])
	seeds = append(seeds, seedsVec{int64(start), int64(l)})
      }
    } else if phase == 2 {
      parseAndInsertMapping(&seedToSoil, scan.Text(), "seed-to-soil")
    } else if phase == 3 {
      parseAndInsertMapping(&soilToFertilizer, scan.Text(), "soil-to")
    } else if phase == 4 {
      parseAndInsertMapping(&fertilizerToWater, scan.Text(), "fertilizer")
    } else if phase == 5 {
      parseAndInsertMapping(&waterToLight, scan.Text(), "water")
    } else if phase == 6 {
      parseAndInsertMapping(&lightToTemperature, scan.Text(), "light")
    } else if phase == 7 {
      parseAndInsertMapping(&temperatureToHumidity, scan.Text(), "temper")
    } else if phase == 8 {
      parseAndInsertMapping(&humidityToLocation, scan.Text(), "humid")
    }
  }

  sortMapping(&seedToSoil)
  sortMapping(&soilToFertilizer)
  sortMapping(&fertilizerToWater)
  sortMapping(&waterToLight)
  sortMapping(&lightToTemperature)
  sortMapping(&temperatureToHumidity)
  sortMapping(&humidityToLocation)

  fmt.Printf("%v\n", seeds)
  seeds = processMapping(seeds, seedToSoil)
  seeds = processMapping(seeds, soilToFertilizer)
  seeds = processMapping(seeds, fertilizerToWater)
  seeds = processMapping(seeds, waterToLight)
  seeds = processMapping(seeds, lightToTemperature)
  seeds = processMapping(seeds, temperatureToHumidity)
  seeds = processMapping(seeds, humidityToLocation)

  sort.Slice(seeds, func(i, j int) bool { return seeds[i].start < seeds[j].start })
  fmt.Printf("%v\n", seeds[0].start)
}
