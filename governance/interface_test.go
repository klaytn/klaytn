package governance

// Compile-time type implementation check
var _ HeaderEngine = (*Governance)(nil)
var _ Engine = (*MixedEngine)(nil)
