byten
=====
Bit the size of file and turn it into human readable format. 

This is a Go Package that convert size of file into human readable format.

Another weekend project by [pyk](http://google.com/+bayualdiyansyah).
### Usage
First thing first, get the remote package
```
$ go get github.com/pyk/byten
```
and import package into your project
```
import "github.com/pyk/byten"
```
then bite the bytes!
```
byten.Size(1024) # => 1.0KB
byten.Size(206848) # => 202KB
byten.Size(10239999998976) # => 9.3TB
byten.Size(6314666666666665984) # => 5.5EB
```
easy huh? =))

### Docs

Nothing fancy, but you can see [here](https://godoc.org/github.com/pyk/byten)

### License
[MIT License](https://github.com/pyk/byten/blob/master/LICENSE)
