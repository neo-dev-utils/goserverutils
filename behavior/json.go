package behavior

// JSONBody json对象
type JSONBody struct {
	Judge []struct {
		Name     string `json:"name,omitempty"`
		Describe string `json:"describe,omitempty"`
	} `json:"judge,omitempty"`
	Action []struct {
		Name     string `json:"name,omitempty"`
		Describe string `json:"describe,omitempty"`
	} `json:"action,omitempty"`
	Behavior []JSONNode `json:"behavior,omitempty"`
}

// JSONNode json节点
type JSONNode struct {
	Kind                      NodeType   `json:"kind,omitempty"`
	IsAction                  bool       `json:"isAction,omitempty"`
	JudeCondition             string     `json:"judeCondition,omitempty"`
	ProbabilityStartCondition uint32     `json:"probabilityStartCondition,omitempty"`
	ProbabilityEndCondition   uint32     `json:"probabilityEndCondition,omitempty"`
	Children                  []JSONNode `json:"children,omitempty"`
	Label                     string     `json:"label,omitempty"`
	Action                    string     `json:"action,omitempty"`
}
