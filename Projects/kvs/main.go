package main

import "fmt"

func main() {

	mapKvs := NewMapKvs()
	redisKvs := NewRedisKvs()

	Set(mapKvs, "name", "labib from map")
	Set(redisKvs, "name", "faisal from redis")

	if val, ok := Get(mapKvs, "name"); ok {
		fmt.Println("MapKvs: ", val)
		Delete(mapKvs, "name")
	}

	if val, ok := Get(redisKvs, "name"); ok {
		fmt.Println("RedisKvs: ", val)
		Delete(redisKvs, "name")
	}
}