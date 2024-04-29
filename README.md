# Go Cov

Go cov is a simple tool to filter out results from the go test coverage profile. 
This usually means that we'll filter out generated files for which we don't want to write unit tests. However, we don't want these files to
negatively impact our overall test coverage.

The rool reads a given coverage profile, filters out the lines that match the specified regulare expressions in the `.coverage-ignore.yaml` and writes
the result back to the coverage report.

# Installing

```bash
go install github.com/wimspaargaren/go-cover-ignore@latest
```

# Running

```bash
GO_COVER_IGNORE_SPEC_PATH=".coverage-ignore.yaml" GO_COVER_IGNORE_COVER_PROFILE_PATH="cover.out" go-cover-ignore
```


# Coverage Ignore Yaml

The program expectes a file called `.coverage-ignore.yaml` to be specified, which has the following structure:

```YAML
module: "github.com/<my_org>/<my_project>"
ignore_rules:
    - <my_regex_to_ignore>
```
