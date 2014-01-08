package main

import (
  "os"
  "flag"
  "fmt"
  "io"
  "log"
  "net/http"
  "net/http/cookiejar"
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
  //var start string
  //flag.StringVar(&start, "start", "", "starting url")
  //flag.Parse()
  
  flag.Parse()
  start := flag.Args()[0]
  
  // http://stackoverflow.com/questions/18414212/golang-how-to-follow-location-with-cookie
  options := cookiejar.Options{}
  cookie_jar, err := cookiejar.New(&options)
  client := &http.Client{
    Jar: cookie_jar,
  }
  res, err := client.Get(start)
  if err != nil {
    log.Fatal(err)
  }
  //fmt.Printf("%s\n", start)
  
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
  res.Body.Close()
  if err != nil {
    log.Fatal(err)
  }
}
