// Package masking provides configurable field-level masking for sensitive
// process and connection metadata emitted by portwatch events.
//
// Masking rules can be applied to fields such as process names, user
// identifiers, command-line arguments, and network addresses. Each rule
// specifies a field selector and a replacement strategy (e.g., redact,
// hash, or truncate).
//
// Example usage:
//
//	masker := masking.New(masking.Config{
//		Rules: []masking.Rule{
//			{Field: "process.args", Strategy: masking.Redact},
//			{Field: "user.name", Strategy: masking.Hash},
//		},
//	})
//	
//	masked := masker.Apply(event)
package masking
