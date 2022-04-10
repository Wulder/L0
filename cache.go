package main

type Cache struct {
	RecentOrders map[string]Model //last orders which have been got

	UnwritedOrdes map[string]Model //last orders which have been got, but haven't ben written in dataBase
}

func (c *Cache) InitCache() {
	c.RecentOrders = make(map[string]Model)
	c.UnwritedOrdes = make(map[string]Model)
}
