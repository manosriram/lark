## lark
> 

### Documentation
> types

- int
- float
- boolean
- string

> variable initialization

```
a <- 1;
b <- "lark";
c <- 3.14;
d <- true;
e <- 123;
f <- (1+2+3) + c + (3*1/10);
unary <- !false;
```

> operators

```
*, /, (
+, -
>, <, >=, <=, !, !=
```

> control-flow

```
a <- (1+3)==2; // false
b <- true;
if (a==b) {
  c <- "ok";
} else {
  c <- "not_ok";
}
```

> swap variables

```
a <- 1;
b <- -123;
a <-> b;
/* a = -123, b = 1 */
```

> functions

```
fn addstatic[] <<
  return 100+500;
>>

fn addwithnoreturn[] <<
  a <- 100;
>>

a <- 500;
staticsum <- addstatic();
addwithnoreturn(); // a = 100
```

### Usage
##### Using Makefile
`make run`

##### Source
`go build -o lark && ./lark <filename>`

[Progress Board](https://trello.com/b/1qAWAjZS/lark)
