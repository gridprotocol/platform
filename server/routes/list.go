package routes

import (
	"encoding/json"
	"fmt"
	"os"
)

// write some nodes into cofig file
func NewList() {

	nodes := []NodeInfo{
		{
			Name:     "node1",
			Entrance: "empty",
			Resource: "empty",
			Price:    "empty",
		},
		{
			Name:     "node2",
			Entrance: "empty",
			Resource: "empty",
			Price:    "empty",
		},
	}

	// open/create file
	file, err := os.Create("nodes.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// objects to json bytes
	jsonData, err := json.MarshalIndent(&nodes, "", "  ")
	if err != nil {
		panic(err)
	}

	// write data
	_, err = file.WriteString(string(jsonData))
	if err != nil {
		panic(err)
	}

	fmt.Println("write file ok.")

}

// read node list from config file
func ReadList() []NodeInfo {
	var nodes []NodeInfo

	// read node list from config file
	data, err := os.ReadFile("nodes.json")
	if err != nil {
		panic(err)
	}

	// decode data into objects
	err = json.Unmarshal(data, &nodes)
	if err != nil {
		panic(err)
	}

	return nodes
}
