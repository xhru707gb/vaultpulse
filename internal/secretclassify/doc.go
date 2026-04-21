// Package secretclassify classifies Vault secret paths into sensitivity tiers
// (public, internal, confidential, secret) using configurable pattern rules.
//
// Usage:
//
//	c, err := secretclassify.New([]secretclassify.Rule{
//		{Pattern: "prod/", Level: secretclassify.LevelSecret},
//		{Pattern: "staging/", Level: secretclassify.LevelConfidential},
//	}, secretclassify.LevelInternal)
//	results := c.ClassifyAll(paths)
package secretclassify
