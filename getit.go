package main

import (
  "os"
  "fmt"
  "io"
  "io/ioutil"
  "log"
  "net/http"
  "strings"
  "path"
  "code.google.com/p/go.net/html"
)

func fetch(url string) {
  res, err := http.Get(url)
  if err != nil {
    log.Fatal(err)
  }
  if err != nil {
    log.Fatal(err)
  }
  // http://stackoverflow.com/questions/1821811/how-to-read-write-from-to-file
  fo, err := os.Create(path.Base(url))
  if err != nil {
    panic(err)
  }
  defer func() {
    if err := fo.Close(); err != nil {
      panic(err)
    }
  }()
  buf := make([]byte, 65536)
  for {
    n, err := res.Body.Read(buf)
    if err != nil && err != io.EOF {
      panic(err)
    }
    if n == 0 {
      break
    }
    
    if _, err := fo.Write(buf[:n]); err != nil {
      panic(err)
    }
  }
  res.Body.Close()
}

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
              fuckedurl := path.Join(path.Dir(start), v)
              // yep, hack it
              // thx go for making me rename the variable
              url := strings.Replace(fuckedurl, ":/", "://", 1)
              fmt.Printf("%s\n", url)
              fetch(url)
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
