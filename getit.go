package main

import (
  "os"
  "flag"
  "fmt"
  "io"
  "log"
  "net/http"
  "net/http/cookiejar"
  "net/url"
  "strings"
  "path"
  "code.google.com/p/go.net/html"
)

func Split2(str, sep string) (string, string) {
  s := strings.Split(str, sep)
  return s[0], s[1]
}

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

func Fuckoff(pres *http.Response) {
}

func FindLinks(body io.Reader) chan string {
  c := make(chan string)
  
  go func() {
    z := html.NewTokenizer(body)
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
              // http://codereview.stackexchange.com/questions/28386/fibonacci-generator-with-golang
              c <- v
            }
            if !more {
              break
            }
          }
        }
      }
    }
    c <- ""
  }()
  
  return c
}

func main() {
  //var start string
  //flag.StringVar(&start, "start", "", "starting url")
  //flag.Parse()
  
  flag.Parse()
  start := flag.Args()[0]
  username, password := Split2(flag.Args()[1], ":")
  
  // http://stackoverflow.com/questions/18414212/golang-how-to-follow-location-with-cookie
  options := cookiejar.Options{}
  cookie_jar, err := cookiejar.New(&options)
  client := &http.Client{
    Jar: cookie_jar,
  }
  pres, err := client.PostForm(start + "index.php?action=login2",
    url.Values{"user": {username}, "passwrd": {password}})
  if err != nil {
    log.Fatal(err)
  }
  Fuckoff(pres)
  
  res, err := client.Get(start)
  if err != nil {
    log.Fatal(err)
  }
  //fmt.Printf("%s\n", start)
  
  board_links := map[string]bool{}
  crawled_board_links := map[string]bool{}
  c := FindLinks(res.Body)
  for {
    v := <- c
    if v == "" {
      break
    }
    fmt.Printf("%v\n", v)
    if strings.Contains(v, start) && strings.Contains(v, "?board=") {
      fmt.Printf("recursing\n")
      _, found := crawled_board_links[v]
      if !found {
        board_links[v] = true
      }
    }
  }
  res.Body.Close()
  if err != nil {
    log.Fatal(err)
  }
}
