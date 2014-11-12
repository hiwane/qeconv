qeconv
======

[![Build Status](https://travis-ci.org/hiwane/qeconv.svg?branch=master)](https://travis-ci.org/hiwane/qeconv)

convert first-order formulas


# install

```sh
go get github.com/hiwane/qeconv/qeconv
```

# usage

```
Usage: qeconv [-f from][-t to][-i inputfile][-o outputfile]
    -f: Use from for input CAS language [syn]
    -t: Use to for output CAS language {math|tex} [math]
    -i: Use inputfile for input [stdin]
    -o: Use outputfile for outpuut [stdout]
```

# Example


### SyNRAC to Mathematica (default)

```
> echo "x>0:" | qeconv
0 < x
> echo "And(x<>0,y>=0):" | qeconv
x != 0 && 0 <= y
> echo "Ex([y],And(x>0,y>=0)):" | qeconv
Exists[{y},0 < x && 0 <= y]
```

### SyNRAC to LaTeX

```
> echo "x>0:" | qeconv -t tex
0 < x
> echo "And(x<>0,y>=0):" | qeconv -t tex
x \neq 0 \land 0 \leq y
> echo "Ex([y],And(x>0,y>=0)):" | qeconv -t tex
\exists y(0 < x \land 0 \leq y)
```

