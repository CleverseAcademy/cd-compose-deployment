package dto

type DeployImageDto struct {
	Priority uint64 `json:"p"`
	Ref      string `json:"r"`
	Service  string `json:"s"`
	Image    string `json:"i"`
}
