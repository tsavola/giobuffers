// Copyright (c) 2022 Timo Savola. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resource

type cache struct {
	synced map[string]struct{}
}

func (c *cache) set(id string) {
	if c.synced == nil {
		c.synced = make(map[string]struct{})
	}
	c.synced[id] = struct{}{}
}

func (c *cache) clear() {
	for id := range c.synced {
		delete(c.synced, id)
	}
}

type Context struct {
	older cache
	newer cache
}

func (c *Context) Synced() {
	c.older, c.newer = c.newer, c.older
	c.newer.clear()
}

type Resource struct {
	c     *Context
	ID    string
	valid bool
}

func MakeResource(c *Context, id string) Resource {
	return Resource{
		c:  c,
		ID: id,
	}
}

func (r *Resource) NeedSync() bool {
	if !r.valid {
		return true
	}
	_, ok := r.c.older.synced[r.ID]
	return !ok
}

func (r *Resource) Synced() {
	r.c.newer.set(r.ID)
	r.valid = true
}
