package dynamic

//easyjson:json
type tracker struct {
	Clicks      []string `json:"clicks,omitempty"`
	Impressions []string `json:"impressions"`
	Views       []string `json:"views"`
}

//easyjson:json
type asset struct {
	Name string `json:"name"`
	Path string `json:"path"`
	Type string `json:"type"`
}

//easyjson:json
type item struct {
	ID         interface{}            `json:"id"`
	Type       string                 `json:"type"`
	URL        string                 `json:"url,omitempty"`
	Content    string                 `json:"content,omitempty"`
	ContentURL string                 `json:"content_url,omitempty"`
	Fields     map[string]interface{} `json:"fields,omitempty"`
	Assets     []asset                `json:"assets,omitempty"`
	Tracker    tracker                `json:"tracker"`
	Debug      interface{}            `json:"debug,omitempty"`
}

//easyjson:json
type group struct {
	ID    string  `json:"id"`
	Items []*item `json:"items"`
}

// Response object description
//easyjson:json
type Response struct {
	Version string   `json:"version"`
	Groups  []*group `json:"groups,omitempty"`
}

func (r *Response) getGroupOrCreate(groupID string) *group {
	for _, g := range r.Groups {
		if g.ID == groupID {
			return g
		}
	}
	g := &group{ID: groupID}
	r.Groups = append(r.Groups, g)
	return g
}
