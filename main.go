package main

import (
  "fmt"
  "net/http"
  "os"
  "bufio"
  "strings"
  "golang.org/x/net/html"
)

func getHref(t html.Token) (ok bool, href string) {
  for _, a := range t.Attr {
    if a.Key == "href" {
      href = a.Val
      ok = true
    }
  }
  // bare return will return the variables sent into the fn (ok, href)
  return
}

func getWords(url string, ch chan string, chFinished chan bool) {
  resp, err := http.Get(url)

  defer func() {
    chFinished <- true
  }()

  if err != nil {
    fmt.Println("Error: Failed to crawl \"" + url + "\"")
    return
  }

  b := resp.Body
  defer b.Close()

  bodyTokens := html.NewTokenizer(b)

  for {
    currentToken := bodyTokens.Next()

    switch {
    case currentToken == html.ErrorToken:
      // end of document
      return

    case currentToken == html.TextToken:
      t := strings.TrimSpace(string(bodyTokens.Text()))
      if t == "" {
        continue
      }
      words := strings.Split(t, " ")
      for _, word := range words {
        if len(word) > 2 {
          ch <- word
        }
      }
      continue
    }
  }
}

func getLinks(url string, ch chan string, chFinished chan bool) {
  resp, err := http.Get(url)

  defer func() {
    chFinished <- true
  }()

  if err != nil {
    fmt.Println("Error: Failed to crawl \"" + url + "\"")
    return
  }

  b := resp.Body
  defer b.Close()

  bodyTokens := html.NewTokenizer(b)

  for {
    currentToken := bodyTokens.Next()

    switch {
    case currentToken == html.ErrorToken:
      // end of document
      return

    case currentToken == html.StartTagToken:
      t := bodyTokens.Token()
      isAnchor := t.Data == "a"
      if !isAnchor {
        continue
      }

      ok, tokenUrl := getHref(t)
      if !ok {
        continue
      }

      hasProtocol := strings.Index(tokenUrl, "http") == 0
      isRootPath := strings.Index(strings.TrimSpace(tokenUrl), "/") == 0
      if hasProtocol {
        ch <- tokenUrl
      } else if isRootPath {
        ch <- url + strings.TrimSpace(tokenUrl)
      }
      continue
    }
  }
}

func writeFile(wordlist []string, path string) error {
  file, err := os.Create(path)
  if err != nil {
    return err
  }
  defer file.Close()

  w := bufio.NewWriter(file)
  for _, word := range wordlist {
    fmt.Fprintln(w, word)
  }
  return w.Flush()
}

func main() {
  foundUrls := make(map[string]bool)
  seedUrls := os.Args[1:]


  chUrls := make(chan string)
  chFinished := make(chan bool)
  for _, url := range seedUrls {
    go getLinks(url, chUrls, chFinished)
  }

  for c := 0; c < len(seedUrls); {
    select {
    case url := <-chUrls:
      foundUrls[url] = true
    case <-chFinished:
      c++
    }
  }
  for _, seedUrl := range(seedUrls) {
    foundUrls[seedUrl] = true
  }

  wordlist := []string{}
  chWords := make(chan string)
  chWordsFinished := make(chan bool)
  for foundUrl := range foundUrls {
    go getWords(foundUrl, chWords, chWordsFinished)
  }
  for i := 0; i < len(foundUrls); {
    select {
    case word := <-chWords:
      wordlist = append(wordlist, word)
    case <-chWordsFinished:
      i++
    }
  }
  writeFile(wordlist, "words.txt")
  close(chUrls)
  close(chWords)
}
