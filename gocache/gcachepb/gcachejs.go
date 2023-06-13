package gcachepb

type GetRequest struct {
	Group string `json:"group"`
	Key   string `json:"key"`
}

type GetResponse struct {
	Value string `json:"value"`
}
