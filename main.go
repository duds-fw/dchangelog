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
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	changes := make(map[string]string)
	files := strings.Split(strings.TrimSpace(out.String()), "\n")

	for _, file := range files {
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
	pdf.SetMargins(10, 10, 10)
	pdf.AddPage()

	// Title
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(190, 10, "Technical Specification Document (TSD)")
<<<<<<< Updated upstream
	pdf.Ln(15)

	// Developer & Project Info
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(190, 10, "Project Information")
	pdf.Ln(10)

	headers := []string{"Developer", "Project Name", "Project Status"}
	values := []string{"John Doe", "Inventory System", "In Progress"}
=======
	pdf.Ln(12)
>>>>>>> Stashed changes

	// Developer & Project Info Table
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(40, 8, "Developer: John Doe")
	pdf.Ln(6)
	pdf.Cell(40, 8, "Project: Inventory System")
	pdf.Ln(6)
	pdf.Cell(40, 8, "Status: In Progress")
	pdf.Ln(12)

	// Code Changes Table
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(190, 8, "Code Changes")
	pdf.Ln(10)

	for file, diff := range changes {
		pdf.SetFont("Arial", "B", 12)
		pdf.Cell(190, 8, "File: "+file)
		pdf.Ln(6)

		pdf.SetFont("Courier", "", 10)
		lines := strings.Split(diff, "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "+") {
				pdf.SetTextColor(0, 128, 0) // Green for additions
			} else if strings.HasPrefix(line, "-") {
				pdf.SetTextColor(255, 0, 0) // Red for deletions
			} else {
				pdf.SetTextColor(0, 0, 0) // Black for unchanged lines
			}
			pdf.Cell(190, 5, line)
			pdf.Ln(4)
		}
		pdf.Ln(6)
	}

	// Query Changes Table (Empty)
	pdf.SetTextColor(0, 0, 0)
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(190, 8, "Query Changes")
	pdf.Ln(10)
	pdf.Cell(190, 30, "(No changes)")
	pdf.Ln(20)

	// Sign Approval Table
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(190, 8, "Sign Approval")
	pdf.Ln(10)
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(190, 8, "Role: Lead Developer")
	pdf.Ln(6)
	pdf.Cell(190, 8, "Name: John Doe")
	pdf.Ln(6)
	pdf.Cell(190, 8, "Signature: [Signed]")

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
