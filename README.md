golib
=====

golang common facilities lib

[![Build Status](https://travis-ci.org/funkygao/golib.png?branch=master)](https://travis-ci.org/funkygao/golib)

### TODO

* [ ] hot reload mechanism in server pkg
* [ ] tweak performance of cache.LruCache
* [ ] merge mutexmap with cmap, lur cache will be sharded
* [ ] merge pqueue and gopqueue for priority queue
* [ ] merge recycler and slab pkg
* [ ] str pkg StringBuilder for better performance
* [X] shard lru cache to lower mutex race
* [ ] replace lru cache container/list with slice
  - https://github.com/golang/go/wiki/SliceTricks
* [ ] https://github.com/dgryski/go-clockpro
