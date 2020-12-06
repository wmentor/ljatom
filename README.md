# Summary

ljatom reads new public posts from LiveJournal via public atom stream interface.

# Install

```
go get github.com/wmentor/ljatom
```

# Usage

All you need is:

```go
package main

import (
  "fmt"

  "github.com/wmentor/ljatom"
)

func main() {

  for msg := range ljatom.Read() {

    fmt.Println("Time: " + msg.Created.String())
    fmt.Println("Journal: " + msg.Journal)
    fmt.Println("Post URL: " + msg.Url)
    fmt.Println("Post title: " + msg.Title)
    fmt.Print("Post body: " + msg.Content + "\n\n\n")

  }
}
```

and wait some times. The atom stream returns data by chunks.

Application print result like this:

```
Time: 2019-11-21 20:57:07 +0000 UTC
Journal: faded_fae
Post URL: https://faded-fae.livejournal.com/967096.html
Post title: DuckTales (2017)
Post body: <p>Canon or AU: AUx3<br>Fic: Kings and Queens</p>
<p>A/N: I am working on other fics, I promise you. XD I’m also trying to balance that
with the new Pokemon game and reading my book, which is awesome.&nbsp;</p>
<p>Oh, on a more serious note--Lena won’t be in this. Otherwise, Webby would be
distracted by her. So, uh, in this universe, I guess Lena doesn’t exist? Sorry, Lena. ^^;</p>
<p>------</p>
<a name="cutid1"></a>
<p>It had taken Wren...

```

*ljatom.Read* return chan of *Entry* refs.  Each *Entry* object has following structure:

```go
type Entry struct {
  Journal      string
  JournalTitle string
  Url          string
  Created      time.Time
  Title        string
  Content      string
}
```

Content is HTML (maybe invalid).

Moreover, *ljatom.Read* makes reconnect if connection is broken. All you need is read chan of Entry objects.
