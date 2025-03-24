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
	cmd := exec.Command("git", "diff", baseBranch+".."+targetBranch, "--name-only")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	changes := make(map[string]string)
	files := strings.Split(out.String(), "\n")

	for _, file := range files {
		if strings.TrimSpace(file) == "" {
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

// generatePDF creates or updates the TSD PDF
func generatePDF(changes map[string]string) error {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Title
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(190, 10, "Technical Specification Document (TSD)")
	pdf.Ln(15)

	// Developer & Project Info Table
	pdf.SetFont("Arial", "", 12)
	pdf.CellFormat(63, 10, "Developer: John Doe", "1", 0, "C", false, 0, "")
	pdf.CellFormat(63, 10, "Project: Inventory System", "1", 0, "C", false, 0, "")
	pdf.CellFormat(64, 10, "Status: In Progress", "1", 1, "C", false, 0, "")
	pdf.Ln(10)

	// Code Changes Table
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(190, 10, "Code Changes")
	pdf.Ln(10)

	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(190, 10, "File Changes", "1", 1, "C", false, 0, "")

	for file, diff := range changes {
		pdf.SetFont("Arial", "B", 12)
		pdf.CellFormat(190, 10, "Filename: "+file, "1", 1, "L", false, 0, "")
		pdf.SetFont("Courier", "", 10)

		pdf.CellFormat(190, 5, "", "1", 1, "L", false, 0, "") // Top border
		pdf.SetLeftMargin(10)
		lines := strings.Split(diff, "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "+") {
				pdf.SetTextColor(0, 128, 0) // Green for additions
			} else if strings.HasPrefix(line, "-") {
				pdf.SetTextColor(255, 0, 0) // Red for deletions
			} else {
				pdf.SetTextColor(0, 0, 0) // Default black
			}
			pdf.MultiCell(190, 5, line, "LR", "L", false)
		}
		pdf.CellFormat(190, 5, "", "1", 1, "L", false, 0, "") // Bottom border
		pdf.SetTextColor(0, 0, 0)                             // Reset color
	}

	pdf.Ln(10)

	// Empty Query Changes Table
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(190, 10, "Query Changes")
	pdf.Ln(10)
	pdf.CellFormat(190, 30, "", "1", 1, "C", false, 0, "")
	pdf.Ln(10)

	// Sign Approval Table
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(190, 10, "Sign Approval")
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 12)
	pdf.CellFormat(95, 10, "Role: Lead Developer", "1", 0, "C", false, 0, "")
	pdf.CellFormat(95, 10, "Name: John Doe", "1", 1, "C", false, 0, "")
	pdf.CellFormat(190, 20, "Signature: [Signed]", "1", 1, "C", false, 0, "")

	// Save PDF
	fileName := fmt.Sprintf("tsd_%s.pdf", time.Now().Format("20060102_150405"))
	return pdf.OutputFileAndClose(fileName)
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

	fmt.Println("TSD PDF generated successfully!")
}
