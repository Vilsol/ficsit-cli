package ficsit

type AllVersionsResponse struct {
	Error   *Error       `json:"error,omitempty"`
	Data    []ModVersion `json:"data,omitempty"`
	Success bool         `json:"success"`
}

type ModVersion struct {
	ID           string       `json:"id"`
	Version      string       `json:"version"`
	Dependencies []Dependency `json:"dependencies"`
	Targets      []Target     `json:"targets"`
}

type Dependency struct {
	ModID     string `json:"mod_id"`
	Condition string `json:"condition"`
	Optional  bool   `json:"optional"`
}

type Target struct {
	VersionID  string `json:"version_id"`
	TargetName string `json:"target_name"`
	Hash       string `json:"hash"`
	Size       int64  `json:"size"`
}

type Error struct {
	Message string `json:"message"`
	Code    int64  `json:"code"`
}
