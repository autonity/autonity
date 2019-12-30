graph TB
    B---A
    C---A
    A---D
    A-->E
    subgraph one
    a1-->a2
    end
    subgraph two
    b1-->b2
    b2---b3
    end
    subgraph three
    c1---c2
    end