# Merkle Damgard experiment

Playing around with the merkle damgard construction.

## Specs

- Block size: 32 bytes
- The padding length (`p`) is `p = 32-len(M)%32` where `M` is the message to be hashed in bytes.
- The padding starts with a `1`. All other bytes are zeros.
- Additionally after the padding (`P`) has been applied to the message `M`, the length of `M` as int64 is appended to `M||P`.
- The compression function `f` is xor.
- The IV is `3.141592653589793238462643383279` (32 bytes).
 
 ![merkle damgard](merkle-damgard.png)
 
 
 ## Use
 1. `go build`
 2. `./md-experiment file.txt`