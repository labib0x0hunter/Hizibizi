package main

type MapKvs struct {
	kvs map[string]string
}

func NewMapKvs() *MapKvs {
	return &MapKvs{
		kvs: make(map[string]string),
	}
}

func (km *MapKvs) Get(key string) (string, bool) {
	val, ok := km.kvs[key]
	return val, ok
}

func (km *MapKvs) Set(key, val string) {
	km.kvs[key] = val
}

func (km *MapKvs) Delete(key string) {
	delete(km.kvs, key)
}
