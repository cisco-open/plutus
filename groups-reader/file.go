package groupsreader

type githubFile struct {
	Type        string `json:"type"`
	Encoding    string `json:"encoding"`
	Size        int    `json:"size"`
	Name        string `json:"name"`
	Path        string `json:"path"`
	Content     string `json:"content"`
	Sha         string `json:"sha"`
	URL         string `json:"url"`
	GitURL      string `json:"git_url"`
	HTMLURL     string `json:"html_url"`
	DownloadURL string `json:"download_url"`
	Links       struct {
		Git  string `json:"git"`
		Self string `json:"self"`
		HTML string `json:"html"`
	} `json:"_links"`
}

type yamlFile struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Spec        struct {
		Type       string        `yaml:"type"`
		Reason     string        `yaml:"reason"`
		Attributes []interface{} `yaml:"attributes"`
		Owners     []string      `yaml:"owners"`
		Members    []string      `yaml:"members"`
	} `yaml:"spec"`
}

func (yf yamlFile) Members() []string {
	return yf.Spec.Members
}
