package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

func generatePDF(outputDir, srcDir string) error {
	resumeHTMLPath := filepath.Join(outputDir, "resume-printable.html")
	pdfOutputPath := filepath.Join(srcDir, "pdfs", "ade-sede.pdf")

	if _, err := os.Stat(resumeHTMLPath); os.IsNotExist(err) {
		return fmt.Errorf("resume-printable.html not found at %s", resumeHTMLPath)
	}

	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	var pdfBuf []byte
	fileURL := fmt.Sprintf("file://%s", resumeHTMLPath)

	err := chromedp.Run(ctx,
		chromedp.Navigate(fileURL),
		chromedp.WaitVisible("body", chromedp.ByQuery),
		chromedp.Sleep(2*time.Second),
		chromedp.ActionFunc(func(ctx context.Context) error {
			buf, _, err := page.PrintToPDF().
				WithPrintBackground(true).
				WithPreferCSSPageSize(true).
				WithMarginTop(0.0).
				WithMarginBottom(0.0).
				WithMarginLeft(0.0).
				WithMarginRight(0.0).
				WithPaperWidth(8.27).  // A4 width in inches
				WithPaperHeight(11.7). // A4 height in inches
				Do(ctx)
			if err != nil {
				return err
			}
			pdfBuf = buf
			return nil
		}),
	)

	if err != nil {
		return fmt.Errorf("failed to generate PDF: %v", err)
	}

	if err := os.MkdirAll(filepath.Dir(pdfOutputPath), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	if err := os.WriteFile(pdfOutputPath, pdfBuf, 0644); err != nil {
		return fmt.Errorf("failed to write PDF file: %v", err)
	}

	log.Printf("PDF generated successfully: %s", pdfOutputPath)
	return nil
}
