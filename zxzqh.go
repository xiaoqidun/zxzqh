package zxzqh

import (
	"bytes"
	_ "embed"
	"github.com/PuerkitoBio/goquery"
	"strconv"
	"strings"
)

//go:embed assets/20201201.html
var data []byte

var (
	err         error
	nodeList    []Node
	nodeTree    Tree
	codeNodeMap = make(map[int]Node)
)

type Node struct {
	Code   int    `json:"code"`
	Name   string `json:"name"`
	Parent int    `json:"parent"`
}

type Tree struct {
	Node
	Children []Tree `json:"children,omitempty"`
}

func init() {
	nodeList, err = generateNodeList(data)
	if err != nil {
		panic(err)
	}
	nodePMap := make(map[int][]Node)
	for i := 0; i < len(nodeList); i++ {
		nodePMap[nodeList[i].Parent] = append(nodePMap[nodeList[i].Parent], nodeList[i])
		codeNodeMap[nodeList[i].Code] = nodeList[i]
	}
	nodeTree = generateNodeTree(nodePMap, 0)
}

// NodeList 获取中华人民共和国行政区划代码扁平列表
func NodeList() []Node {
	return nodeList
}

// NodeTree 获取中华人民共和国行政区划代码树状列表
func NodeTree() Tree {
	return nodeTree
}

// CodeNode 通过中华人民共和国行政区划代码查找节点
func CodeNode(code int) *Node {
	v, ok := codeNodeMap[code]
	if !ok {
		return nil
	}
	return &v
}

func generateNodeList(html []byte) (nodes []Node, err error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(html))
	if err != nil {
		return
	}
	top, parent := 0, 0
	doc.Find(`tr[height$="19"]`).Each(func(i int, selection *goquery.Selection) {
		id := 0
		name := ""
		selection.Find(".xl7228320").Each(func(i int, selection *goquery.Selection) {
			switch i {
			case 0:
				id, _ = strconv.Atoi(selection.Text())
			case 1:
				name = selection.Text()
			}
		})
		if id != 0 && name != "" {
			parent = id
			if !strings.HasPrefix(name, " ") {
				nodes = append(nodes, Node{
					Code:   id,
					Parent: 0,
					Name:   name,
				})
				top = id
			} else {
				nodes = append(nodes, Node{
					Code:   id,
					Parent: top,
					Name:   strings.TrimSpace(name),
				})
			}
			return
		}
		selection.Find(".xl7328320").Each(func(i int, selection *goquery.Selection) {
			switch i {
			case 0:
				id, _ = strconv.Atoi(selection.Text())
			case 1:
				name = strings.TrimSpace(selection.Text())
			}
		})
		if id != 0 && name != "" {
			nodes = append(nodes, Node{
				Code:   id,
				Parent: parent,
				Name:   name,
			})
		}
	})
	return
}

func generateNodeTree(data map[int][]Node, root int) Tree {
	tree := Tree{
		Node: Node{
			Code:   0,
			Parent: 0,
			Name:   "根节点",
		},
		Children: nil,
	}
	if v, ok := data[root]; ok {
		for i := 0; i < len(v); i++ {
			item := Tree{
				Node:     v[i],
				Children: nil,
			}
			if _, ok2 := data[v[i].Code]; ok2 {
				item.Children = append(item.Children, generateNodeTree(data, v[i].Code).Children...)
			}
			tree.Children = append(tree.Children, item)
		}
	}
	return tree
}
