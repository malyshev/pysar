package editorial

// intake, research, draft, staff-edit, sharpen, humanize, export, and
// discoverability are real passes, each in its own file (dec-20260718-
// ffbfb04b, dec-20260725-35fa2d24, dec-20260804-e3234e50).

func init() {
	Register(intakePass{})
	Register(researchPass{})
	Register(draftPass{})
	Register(staffEditPass{})
	Register(sharpenPass{})
	Register(humanizePass{})
	Register(exportPass{})
	Register(discoverabilityPass{})
}
