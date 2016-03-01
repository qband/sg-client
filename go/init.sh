go get -u github.com/constabulary/gb/...

#gb vendor fetch src.sourcegraph.com/sourcegraph
git clone https://src.sourcegraph.com/sourcegraph vendor/src/src.sourcegraph.com/sourcegraph
cp -r vendor/src/src.sourcegraph.com/sourcegraph/vendor/* vendor/src/
