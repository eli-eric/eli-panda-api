package models

// type GraphNode = {
// 	uid: string
// 	name: string
// 	properties: Record<string,string>
//   }

//   type GraphLink = {
// 	target: string,
// 	relationship: string,
// 	source: string
//   }

//   type SystemGraphResponse = {
// 	nodes: GraphNode[]
// 	links: GraphLink[]
//   }

type GraphNode struct {
	Uid        string            `json:"uid"`
	Name       string            `json:"name"`
	Label      string            `json:"label"`
	Properties map[string]string `json:"properties"`
}

type GraphLink struct {
	Target       string `json:"target"`
	Relationship string `json:"relationship"`
	Source       string `json:"source"`
}

type GraphResponse struct {
	Nodes []GraphNode `json:"nodes"`
	Links []GraphLink `json:"links"`
}

type GlobalSearchResult struct {
	Uid         string `json:"uid"`
	Name        string `json:"name"`
	Description string `json:"description"`
	NodeType    string `json:"nodeType"`
}

type GlobalSearchResponse struct {
	TotalCount int64                `json:"totalCount"`
	Data       []GlobalSearchResult `json:"data"`
}
