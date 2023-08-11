package cache

import (
    "fmt"
    "log"
)

var cache *CacheRepository

func Exists(key, value string) bool {
    initCache()
    item, err := cache.GetItem(key)
    if err != nil {
        fmt.Printf("Cache Exists error: %v\n", err)
        return false
    }
    return item != nil && item.Value == value
}

func Update(key, value string) error {
    initCache()
    item, err := cache.GetItem(key)
    if err != nil {
        fmt.Printf("GetItem error: %v\n", err)
        return err
    }
    
    if item != nil {
        err = cache.DeleteItem(key)
        if err != nil {
            fmt.Printf("DeleteItem error: %v\n", err)
            return err
        }
    }
    
    err = cache.SaveItem(key, value)
    if err != nil {
        fmt.Printf("SaveItem error: %v\n", err)
        return err
    }
    return nil
}

func initCache() {
    var err error
    if cache == nil {
        cache, err = NewCacheRepository()
        if err != nil {
            log.Fatalf("Got error creating the cache repository: %s", err)
        }
    }
}