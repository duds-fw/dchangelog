package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/jung-kurt/gofpdf"
)

// getGitChanges fetches all modified files and their diffs
func getGitChanges(baseBranch, targetBranch string) (map[string]string, error) {
	cmd := exec.Command("git", "diff", baseBranch+".."+targetBranch, "--name-only --pretty=format:* %h - %s (%an, %ad)")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	changes := make(map[string]string)
	files := strings.Split(out.String(), "\n")

	for _, file := range files {
		file = strings.TrimSpace(file)
		if file == "" {
			continue
		}

		diffCmd := exec.Command("git", "diff", baseBranch+".."+targetBranch, "--", file)
		var diffOut bytes.Buffer
		diffCmd.Stdout = &diffOut
		diffCmd.Run()

		changes[file] = diffOut.String()
	}

	return changes, nil
}

// generatePDF creates the TSD PDF
func generatePDF(changes map[string]string) error {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(15, 15, 15)
	pdf.AddPage()

	// Title
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(190, 10, "Technical Specification Document (TSD)")
	pdf.Ln(15)

	// Developer & Project Info
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(190, 10, "Project Information")
	pdf.Ln(10)

	headers := []string{"Developer", "Project Name", "Project Status"}
	values := []string{"John Doe", "Inventory System", "In Progress"}

	pdf.SetFont("Arial", "", 12)
	for i := 0; i < len(headers); i++ {
		pdf.CellFormat(63, 10, headers[i], "1", 0, "C", false, 0, "")
	}
	pdf.Ln(-1)
	for i := 0; i < len(values); i++ {
		pdf.CellFormat(63, 10, values[i], "1", 0, "C", false, 0, "")
	}
	pdf.Ln(15)

	// Code Changes Section
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(190, 10, "Code Changes")
	pdf.Ln(10)

	if len(changes) == 0 {
		pdf.SetFont("Arial", "", 12)
		pdf.Cell(190, 10, "No code changes detected.")
		pdf.Ln(10)
	} else {
		for file, diff := range changes {
			// Filename
			pdf.SetFont("Arial", "B", 12)
			pdf.Cell(190, 10, "Filename: "+file)
			pdf.Ln(5)

			// Code diff (inside a bordered multi-line cell)
			pdf.SetFont("Courier", "", 10)
			pdf.MultiCell(190, 6, diff, "1", "L", false)
			pdf.Ln(5)
		}
	}

	// Empty Query Changes Table
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(190, 10, "Query Changes")
	pdf.Ln(10)
	pdf.CellFormat(190, 20, "No query changes detected", "1", 0, "C", false, 0, "")
	pdf.Ln(15)

	// Sign Approval Section
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(190, 10, "Sign Approval")
	pdf.Ln(10)

	signHeaders := []string{"Role", "Name", "Signature"}
	signValues := []string{"Lead Developer", "John Doe", "[Signed]"}

	pdf.SetFont("Arial", "", 12)
	for i := 0; i < len(signHeaders); i++ {
		pdf.CellFormat(63, 10, signHeaders[i], "1", 0, "C", false, 0, "")
	}
	pdf.Ln(-1)
	for i := 0; i < len(signValues); i++ {
		pdf.CellFormat(63, 10, signValues[i], "1", 0, "C", false, 0, "")
	}

	// Save PDF
	fileName := fmt.Sprintf("TSD_%s.pdf", time.Now().Format("20060102_150405"))
	err := pdf.OutputFileAndClose(fileName)
	if err != nil {
		return err
	}

	fmt.Println("TSD PDF generated successfully:", fileName)
	return nil
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run main.go <base_branch> <target_branch>")
		return
	}

	baseBranch := os.Args[1]
	targetBranch := os.Args[2]

	fmt.Println("Fetching Git changes between", baseBranch, "and", targetBranch)
	changes, err := getGitChanges(baseBranch, targetBranch)
	if err != nil {
		fmt.Println("Error fetching Git changes:", err)
		return
	}

	fmt.Println("Generating TSD PDF...")
	err = generatePDF(changes)
	if err != nil {
		fmt.Println("Error generating PDF:", err)
		return
	}
}
