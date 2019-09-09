# climbcomp-api

[![Build Status](https://travis-ci.com/climbcomp/climbcomp-api.svg?branch=master)](https://travis-ci.com/climbcomp/climbcomp-api)

Climbcomp gRPC server

## Install

In one window, checkout and run the API server:

```
git clone git@github.com:climbcomp/climbcomp-api.git
cd climbcomp-api

make build
make run
```

In another, use the CLI to make requests:

```
cd climbcomp-api

make bash

# Returns the version of the API server
# Run `climbcomp help` to see all options
climbcomp meta version
```

## Testing

```
make test
```
