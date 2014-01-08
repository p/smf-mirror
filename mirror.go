package main

import (
  "encoding/json"
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
  "github.com/steveyen/gkvlite"
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

const (
  _ = iota
  linkfound
  linkcrawled
)

func findBoardLinks(client *http.Client, start string, board_links map[string]int) (board_links_out map[string]int) {
  res, err := client.Get(start)
  if err != nil {
    log.Fatal(err)
  }
  //fmt.Printf("%s\n", start)
  
  c := FindLinks(res.Body)
  for {
    v := <- c
    if v == "" {
      break
    }
    //fmt.Printf("%v\n", v)
    if strings.Contains(v, start) && strings.Contains(v, "?board=") {
      // index.php?board=18.0;sort=starter -- drop ;.*
      if strings.Contains(v, ";") {
        pre, _ := Split2(v, ";")
        v = pre
      }
      //fmt.Printf("recursing\n")
      _, found := board_links[v]
      if !found {
        board_links[v] = linkfound
      }
    }
  }
  res.Body.Close()
  return board_links
}

func loadBoardLinks(client *http.Client, start string) (flat_links []string) {
  board_links := map[string]int{}
  board_links[start] = linkfound
  for {
    copy := board_links
    any := false
    for key, value := range copy {
      if value == linkcrawled {
        continue
      }
      fmt.Printf("%v\n", key)
      board_links = findBoardLinks(client, key, board_links)
      board_links[key] = linkcrawled
      any = true
    }
    if !any {
      break
    }
  }
  
  flat_links = []string{}
  for key, _ := range board_links { 
    flat_links = append(flat_links, key)
  }
  return flat_links
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
  
  f, err := os.OpenFile("smfmirror.gkvlite", os.O_RDWR, 0666)
  if err != nil {
    f, err = os.Create("smfmirror.gkvlite")
  }
  if err != nil {
    panic(err)
  }
  s, err := gkvlite.NewStore(f)
  if err != nil {
    panic(err)
  }
  c := s.GetCollection("smfmirror")
  if c == nil {
    c = s.SetCollection("smfmirror", nil)
  }
  pboard_links, err := c.Get([]byte("board_links"))
  var board_links []string
  if pboard_links != nil {
    err := json.Unmarshal(pboard_links, &board_links)
    if err != nil {
      panic(err)
    }
  } else {
    board_links = loadBoardLinks(client, start)
    b, err := json.Marshal(board_links)
    if err != nil {
      panic(err)
    }
    //c.Set([]byte("hello"), []byte("world"))
    //a, x := c.Get([]byte("hello"))
    //fmt.Printf("%v %v\n", a, x)
    c.Set([]byte("board_links"), b)
    s.Flush()
    s.Close()
    f.Sync()
    f.Close()
  }
  fmt.Printf("%v\n", board_links)
}
