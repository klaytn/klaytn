package governance

// Compile-time type implementation check
var (
	_ HeaderEngine = (*Governance)(nil)
	_ ReaderEngine = (*Governance)(nil)
	_ Engine       = (*MixedEngine)(nil)
)
