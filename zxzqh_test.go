package zxzqh

import (
	"encoding/json"
	"os"
	"testing"
)

var (
	generate     = "./generate"
	nodeListFile = generate + "/node_list.json"
	nodeTreeFile = generate + "/node_tree.json"
)

func TestNodeList(t *testing.T) {
	b, err := json.MarshalIndent(NodeList(), "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(nodeListFile, b, 0600)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNodeTree(t *testing.T) {
	b, err := json.MarshalIndent(NodeTree(), "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(nodeTreeFile, b, 0600)
	if err != nil {
		t.Fatal(err)
	}
}

func TestCodeNode(t *testing.T) {
	t.Log(CodeNode(110000))
}
