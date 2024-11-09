package types

type container struct {
	language string
	name     string
}

type ListOfContainers struct {
}

type Cases struct {
	Input    any `json:"input"`
	Expected any `json:"expected"`
}

type Body struct {
	Language string
}
