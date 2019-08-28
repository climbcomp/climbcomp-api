# climbcomp-api

Climbcomp gRPC server

## Install

In one window, checkout and run the API server:

```
git clone git@github.com:climbcomp/climbcomp-api.git
cd climbcomp-api

make build
make run
```

In another, use prototool to make gRPC requests:

```
git clone git@github.com:climbcomp/climbcomp-proto.git
cd climbcomp-proto
brew install prototool

prototool grpc --address 0.0.0.0:3000 --method climbcomp.meta.v1.MetaAPI/GetVersion --data '{}'
```

## Testing

```
make test
```
