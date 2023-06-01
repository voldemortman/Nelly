package furtado

type IStalkerConfig interface {
	RemoteStalkerConfig | LocalStalkerConfig
}
