package main

import (
  "fmt"
  "io/ioutil"
  "log"
  "net/http"
  "strings"
  "path"
  "code.google.com/p/go.net/html"
)

func main() {
  start := "http://www.panynj.gov/path/full-schedules.html"
  res, err := http.Get(start)
  if err != nil {
    log.Fatal(err)
  }
  z := html.NewTokenizer(res.Body)
  for {
    tt := z.Next()
    if tt == html.ErrorToken {
      break
    }
    if tt == html.StartTagToken {
      tn, _ := z.TagName()
      if len(tn) == 1 && tn[0] == 'a' {
        for {
          key, value, more := z.TagAttr()
          // http://stackoverflow.com/questions/14230145/what-is-the-best-way-to-convert-byte-array-to-string
          if string(key) == "href" {
            v := string(value)
            if strings.HasPrefix(v, "schedules/") {
              url := path.Join(path.Dir(start), v)
              fmt.Printf("%s\n", url)
            }
          }
          if !more {
            break
          }
        }
      }
      // ...
      //return ...
    }
    // Process the current token.
  }
  robots, err := ioutil.ReadAll(res.Body)
  res.Body.Close()
  if err != nil {
    log.Fatal(err)
  }
  fmt.Printf("%s", robots)
}
