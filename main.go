package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/jung-kurt/gofpdf"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"golang.org/x/exp/slices"
)

// Config represents the structure of the JSON config file
type Config struct {
	PR struct {
		Link string `json:"link"`
	} `json:"pr"`
	Jira struct {
		Link        string `json:"link"`
		Title       string `json:"title"`
		Description string `json:"description"`
	} `json:"jira"`
	Developer struct {
		Name  string `json:"name"`
		Title string `json:"title"`
	} `json:"developer"`
	Project struct {
		Title string `json:"title"`
		Value string `json:"value"`
	} `json:"project"`
	Status struct {
		Title string `json:"title"`
	} `json:"status"`
	SignApproval struct {
		Name string `json:"name"`
		Role string `json:"role"`
	} `json:"sign_approval"`
}

// LoadConfig reads the JSON configuration file
func LoadConfig(filename string) (*Config, error) {
	file, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(file, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

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
func generatePDF(config *Config, changes map[string]string) error {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Title
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(190, 10, "Technical Specification Document (TSD)")
	pdf.Ln(15)

	// Title
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(190, 10, config.Jira.Title)
	pdf.Ln(7.5)

	pdf.SetFont("Arial", "", 12)
	pdf.Cell(190, 10, "Link JIRA: "+config.Jira.Link)
	pdf.Ln(7.5)

	pdf.SetFont("Arial", "", 12)
	pdf.Cell(190, 10, "Link PR: "+config.PR.Link)
	pdf.Ln(7.5)

	pdf.SetFont("Arial", "", 12)
	pdf.Cell(190, 10, "Task Description: "+config.Jira.Description)
	pdf.Ln(15)

	// Developer & Project Info Table
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(63, 10, "Developer", "1", 0, "C", false, 0, "")
	pdf.CellFormat(63, 10, config.Project.Title, "1", 0, "C", false, 0, "")
	pdf.CellFormat(64, 10, "Status", "1", 1, "C", false, 0, "")
	pdf.SetFont("Arial", "", 12)
	pdf.CellFormat(63, 10, config.Developer.Name, "1", 0, "C", false, 0, "")
	pdf.CellFormat(63, 10, config.Project.Value, "1", 0, "C", false, 0, "")
	pdf.CellFormat(64, 10, config.Status.Title, "1", 1, "C", false, 0, "")
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

	// Sign Approval Table
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(190, 10, "Sign Approval")
	pdf.Ln(10)
	pdf.SetFont("Arial", "", 12)
	pdf.CellFormat(95, 10, "Role: "+config.SignApproval.Role, "1", 0, "C", false, 0, "")
	pdf.CellFormat(95, 10, "Name: "+config.SignApproval.Name, "1", 1, "C", false, 0, "")
	pdf.CellFormat(190, 20, "Signature: [Signed]", "1", 1, "C", false, 0, "")

	// Save PDF
	fileName := fmt.Sprintf("tsd_%s.pdf", time.Now().Format("20060102_150405"))
	return pdf.OutputFileAndClose(fileName)
}

func mergePDFs(folderPath, outputFileName string) error {
	files, err := filepath.Glob(filepath.Join(folderPath, "*.pdf"))
	if err != nil {
		return err
	}

	if len(files) == 0 {
		return fmt.Errorf("no PDF files found in the directory")
	}

	err = api.MergeCreateFile(files, outputFileName, false, nil)
	if err != nil {
		return fmt.Errorf("error merging PDFs: %v", err)
	}

	fmt.Println("Merged PDF saved as:", outputFileName)
	return nil
}

func main() {
	usage := []string{"generate", "merge"}
	args1 := os.Args[1]

	baseBranch := flag.String("dest", "", "Base branch for Git diff")
	targetBranch := flag.String("src", "", "Target branch for Git diff")
	folderPath := flag.String("folder", "", "Folder containing PDFs to merge")
	configFile := flag.String("config", "", "Config json")
	flag.CommandLine.Parse(os.Args[2:])

	cliScript := "dchangelog"

	// if len(os.Args) < 3 {
	// 	fmt.Println("Usage:")
	// 	fmt.Println("  Generate TSD: go run main.go generate <json_file> <base_branch>  <target_branch>")
	// 	fmt.Println("  Merge PDFs: go run main.go merge <folder_path>")
	// 	return
	// }

	if slices.Contains(usage, args1) {
		if args1 == "generate" {

			if *configFile != "" && *baseBranch != "" && *targetBranch != "" {
				fmt.Println("Load config ", *configFile, " ...")
				config, err := LoadConfig(*configFile)
				if err != nil {
					log.Fatal("Error loading config:", err)
					return
				}

				fmt.Println("Fetching Git changes between", *baseBranch, "and", *targetBranch)
				changes, err := getGitChanges(*baseBranch, *targetBranch)
				if err != nil {
					log.Fatal("Error fetching Git changes:", err)
					return
				}
				fmt.Println("Generating TSD PDF...")
				err = generatePDF(config, changes)
				if err != nil {
					log.Fatal("Error generating PDF:", err)
					return
				}
				fmt.Println("TSD PDF generated successfully!")
			} else {
				log.Fatalf("Error Usage\nGenerate TSD: %s generate --config=<jsonfile> --dest=<base_branch> --src=<target_branch>", cliScript)
			}

		}
		if args1 == "merge" {
			if *folderPath != "" {
				fmt.Println("Merging PDFs in folder:", *folderPath)
				err := mergePDFs(*folderPath, "merged_output.pdf")
				if err != nil {
					fmt.Println("Error merging PDFs:", err)
				}
				return
			} else {
				log.Fatalf("Error Usage\nMerge PDF: %s merge --folder=<folder> --dest=<base_branch> --src=<target_branch>", cliScript)
			}
		}
	} else {
		log.Fatalf("Error Usage\nGenerate TSD: %s generate --config=<jsonfile> --dest=<base_branch> --src=<target_branch>\nMerge PDF: %s merge --folder=<folder> --dest=<base_branch> --src=<target_branch>", cliScript, cliScript)
	}

}
