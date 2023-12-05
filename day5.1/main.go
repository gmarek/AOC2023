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

type mapping struct {
  src, dst, l int64
}

func parseAndInsertMapping(m *[]mapping, s, prefix string) {
  if strings.HasPrefix(s, prefix) {
    return
  }
  nums := strings.Split(s, " ")
  dst, _ := strconv.Atoi(nums[0])
  src, _ := strconv.Atoi(nums[1])
  l, _ := strconv.Atoi(nums[2])
  *m = append(*m, mapping{int64(src), int64(dst), int64(l)})
}

func processMapping(seeds *[]int64, mpg []mapping){
  for i, s := range *seeds {
    for _, m := range mpg {
      if s >= m.src && s <= m.src + m.l - 1 {
	(*seeds)[i] = m.dst + (s - m.src)
	break
      }
    }
  }
  if debug {
    fmt.Printf("%v\n", *seeds)
  }
}

func sortMapping(m *[]mapping) {
  sort.Slice(*m, func(i, j int) bool {
    return (*m)[i].src < (*m)[j].src || ((*m)[i].src == (*m)[j].src && (*m)[i].dst < (*m)[j].dst)
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

  seeds := []int64{}
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
      for _, v := range numS {
	if v != "" {
	  val, err := strconv.Atoi(v)
	  if err != nil {
	    panic(err)
	  }
	  seeds = append(seeds, int64(val))
	}
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

  processMapping(&seeds, seedToSoil)
  processMapping(&seeds, soilToFertilizer)
  processMapping(&seeds, fertilizerToWater)
  processMapping(&seeds, waterToLight)
  processMapping(&seeds, lightToTemperature)
  processMapping(&seeds, temperatureToHumidity)
  processMapping(&seeds, humidityToLocation)

  sort.Slice(seeds, func(i, j int) bool { return seeds[i] < seeds[j] })
  fmt.Printf("%v\n", seeds[0])
}
