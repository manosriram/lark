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
fn addStaticVars[] <<
  return 100+500;
>>

fn addWithNoReturn[] <<
  a <- 100;
>>

fn addWithArgs[a,b] <<
  return a+b;
>>

fn addWithArgsAndLocalVar[a,b] <<
  local c <- 3;
  return a+b+c;
>>

fn addWithDynamicArgs[a,b] <<
  return a+b;
>>

addWithNoReturn();
staticsum <- addStaticVars();
sumOne <- addWithArgs(1, 2); // 3
sumTwo <- addWithArgsAndLocalVar(1, 2); // 6

first <- 100;
second <- 200;
sumThree <- addWithDynamicArgs(first+(5-3)*second, second); // 700
```

### Usage
##### Using Makefile
`make run`

##### Source
`go build -o lark && ./lark <filename>`

[Progress Board](https://trello.com/b/1qAWAjZS/lark)
