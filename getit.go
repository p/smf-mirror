package main

import (
  "fmt"
  "io/ioutil"
  "log"
  "net/http"
)

func main() {
  res, err := http.Get("http://www.panynj.gov/path/full-schedules.html")
  if err != nil {
    log.Fatal(err)
  }
  robots, err := ioutil.ReadAll(res.Body)
  res.Body.Close()
  if err != nil {
    log.Fatal(err)
  }
  fmt.Printf("%s", robots)
}
