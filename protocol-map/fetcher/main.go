package main

import (
	"os"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"strings"
	"strconv"
	"encoding/json"
)

type Map struct {
	Clientbound map[string]string `json:"clientbound"`
	Serverbound map[string]string `json:"serverbound"`
}

type Content struct {
	Name string `json:"name"`
	Protocol int `json:"protocol"`
	Base int `json:"base"`
	Map Map `json:"map"`
}

type Type struct {
	Name string `json:"name"`
	Version int `json:"version"`
}

type Module struct {
	Type Type `json:"type"`
	Content Content `json:"content"`
}

func main() {
	args := os.Args[1:]

	if len(args) < 1 {
		fmt.Println("You have to specify a link")
		os.Exit(0)
	}

	doc, err := goquery.NewDocument(args[0])
	if err != nil {
		panic(err)
	}

	content := Content{}
	content.Map.Clientbound = make(map[string]string)
	content.Map.Serverbound = make(map[string]string)

	doc.Find(`a[title="Protocol version numbers"]`).Each(func(i int, s *goquery.Selection) {
		html, err := s.Html()
		if err != nil {
			panic(err)
		}

		splt := strings.Split(string(html), ", protocol ")

		if len(splt) > 1 {
			if content.Base != 0 {
				content.Name = splt[0]
				str, err := strconv.Atoi(splt[1])
				if err != nil {
					panic(err)
				}
				content.Protocol = str
			} else {
				str, err := strconv.Atoi(splt[1])
				if err != nil {
					panic(err)
				}
				content.Base = str
			}
		}
	})

	clientbound := false
	serverbound := false
	doc.Find(`h3 > #Packets`).Parent().Next().Find(".wikitable tbody tr").Each(func(i int, s *goquery.Selection) {
		del := s.Find("td del")
		ins := s.Find("td ins")
		if del.Size() > 0 && ins.Size() > 0 {
			htmlDel, err := del.Html()
			if err != nil {
				panic(err)
			}

			htmlIns, err := ins.Html()
			if err != nil {
				panic(err)
			}

			if clientbound {
				content.Map.Clientbound[htmlDel] = htmlIns
			} else if serverbound {
				content.Map.Serverbound[htmlDel] = htmlIns
			}
		} else {
			head := s.Find("th")

			html, err := head.Html()
			if err != nil {
				panic(err)
			}

			if strings.Contains(html,"Play clientbound") {
				clientbound = true
			} else if strings.Contains(html,"Play serverbound") {
				serverbound = true
				clientbound = false
			}
		}
	})

	module := Module{
		Type{"protocol-map", 1},
		content,
	}

	str, err := json.MarshalIndent(module, "", "    ")
	if err != nil {
		panic(err)
	}

	path := module.Content.Name + ".json"
	_, err = os.Create(path)
	if err != nil {
		panic(err)
	}

	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		panic(err)
	}

	_, err = file.Write(str)
	if err != nil {
		panic(err)
	}
}
