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
		Url    *string `yaml:"url,omitempty"`
		Name   *string `yaml:"name,omitempty"`
		Path   *string `yaml:"path,omitempty"`
		Branch *string `yaml:"branch,omitempty"`
		Tag    *string `yaml:"tag,omitempty"`
	} `yaml:"golden-path"`

	Deployment struct {
		Kubernetes struct {
			ClusterUrl string `yaml:"cluster-url"`
			Namespace  string `yaml:"namespace"`
			Logs       bool   `yaml:"logs"`
		} `yaml:"kubernetes"`
	} `yaml:"deployment"`
}
