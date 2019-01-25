## setup

remove template gRPC codes
```
$ git submodule deinit go-grpc-template/
$ git rm go-grpc-template
```

add your gRPC codes
```
$ git submodule add <YOUR_gRPC_REPOSITORY_URL>
(Ex. $ git submodule add -b generated/go git@github.com:kaito2/grpc-gen-circleci-template.git)
```

replace your code (`helloworld.pb.go` -> your pb)