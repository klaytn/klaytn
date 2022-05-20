package governance

// Compile-time type implementation check
var (
	_ HeaderEngine = (*Governance)(nil)
	_ ReaderEngine = (*ContractEngine)(nil)
	_ Engine       = (*MixedEngine)(nil)
)
