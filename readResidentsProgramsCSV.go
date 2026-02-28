package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// The Resident data type
type Resident struct {
	residentID     int
	firstname      string
	lastname       string
	rol            []string // resident rank order list
	matchedProgram string   // will be "" for unmatched resident
}

// The Program data type
type Program struct {
	programID         string
	name              string
	nPositions        int   // number of positions available (quota)
	rol               []int // program rank order list
	selectedResidents *residentHeap
}

// Parse a resident's ROL
func parseRol(s string) []string {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "[")
	s = strings.TrimSuffix(s, "]")
	if s == "" {
		return []string{}
	}
	parts := strings.Split(s, ",")
	for i, part := range parts {
		parts[i] = strings.TrimSpace(part)
	}
	return parts
}

// Parse a program's ROL
func parseIntRol(s string) []int {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "[")
	s = strings.TrimSuffix(s, "]")
	if s == "" {
		return []int{}
	}
	parts := strings.Split(s, ",")
	var ints []int
	for _, part := range parts {
		pid, _ := strconv.Atoi(strings.TrimSpace(part))
		ints = append(ints, pid)
	}
	return ints
}

// ReadCSV reads a CSV file into a map of Resident
func ReadResidentsCSV(filename string) (map[int]*Resident, error) {

	// map to store residents by ID
	residents := make(map[int]*Resident)

	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("unable to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Read all records
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error reading CSV: %w", err)
	}

	// Skip header if present (assuming it is)
	for i, record := range records {
		if i == 0 && record[0] == "id" {
			continue
		}
		if len(record) < 4 {
			return nil, fmt.Errorf("invalid record at line %d: %v", i+1, record)
		}

		// Parse ID
		id, err := strconv.Atoi(record[0])
		if err != nil {
			return nil, fmt.Errorf("invalid ID at line %d: %w", i+1, err)
		}

		if _, exists := residents[id]; exists {
			fmt.Println(id)
		}

		residents[id] = &Resident{
			residentID:     id,
			firstname:      record[1],
			lastname:       record[2],
			rol:            parseRol(record[3]),
			matchedProgram: "",
		}
	}

	return residents, nil
}

// reads a CSV file into a map of Program
func ReadProgramsCSV(filename string) (map[string]*Program, error) {

	// map to store programs by ID
	programs := make(map[string]*Program)

	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("unable to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Read all records
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error reading CSV: %w", err)
	}

	// Skip header if present (assuming it is)
	for i, record := range records {
		if i == 0 && record[0] == "id" {
			continue
		}
		if len(record) < 4 {
			return nil, fmt.Errorf("invalid record at line %d: %v", i+1, record)
		}

		// Parse number of positions
		np, err := strconv.Atoi(record[2])
		if err != nil {
			return nil, fmt.Errorf("invalid number at line %d: %w", i+1, err)
		}

		programs[record[0]] = &Program{
			programID:  record[0],
			name:       record[1],
			nPositions: np,
			rol:        parseIntRol(record[3]),
		}

	}

	programs[""] = &Program{
		programID:  "XXX",
		name:       "NOT_MATCHED",
		nPositions: 0,
		rol:        []int{0},
	}

	return programs, nil
}

// Example usage
func main() {

	// read residents
	residents, err := ReadResidentsCSV("residents4000.csv")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	for _, p := range residents {
		fmt.Printf("ID: %d, Name: %s %s, Rol: %v\n", p.residentID, p.firstname, p.lastname, p.rol)
	}

	programs, err := ReadProgramsCSV("programs4000.csv")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	for _, p := range programs {
		fmt.Printf("ID: %s, Name: %s, Number of pos: %d, Number of applicants: %d\n", p.programID, p.name, p.nPositions, len(p.rol))
	}

	//TODO: uncomment this line fmt.Printf("\nNMD: %v", programs["NMD"])
	fmt.Printf("\nMMI: %v", programs["MMI"])

	start := time.Now()

	//fmt.Println("\nright before McVittieWilson")
	McVittieWilson(residents, programs)
	//fmt.Println("right after McVittieWilson")

	end := time.Now()

	var unmatched int
	fmt.Println("\n\nlastname, firstname, residentID, programID, name")
	for _, r := range residents {
		fmt.Printf("\n%s, %s, %d, %s, %s, ", r.lastname, r.firstname, r.residentID, programs[r.matchedProgram].programID, programs[r.matchedProgram].name)
		if r.matchedProgram == "" {
			unmatched++
		}
	}

	/*
		fmt.Println("\n\nid, name, numberofPos, rol")
		for _, p := range programs {
			fmt.Printf("\n%s, %s, %d, %v", p.programID, p.name, p.nPositions, p.rol)
		}
	*/

	fmt.Printf("\n\nExecution Time: %s", end.Sub(start))

}
