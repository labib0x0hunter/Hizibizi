package main

type Storage interface {
	Set(key string, value string)
	Get(key string) (string, bool)
	Delete(key string)
}

func Set(s Storage, key string, val string) {
	s.Set(key, val)
}

func Get(s Storage, key string) (string, bool) {
	return s.Get(key)
}

func Delete(s Storage, key string) {
	s.Delete(key)
}
