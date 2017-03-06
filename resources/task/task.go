package task

type ApplicationJSON struct {
	Name        string              `json:"name"`
	Resources   *ResourceJSON       `json:"resources"`
	Command     *CommandJSON        `json:"command"`
	Container   *ContainerJSON      `json:"container"`
	HealthCheck *HealthCheckJSON    `json:"healthcheck"`
	Labels      []map[string]string `json:"labels"`
}

type HealthCheckJSON struct {
	Endpoint *string `json:"endpoint"`
}

type KillJson struct {
	Name *string `json:"name"`
}

type ResourceJSON struct {
	Mem  float64 `json:"mem"`
	Cpu  float64 `json:"cpu"`
	Disk float64 `json:"disk"`
}

type CommandJSON struct {
	Cmd  *string   `json:"cmd"`
	Uris []UriJSON `json:"uris"`
}

type ContainerJSON struct {
	ImageName *string       `json:"image"`
	Tag       *string       `json:"tag"`
	Network   []NetworkJSON `json:"network"`
}

type NetworkJSON struct {
	IpAddresses []IpAddressJSON `json:"ipaddress,omitempty"`
	Name        *string         `json:"name"`
}

type IpAddressJSON struct {
	IP       *string `json:"ip"`
	Protocol *string `json:"protocol"`
}

type UriJSON struct {
	Uri     *string `json:"uri"`
	Extract *bool   `json:"extract"`
	Execute *bool   `json:"execute"`
}