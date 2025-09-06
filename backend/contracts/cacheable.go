package contracts


type Cacheable interface {
	CacheKey() string
}
