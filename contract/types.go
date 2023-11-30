package contract

type Contract struct {
	Action string `yaml:"action"`

	Code struct {
		Tool   string  `yaml:"tool"`
		Org    *string `yaml:"org,omitempty"`
		Repo   string  `yaml:"repo"`
		Branch string  `yaml:"branch"`
	} `yaml:"code"`

	GoldenPath struct {
		Git    string `yaml:"git"`
		Name   string `yaml:"name"`
		Branch string `yaml:"branch"`
		Tag    string `yaml:"tag"`
		Path   string `yaml:"path"`
	} `yaml:"golden-path"`

	Deployment struct {
		Kubernetes struct {
			ClusterUrl string `yaml:"cluster-url"`
			Namespace  string `yaml:"namespace"`
			Logs       bool   `yaml:"logs"`
		} `yaml:"kubernetes"`
	} `yaml:"deployment"`
}
