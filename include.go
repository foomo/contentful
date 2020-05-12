package contentful

//IncludeEntry model
type IncludeEntry struct {
	Fields map[string]interface{} `json:"fields,omitempty"`
	Sys    *Sys                   `json:"sys"`
}

// IncludeFileFields model
type IncludeFileFields struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	File        *File  `json:"file,omitempty"`
}

// IncludeAsset model
type IncludeAsset struct {
	Fields *IncludeFileFields `json:"fields"`
	Sys    *Sys               `json:"sys"`
}
