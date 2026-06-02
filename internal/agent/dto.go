package agent

type PromptTypeResolving struct {
	NotProject     bool `json:"not_project"`
	DiscussProject bool `json:"discuss_project"`
	EditProject    bool `json:"edit_project"`

	WantedFiles []string `json:"wanted_files"`
}
