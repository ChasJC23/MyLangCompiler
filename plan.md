# Language Plan

All examples shown are standard operators used in modern mathematical notation or in the `C` programming language.

## Operator Tesseract

Each and every operator is explicitly identifiable by a number of factors:

- The symbol used to represent the operator
- The precedence of the operator, which defines behaviours of the operator
- The types used by the operator.

Each precedence level holds a set of operators with known properties, and some support operation requiring no symbol.

Therefore, the set of all available operators can be represented using a table such as:

| Precedence | Layer type                                           | NULL                    | `==`             | `>`                  | `<`                  | `+`                     | `-`                     | `*`                     | `/`                     | `%`                     |
|------------|------------------------------------------------------|-------------------------|------------------|----------------------|----------------------|-------------------------|-------------------------|-------------------------|-------------------------|-------------------------|
| 0          | Implied Operation Weak Left Associative Infix Binary | `bool_x2 -> bool`       | `any_x2 -> bool` | `numeric_x2 -> bool` | `numeric_x2 -> bool` | -                       | -                       | -                       | -                       | -                       |
| 1          | Left Associative Infix Binary                        | -                       | -                | -                    | -                    | `numeric_x2 -> numeric` | `numeric_x2 -> numeric` | -                       | -                       | -                       |
| 2          | Repeatable Prefix Unary                              | -                       | -                | -                    | -                    | `numeric -> numeric`    | `numeric -> numeric`    | -                       | -                       | -                       |
| 3          | Left Associative Infix Binary                        | `numeric_x2 -> numeric` | -                | -                    | -                    | -                       | -                       | `numeric_x2 -> numeric` | `numeric_x2 -> numeric` | `numeric_x2 -> numeric` |

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
*Note: the chosen interpretation for the precedence of juxtaposition places it at the same precedence as multiplication, meaning an expression such as `a/3b` would be interpreted `((a / 3) * b)`. This was chosen as the alternative interpretation would require an entirely new level of precedence used for the juxtaposition alone.*

Implied operations must always use binary infix notation.

#### Prefix / Infix / Postfix

A binary operator can be defined as prefix (`+ x y`), infix (`x + y`), or postfix (`x y +`). An operator level must enforce one of these formats. For implied operations to work as intended, a binary operator level must enforce the infix format.

### Extended Operators

Operators making use of the prefix or postfix notation can be extended indefinitely as long as there do not exist any implicit operations. For operators with multiple arguments, associativity works identically to binary operators. In order to support extended infix operators, a symbol must separate all arguments. If not, at least one is necessary, and the same implicit operation restriction applies.

### Prefix Flush

During parsing, when encountering a prefix operator of lower precedence than the most recently parsed operator, a prefix flush occurrs. This phenomenon only concerns a single term in the expression, applying the prefix operator and returning to the original precedence level the parser was working at. Functionally, this means for an expression such as `2 * -x * y`, it is parsed `2 * (-x) * y`. However, as `-x` still has a lower precedence, `-2 * x * y` is parsed `-(2 * x * y)`.

### Colliding Operator Symbols

Multiple operators may have the same symbol. In some cases, this results in one or more operators which are left unused, and therefore will prematurely end translation. The simplest case where multiple operators can have the same symbol is when one is a prefix operator, and the other is not. Therefore, there can only ever be at most two operators a symbol can represent.

### Many Symbol Operators

Operators which enforce prefix or infix notation may optionally require additional symbols between arguments. This syntax allows for operators such as the ternary operator present in C: `predicate ? expr1 : expr2`.

## Default Operators

Many of the default operators have special syntax features which cannot be used in user-defined operators.

| Operator    | ASCII Alias    | Type                     | Description                                                                                                                                                                                                                                                     |
|-------------|----------------|--------------------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `_ ≔ _`     | `_ := _`       | Non-Associative Binary   | This operator declares and defines the left side variable using the expression supplied on the right side.                                                                                                                                                      |
| `_ ← _`     | `_ <- _`       | Non-Associative Binary   | This operator assigns the right side expression to the left side variable.                                                                                                                                                                                      |
| `_ → _`     | `_ -> _`       | Right-Associative Binary | This type operator represents a mapping from one set to another, specified by its two arguments.                                                                                                                                                                |
| `_ ∈ _`     | `_ E _`        | Non-Associative Binary   | This operator specifies the set, or type, a variable belongs to.                                                                                                                                                                                                |
| `_ ? _ : _` | Not Applicable | Right-Associative Binary | This operator, most commonly written in the form `_ ? {} : {}`, represents an if-else statement. The else clause is optional, and if both clauses are single expressions, behaviour is identical to the ternary operator present in the C programming language. |
| `_ ⁇ _`     | `_ ?? _`       | Right-Associative Binary | This operator, most commonly written in the form `_ ?? {}`, represents a while loop.                                                                                                                                                                            |
| `_ _ ⁇`     | `_ _ ??`       | Postfix                  | This operator, most commonly written in the form `{} _ ??`, represents a do-while loop.                                                                                                                                                                         |
| `prog {}`   | Not Applicable | Prefix                   | This operator specifies the entrypoint to a program compiled using this source file.                                                                                                                                                                            |
| `before _`  | Not Applicable | Prefix                   | Operator for selecting a higher precedence than the one supplied. Can be applied to a precedence level or a particular operator within a precedence level using the `_` syntax. Refers to a new precedence level when combined with `prec _`.                   |
| `with _`    | Not Applicable | Prefix                   | Operator for selecting the precedence level of the supplied operator. Can be applied to a precedence level or a particular operator within a precedence level using the `_` syntax.                                                                             |
| `after _`   | Not Applicable | Prefix                   | Operator for selecting a lower precedence than the one supplied. Can be applied to a precedence level or a particular operator within a precedence level using the `_` syntax. Refers to a new precedence level when combined with `prec _`.                    |
| `prec _`    | Not Applicable | Prefix                   | Operator for defining a new precedence level.                                                                                                                                                                                                                   |

