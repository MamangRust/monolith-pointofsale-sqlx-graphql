package utils

import (
	"regexp"
	"testing"
)

var slugPattern = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*-[a-z0-9]{4}$`)
var barcodePattern = regexp.MustCompile(`^[A-Z]+[0-9]{6}$`)

func TestGenerateSlug(t *testing.T) {
	slug := GenerateSlug("Hello World 123")
	if !slugPattern.MatchString(slug) {
		t.Fatalf("slug %q does not match pattern", slug)
	}
}

func TestGenerateSlugLowercasesAndSlugifies(t *testing.T) {
	slug := GenerateSlug("Kopi Susu Gula Aren")
	// Base portion must be lowercase slug of the words (dash-separated).
	if !regexp.MustCompile(`^kopi-susu-gula-aren-[a-z0-9]{4}$`).MatchString(slug) {
		t.Fatalf("slug %q not derived from 'Kopi Susu Gula Aren'", slug)
	}
}

func TestGenerateBarcode(t *testing.T) {
	barcode := GenerateBarcode("Kopi Susu")
	if !barcodePattern.MatchString(barcode) {
		t.Fatalf("barcode %q does not match pattern", barcode)
	}
}

func TestGenerateBarcodeUppercaseInitials(t *testing.T) {
	barcode := GenerateBarcode("es teh manis")
	if !regexp.MustCompile(`^ETM[0-9]{6}$`).MatchString(barcode) {
		t.Fatalf("barcode %q not derived from 'es teh manis'", barcode)
	}
}

func TestGenerateBarcodeEmptyNameFallback(t *testing.T) {
	barcode := GenerateBarcode("")
	if !regexp.MustCompile(`^PRD[0-9]{6}$`).MatchString(barcode) {
		t.Fatalf("barcode %q does not use PRD fallback", barcode)
	}
}
