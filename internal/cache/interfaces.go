package cache

type Stats interface {
	Type() string
}

type CacheManager interface {
	GetStats() (Stats, error)
	Clear() (int, int64, error)
	GetLocation() string
}
