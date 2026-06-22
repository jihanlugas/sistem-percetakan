package utils

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/labstack/gommon/log"
	"github.com/playwright-community/playwright-go"
)

// GeneratePDFWithChromedp masih dipertahankan nama fungsi untuk kompatibilitas,
// tapi sekarang menggunakan Playwright untuk menghasilkan PDF.
func GeneratePDFWithChromedp(inputHTML string) ([]byte, error) {
	absPath, err := filepath.Abs(inputHTML)
	if err != nil {
		log.Print("Err filepath.Abs ", err)
		return nil, err
	}

	if _, err := os.Stat(absPath); err != nil {
		log.Print("Err os.Stat ", err)
		return nil, fmt.Errorf("input HTML file not found: %w", err)
	}

	fileURL := (&url.URL{
		Scheme: "file",
		Path:   filepath.ToSlash(absPath),
	}).String()

	// Start Playwright
	pw, err := playwright.Run()
	if err != nil {
		log.Print("playwright.Run ", err)
		return nil, err
	}
	defer func() {
		_ = pw.Stop()
	}()

	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(true),
		Args:     []string{"--no-sandbox", "--disable-setuid-sandbox", "--disable-dev-shm-usage"},
	})
	if err != nil {
		log.Print("browser.Launch ", err)
		return nil, err
	}
	defer func() {
		_ = browser.Close()
	}()

	page, err := browser.NewPage()
	if err != nil {
		log.Print("browser.NewPage ", err)
		return nil, err
	}

	// Navigate to the file URL
	if _, err = page.Goto(fileURL, playwright.PageGotoOptions{WaitUntil: playwright.WaitUntilStateNetworkidle}); err != nil {
		log.Print("page.Goto ", err)
		return nil, err
	}

	// Give the page a short moment to settle
	time.Sleep(300 * time.Millisecond)

	// Get content height
	heightVal, err := page.Evaluate("() => document.body.scrollHeight")
	if err != nil {
		log.Print("page.Evaluate height ", err)
		return nil, err
	}

	var heightPx int
	switch v := heightVal.(type) {
	case float64:
		heightPx = int(v)
	case int:
		heightPx = v
	default:
		heightPx = 0
	}
	if heightPx <= 0 {
		heightPx = 1122 // fallback height
	}

	// Generate PDF using Playwright. Use pixel height to match content length.
	pdfBuf, err := page.PDF(playwright.PagePdfOptions{
		PrintBackground: playwright.Bool(true),
		Width:           playwright.String("5.9in"),
		Height:          playwright.String(fmt.Sprintf("%dpx", heightPx)),
		Margin: &playwright.Margin{
			Top:    playwright.String("0in"),
			Bottom: playwright.String("0in"),
			Left:   playwright.String("0in"),
			Right:  playwright.String("0in"),
		},
	})
	if err != nil {
		log.Print("page.PDF ", err)
		return nil, err
	}

	return pdfBuf, nil
}
