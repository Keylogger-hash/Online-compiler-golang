package storage

import (
	"github.com/bradfitz/gomemcache/memcache"
)




func NewMemcachedClient() (*memcache.Client,error){
	c := memcache.New("localhost:11211")
	err := c.Ping()
	if err != nil {
		return nil, err
	}
	return c,nil
	
}

func MemcachedGetValue(client *memcache.Client, sha256key string) ([]byte,error) {
	item, err := client.Get(sha256key)
	if err != nil {
		return nil,err
	}
	return item.Value,nil
}
func MemcachedAddValue(client *memcache.Client,sha256key string, body []byte) error{
	item := &memcache.Item{Key:sha256key,Value:body,Expiration:0}
	err := client.Set(item)
	if err != nil {
		return err
	}
	return nil
}