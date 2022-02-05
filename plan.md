# Language Plan

All examples shown are standard operators used in mathematics or in the `C` programming language.

## Operator Tesseract

Each and every operator is explicitly identifiable by a number of factors:

- The symbol used to represent the operator
- The precedence of the operator, which defines behaviours of the operator
- The types used by the operator.

Each precedence level holds a set of operators with known properties, and some support operation requiring no symbol.

Therefore, the set of all available operators can be represented using a table such as:

Precedence | Layer type | NULL | `==` | `>` | `<` | `+` | `-` | `*` | `/` | `%`
---|---|---|---|---|---|---|---|---|---|---
0 | Implied Operation Weak Left Associative Infix Binary | `bool_x2 -> bool` | `any_x2 -> bool` | `numeric_x2 -> bool` | `numeric_x2 -> bool` | - | - | - | - | -
1 | Left Associative Infix Binary | - | - | - | - | `numeric_x2 -> numeric` | `numeric_x2 -> numeric` | - | - | -
2 | Repeatable Prefix Unary | - | - | - | - | `numeric -> numeric` | `numeric -> numeric` | - | - | -
3 | Left Associative Infix Binary | `numeric_x2 -> numeric` | - | - | - | - | - | `numeric_x2 -> numeric` | `numeric_x2 -> numeric` | `numeric_x2 -> numeric`

## Operators

### Unary Operators

A Unary operator is one which concerns a single argument.

#### Prefix / Postfix

Any unary operator can be defined as prefix (`+x`) or postfix (`x!`). An operator level must enforce one of these two formats.

#### Repeatability

Unary operators which return the same type as the one they take in can be considered repeatable. For a layer to be repeatable, all operators within it must make use of the same type. Eg. `++-+--++---x`

A layer does not need to be repeatable for many operations within said layer to be chained together.

### Binary Operators

A binary operator is one which concerns two arguments.

#### Associativity

A binary operator can be left-associative, right-associative, or just not associative at all. For such an operator to be associative in any form, it must concern arguments of a single type.

##### Associativity Strength

The "strength" of an operator level's associativity gives a sense of whether operators within it require directional associativity.

- Weak associativity is where all operators within the level can be parsed in either direction. ie. they are all entirely associative.
- Strong associativity is where all operators within the level have to be parsed in a specific direction.

#### Implied Operations

Some operators require some implied operation for chaining and therefore associativity to exist. For example, the operator `<` is chained like `x < y < z < w` despite the use of regular associativity rules not holding. The implied operation in this case is conjunction. Showing the same expression in `C` makes this clear: `x < y & y < z & z < w`.

These implied operations can also be used with ordinary associative operators, in which case they act when the two arguments are placed adjacent to one another. For example: `3 * x` can be rephrased as `3x` using the implied operation for multiplication.

#### Prefix / Infix / Postfix

A binary operator can be defined as prefix (`+ x y`), infix (`x + y`), or postfix (`x y +`). An operator level must enforce one of these formats. For implied operations to work as intended, a binary operator level must enforce the infix format.

### Ternary Operators

A ternary operator is one which concerns three arguments.

### Extended Operators

Operators making use of the prefix or postfix notation can be extended indefinitely as long as there do not exist any implicit operations. For operators with multiple arguments, associativity works identically to binary operators. In order to support extended infix operators, a symbol must separate all arguments. If not, at least one is necessary, and the same implicit operation restriction applies.

### Prefix Flush

During parsing, when encountering a prefix operator of lower precedence than the most recently parsed operator, a prefix flush occurrs. This phenomenon only concerns a single term in the expression, applying the prefix operator and returning to the original precedence level the parser was working at. Functionally, this means for an expression such as `2 * -x * y`, it is parsed `2 * (-x) * y`. However, as `-x` still has a lower precedence, `-2 * x * y` is parsed `-(2 * x * y)`.

### Colliding Operator Symbols

Multiple operators may have the same symbol. In some cases, this results in one or more operators which are left unused, and therefore will prematurely end translation. The simplest case where multiple operators can have the same symbol is when one is a prefix operator, and the other is not. Therefore, there can only ever be at most two operators a symbol can represent.
