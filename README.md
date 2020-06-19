# replacer

A simple little command utility to run a basic replacement of strings a file.

Does no fancy pattern matching, no crazy syntax. Just give it a list of things to replace, and a list of things to replace them with, and it does it.

### Example

Given this file:

```Go
// +build ignore

package somemap

import "fmt"

func GetThing(things map[KEYTYPE]VALTYPE, key KEYTYPE) (val VALTYPE, err error) {
    if val, exists := things[key]; !exists {
        err = fmt.Errorf("Could not find thing!")
    }
    return
}
```

Run this command:

```BASH
$ replacer --in=thingmap.go --out=thingmap_stringint.go --from="KEYTYPE,// +build ignore,VALTYPE" --to=string,,int
```

(Note that empty comma in `--to=` will make it just remove that 2nd entry by replacing it with an empty string. Same applies if the comma is trailing, then the last thing in `--from=` will effectively be removed.)

Get this output:

```Go

package somemap

import "fmt"

func GetThing(things map[string]int, key string) (val int, err error) {
    if val, exists := things[key]; !exists {
        err = fmt.Errorf("Could not find thing!")
    }
    return
}
```
