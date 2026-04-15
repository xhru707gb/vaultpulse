// Package health provides functionality for checking the operational
// status of a HashiCorp Vault instance, including sealed/unsealed state
// and response latency measurements.
//
// Usage:
//
//	checker := health.NewChecker(vaultClient)
//	status, err := checker.Check(ctx)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println(health.FormatTable([]health.Status{status}))
package health
