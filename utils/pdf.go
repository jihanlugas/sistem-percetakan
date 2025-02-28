package utils

import (
	"context"
	"fmt"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"os"
	"time"
)

// Generate PDF menggunakan chromedp dan cdproto/page
func GeneratePDFWithChromedp(inputHTML string) ([]byte, error) {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	fullPath, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	htmlFile := fmt.Sprintf("file://%s/%s", fullPath, inputHTML)
	var pdfBuffer []byte

	// Menggunakan cdproto/page untuk generate PDF
	err = chromedp.Run(ctx,
		chromedp.Navigate(htmlFile),
		chromedp.ActionFunc(func(ctx context.Context) error {
			buf, _, err := page.PrintToPDF().
				WithPrintBackground(true).
				WithPaperWidth(4.1).  // A6 width in inches
				WithPaperHeight(5.8). // A6 height in inches
				//WithScale(0.5).
				Do(ctx)
			if err != nil {
				return err
			}
			pdfBuffer = buf
			return nil
		}),
	)
	if err != nil {
		return nil, err
	}

	return pdfBuffer, nil
}
