## ./pachctl glob-file

Return files that match a glob pattern in a commit.

### Synopsis


Return files that match a glob pattern in a commit (that is, match a glob pattern
in a repo at the state represented by a commit). Glob patterns are
documented [here](https://golang.org/pkg/path/filepath/#Match).

Examples:

```sh

# Return files in repo "foo" on branch "master" that start
# with the character "A".  Note how the double quotation marks around "A*" are
# necessary because otherwise your shell might interpret the "*".
$ pachctl glob-file foo master "A*"

# Return files in repo "foo" on branch "master" under directory "data".
$ pachctl glob-file foo master "data/*"

```

```
./pachctl glob-file repo-name commit-id pattern
```

### Options

```
      --raw   disable pretty printing, print raw json
```

### Options inherited from parent commands

```
      --no-metrics   Don't report user metrics for this command
  -v, --verbose      Output verbose logs
```

### SEE ALSO
* [./pachctl](./pachctl.md)	 - 

###### Auto generated by spf13/cobra on 2-Aug-2017
