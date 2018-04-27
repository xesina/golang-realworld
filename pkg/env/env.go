package env

const (
	// Development environment
	Development = "development"
	// Test environment
	Test = "test"
	// Staging environment
	Staging = "staging"
	// Production environment
	Production = "production"
)

// IsRelease return true if environment is either staging or production
func IsRelease(e string) bool {
	//panic(e)
	return e == Staging || e == Production
}

// IsDevelopment return true if environment is development
func IsDevelopment(e string) bool {
	return e == Development
}

// IsStaging return true if environment is staging
func IsStaging(e string) bool {
	return e == Staging
}
