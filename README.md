# gRPC server reflection test

This codebase demonstrates an issue with a private Armeria gRPC reflection setup.

Run as

	go run refl.go localhost:8080

to see how second request receives two identical responses.

Wire shark shows that 3 requests have 5 (sometimes only 4) responses
instead of the expected 3 responses:

![wireshark screenshot][wireshark]

[wireshark]: wireshark.png "Wireshark capture of single refl.go execution"

This issue could not be reproduced when run against a bare bones Java example following

* https://grpc.io/docs/languages/java/quickstart/
* https://github.com/grpc/grpc-java/blob/master/documentation/server-reflection-tutorial.md

we still need to work out if this is:

* related to a local Armeria setup
* gRPC reflection specifc
* a generic Armeria streaming issue
