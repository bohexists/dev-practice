package main

import (
	"fmt"
	"github.com/bohexists/cache-library/cache"
	"time"
)

func main() {
	// Инициализация кеша с TTL 1 минуту
	c := cache.New(cache.CacheConfig{MaxSize: 100, DefaultTTL: time.Minute})

	key := "testKey"
	value := "testValue"

	err := c.Set(key, value)
	if err != nil {
		fmt.Println("Error setting value:", err)
		return
	} else {
		fmt.Println("Value set successfully")
	}

	// Попробуйте получить значение сразу после его установки
	item, err := c.Get(key)
	if err != nil {
		fmt.Println("Error getting value:", err)
		return
	}

	if item != nil {
		fmt.Printf("Got value: %v\n", item)
	} else {
		fmt.Println("Got value: <nil>")
	}

	// Подождите, чтобы проверить TTL (например, 2 минуты)
	time.Sleep(2 * time.Minute)

	// Попробуйте получить значение после истечения TTL
	item, err = c.Get(key)
	if err != nil {
		fmt.Println("Error getting value after TTL:", err)
		return
	}

	if item != nil {
		fmt.Printf("Got value after TTL: %v\n", item)
	} else {
		fmt.Println("Got value after TTL: <nil>")
	}
}
