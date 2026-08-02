package app

import "github.com/SergioLacerda/rpg-system-rgb/internal/components/publication"

// LibraryOptions configures static HTML Library generation.
type LibraryOptions = publication.LibraryOptions

// PDFOptions configures PDF publication.
type PDFOptions = publication.PDFOptions

// PublicationCheckOptions configures publication smoke checks.
type PublicationCheckOptions = publication.CheckOptions

// BuildLibrary renders docs/core/** into the public Library directory.
func BuildLibrary(options LibraryOptions) error {
	return publication.BuildLibrary(options)
}

// PublishPDFs publishes curated PDFs to latest and versioned download names.
func PublishPDFs(options PDFOptions) error {
	return publication.PublishPDFs(options)
}

// CheckPublication validates generated public publication artifacts.
func CheckPublication(options PublicationCheckOptions) error {
	return publication.Check(options)
}
