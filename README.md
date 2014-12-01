qeconv
======

[![Build Status](https://travis-ci.org/hiwane/qeconv.svg?branch=master)](https://travis-ci.org/hiwane/qeconv)
[![Coverage Status](https://img.shields.io/coveralls/hiwane/qeconv.svg)](https://coveralls.io/r/hiwane/qeconv?branch=master)


convert first-order formulas


# Install

After installing `golang` and setting up your `GOPATH`, execute follows: 
```sh
go get github.com/hiwane/qeconv/qeconv
```

# Usage

```
Usage: qeconv [-f from][-t to][-i inputfile][-o outputfile]
    -f: Use from for input format           [syn]
    -t: Use to for output format {math|tex} [math]
    -i: Use inputfile for input             [stdin]
    -o: Use outputfile for outpuut          [stdout]
```

# Examples


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

