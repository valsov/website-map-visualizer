package visualizer

import (
	"encoding/json"
	"os"
	"text/template"

	"github.com/valsov/websitemapper/scraper"
)

type Node struct {
	Id       int    `json:"id"`
	Label    string `json:"label"`
	Internal bool   `json:"internal"`
	Error    bool   `json:"error"`
}

type Edge struct {
	Source      int `json:"from"`
	Destination int `json:"to"`
}

type JsonData struct {
	Nodes string
	Edges string
}

func GenerateVisualizer(pages []*scraper.Page, writePath string) {
	jsonData := toJsonData(pages)
	writeVisualizer(jsonData, writePath)
}

func toJsonData(pages []*scraper.Page) *JsonData {
	data := struct {
		Nodes []Node
		Edges []Edge
	}{Nodes: make([]Node, 0, len(pages))}

	// Generate data structures
	for _, page := range pages {
		data.Nodes = append(data.Nodes, Node{Id: page.Id, Label: page.Url, Internal: page.IsInternalUrl, Error: page.Failed})
		if len(page.OutgoingLinks) != 0 {
			for _, linkId := range page.OutgoingLinks {
				data.Edges = append(data.Edges, Edge{Source: page.Id, Destination: linkId})
			}
		}
	}

	// To JSON
	nodesBytes, err := json.Marshal(data.Nodes)
	if err != nil {
		panic(err)
	}
	edgesBytes, err := json.Marshal(data.Edges)
	if err != nil {
		panic(err)
	}
	return &JsonData{
		Nodes: string(nodesBytes),
		Edges: string(edgesBytes),
	}
}

func writeVisualizer(data *JsonData, writePath string) {
	var tmplFile = "view.tmpl"
	tmpl, err := template.New(tmplFile).ParseFiles(tmplFile)
	if err != nil {
		panic(err)
	}

	fileHandle, err := os.Create(writePath)
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(fileHandle, data)
	if err != nil {
		panic(err)
	}
}
