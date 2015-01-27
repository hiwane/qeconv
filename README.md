qeconv
======

[![Build Status](https://travis-ci.org/hiwane/qeconv.svg?branch=master)](https://travis-ci.org/hiwane/qeconv)
[![Coverage Status](https://coveralls.io/repos/hiwane/qeconv/badge.svg?branch=master)](https://coveralls.io/r/hiwane/qeconv?branch=master)


Convert first-order formulas


# Install

After installing `golang` and setting up your `GOPATH`, execute follows: 
```sh
go get github.com/hiwane/qeconv/qeconv
```
`qeconv` will be installed to the `${GOPATH}/bin` directory.
For convenience, add the `${GOPATH}/bin` to your `PATH`.

# Usage

```
Usage: qeconv [-f {from}][-t {to}][-i {inputfile}][-o {outputfile}][-s][-n {{index}]
    -f: Use {from} for input format                       [syn]
    -t: Use {to} for output format                        [syn]
		 syn:  SyNRAC format
		 tex:  LaTeX format
         math: MATHEMATICA format
		 qep:  QEPCAD format
		 red:  Redlog format
		 rc:   RegularChains format
		 smt2: SMT2 format (-n option is reguired)
    -i: Use {inputfile} for input                         [stdin]
    -o: Use {outputfile} for output                       [stdout]
	-s: Remove duplicate formulas
	-n: Use {index}th formula                             [0]
	      0:   Use all formula
		 -1:   Get the number of first-order formulas in input
```


# Examples


### SyNRAC to SyNRAC (default)

```sh
% echo "x>0:" | qeconv
0 < x:
% echo "And(x<>0,y>=0):" | qconv
And(x <> 0,0 <= y):
% echo "Ex([y],And(x>0,y>=0)):" | qconv
Ex([y],And(0 < x,0 <= y)):
```

### SyNRAC to LaTeX

```sh
% echo "x>0:" | qeconv -t tex
0 < x
% echo "And(x<>0,y>=0):" | qeconv -t tex
x \neq 0 \land 0 \leq y
% echo "Ex([y],And(x>0,y>=0)):" | qeconv -t tex
\exists y(0 < x \land 0 \leq y)
```

### SyNRAC to Mathematica

```sh
% echo "x>0:" | qeconv -t math
0 < x
% echo "And(x<>0,y>=0):" | qeconv -t math
x != 0 && 0 <= y
% echo "Ex([y],And(x>0,y>=0)):" | qeconv -t math
Exists[{y},0 < x && 0 <= y]
```

### SyNRAC to QEPCAD

- http://www.usna.edu/CS/qepcadweb/B/QEPCAD.html
- https://github.com/hiwane/qepcad/

```sh
% echo "x>0:" | qeconv -t qep
0 < x
% echo "And(x<>0,y>=0):" | qeconv -t qep
x /= 0 /\ 0 <= y
% echo "Ex([y],And(x>0,y>=0)):" | qeconv -t qep
(E y)[0 < x /\ 0 <= y]
```

### SyNRAC to Redlog/REDUCE

- http://www.redlog.eu/
- http://reduce-algebra.sourceforge.net/


```sh
% echo "x>0:" | qeconv -t red
0 < x
% echo "And(x<>0,y>=0):" | qeconv -t red
x <> 0 and 0 <= y
% echo "Ex([y],And(x>0,y>=0)):" | qeconv -t red
ex([y],0 < x and 0 <= y)
```

### SyNRAC to RegularChains

- http://regularchains.org/

```sh
% echo "x>0:" | qeconv -t rc
0 < x
% echo "And(x<>0,y>=0):" | qeconv -t rc
`&and`(x <> 0, 0 <= y)
% echo "Ex([y],And(x>0,y>=0)):" | qeconv -t rc
`&E`([y]), `&and`(0 < x,0 <= y)
```


<!-- vim: set spell: -->
