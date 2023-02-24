# b00m-in/data

# Contents

Simple postgres db interface for b00m-in/gin.

Uses an rtree to provide clustering by location of pubs/devices.

# Usage

`git clone https://github.com/b00m-in/data` and use `replace` in your `go.mod`, so for example if you import this package as `b00m.in/data` in your source:

`import b00m.in/data` 

Then replace `b00m.in/data` in your `go.mod`

```
require b00m.in/data v0.0.0
replace b00m.in/data => ../data` 
```

Don't use `go get github.com/b00m-in/data` or use this package from module cache as this module declares itself as `b00m.in/data`. All release versions have been deleted.
