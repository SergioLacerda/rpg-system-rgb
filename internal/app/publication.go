package app

import "github.com/SergioLacerda/rpg-system-rgb/internal/components/publication"

// PDFOptions configures PDF publication.
type PDFOptions = publication.PDFOptions

// PublicationCheckOptions configures publication smoke checks.
type PublicationCheckOptions = publication.CheckOptions

// PublishPDFs publishes curated PDFs to latest and versioned download names.
func PublishPDFs(options PDFOptions) error {
	return publication.PublishPDFs(options)
}

// CheckPublication validates generated public publication artifacts.
func CheckPublication(options PublicationCheckOptions) error {
	return publication.Check(options)
}
