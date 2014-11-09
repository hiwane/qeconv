qeconv
======

convert first-order formulas

# install

```sh
go get github.com/hiwane/qeconv/qeconv
```

# usage

```
Usage: qeconv [-f from][-t to][-i inputfile][-o outputfile]
    -f: Use from for input CAS language [syn]
    -t: Use to for output CAS language [math]
    -i: Use inputfile for input [stdin]
    -o: Use outputfile for outpuut [stdout]
```

# Example

```
> echo "x>0:" | qeconv
0 < x
> echo "And(x>0,y>=0):" | qeconv
0 < x && 0 <= y
> echo "Ex([y],And(x>0,y>=0)):" | qeconv
Exists[{y},0 < x && 0 <= y]
```

