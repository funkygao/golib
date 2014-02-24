/*
A slab allocator implementation.

Arena class2's chunk size = growthFactor * class1's chunk size 

slabClass1
    |<----------------------- slabSize ------------>|
    +-----------------------------------------------+
    | chunk | chunk | chunk | ...           | chunk |
    +-----------------------------------------------+

slabClass2
    +-----------------------------------------------+
    | chunk     | chunk     | ...       | chunk     |
    +-----------------------------------------------+

slabClassN
    |<----------------------- slabSize ------------>|
    +-----------------------------------------------+
    | chunk             | ...         | chunk       |
    +-----------------------------------------------+
    |<-   chunkSize --->|
*/
package slab
