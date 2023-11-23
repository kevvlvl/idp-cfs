package request

type Contract struct {
	Git struct {
		Tool   string `yaml:"tool"`
		Org    string `yaml:"org"`
		Repo   string `yaml:"repo"`
		Branch string `yaml:"branch"`
	} `yaml:"git"`

	GoldenPath struct {
		Name    string `yaml:"name"`
		Version string `yaml:"version"`
	} `yaml:"golden-path"`

	Deployment struct {
		Kubernetes struct {
			ClusterUrl string `yaml:"cluster-url"`
			Namespace  string `yaml:"namespace"`
			Logs       bool   `yaml:"logs"`
		} `yaml:"kubernetes"`
	} `yaml:"deployment"`
}
