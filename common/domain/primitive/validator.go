package primitive

var (
	licenseValidator licensePrimitiveValidator
)

func Init(
	license licensePrimitiveValidator,
) {
	licenseValidator = license
}

type licensePrimitiveValidator interface {
	IsValidLicense(string) bool
}
