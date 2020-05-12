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

// IncludeFileLocalizedFields model
type IncludeFileLocalizedFields struct {
	Title       map[string]string `json:"title,omitempty"`
	Description map[string]string `json:"description,omitempty"`
	File        map[string]*File  `json:"file,omitempty"`
}

// IncludeLocalizedAsset model
type IncludeLocalizedAsset struct {
	Fields *IncludeFileLocalizedFields `json:"fields"`
	Sys    *Sys                        `json:"sys"`
}
