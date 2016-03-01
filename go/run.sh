set -eu
gb build

bin/repo "$repoURI" "$cloneURL"
