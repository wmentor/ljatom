package ljatom

import (
	"bufio"
	"bytes"
	"encoding/xml"
	"net"
	"strings"
	"time"

	"github.com/wmentor/epoch"
	"github.com/wmentor/log"
)

type Entry struct {
	Journal      string
	JournalTitle string
	Url          string
	Created      time.Time
	Title        string
	Content      string
}

func Read() <-chan *Entry {

	out := make(chan *Entry, 256)

	go func() {
		for {
			tact(out)
		}
	}()

	return out
}

type feed struct {
	Author struct {
		Journal string `xml:"journal"`
		Title   string `xml:"name"`
	} `xml:"author"`
	Entry struct {
		Title     string `xml:"title"`
		Published string `xml:"published"`
		Link      struct {
			Href string `xml:"href,attr"`
		} `xml:"link"`
		Content string `xml:"content"`
	} `xml:"entry"`
}

func tact(out chan *Entry) {

	defer func() {

		if r := recover(); r != nil {
			log.Error("ljatom: connect broken")
		}

	}()

	con, err := net.DialTimeout("tcp", "atom.services.livejournal.com:80", time.Hour)
	if err != nil {
		panic(err)
	}
	defer con.Close()

	msg := `GET /atom-stream.xml HTTP/1.1
Host: atom.services.livejournal.com
User-Agent: ljatom


`

	con.Write([]byte(msg))

	br := bufio.NewReaderSize(con, 1024*1024)

	start := false
	builder := bytes.NewBuffer(nil)

	for {

		str, err := br.ReadString('\n')
		if err != nil && str == "" {
			break
		}

		if !start {

			if strings.Index(str, "<feed ") != -1 {
				start = true
				builder.Reset()
				builder.WriteString("<feed>\n")
			}

		} else {

			if strings.Index(str, "</feed>") != -1 {
				builder.WriteString("</feed>")
				start = false

				f := new(feed)

				xml := xml.NewDecoder(builder)

				if err = xml.Decode(f); err != nil {
					continue
				} else {
					e := &Entry{
						Journal:      f.Author.Journal,
						JournalTitle: f.Author.Title,
						Url:          f.Entry.Link.Href,
						Title:        f.Entry.Title,
						Content:      f.Entry.Content,
					}

					if tm, err := epoch.Parse(f.Entry.Published); err == nil {
						e.Created = tm
					} else {
						e.Created = time.Now()
					}

					out <- e
				}

			} else {
				str = strings.ReplaceAll(str, "<lj:", "<")
				str = strings.ReplaceAll(str, "</lj:", "</")
				str = strings.ReplaceAll(str, "'", "\"")
				builder.WriteString(str)
			}

		}
	}

	time.Sleep(time.Second)
}
