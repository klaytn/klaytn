package core

// Constants containing the genesis allocation of built-in genesis blocks.
// Their content is an RLP-encoded list of (address, balance) tuples.
// Use mkalloc.go to create/update them.

// nolint: misspell
const mainnetAllocData = "\xda\u0654\x19t\x18\x1a?\x1bE+\x00\xec\u06ab\x19\xa2j\xf3\x8f\xcb\xff\x1e\x83\x98\x96\x80"

const testnetAllocData = "\xda\u0654\x19t\x18\x1a?\x1bE+\x00\xec\u06ab\x19\xa2j\xf3\x8f\xcb\xff\x1e\x83\x98\x96\x80"
