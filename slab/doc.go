/*
A slab allocator implementation.

    Arena

        slabClass
        +-----------+
        |   |   |   |
        +-----------+

            slab
            +-----------+
            |   |   |   |
            +-----------+


            chunk
            +------+
            | refs |
            +------+
            | self |
            | next |
            +------+

*/
package slab
