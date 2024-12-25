// Package main provides a web-based port scanner application
package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// ScanResult represents the result of a port scan
type ScanResult struct {
	IP   string
	Port int
	Open bool
}

// PageData represents the data passed to HTML templates
type PageData struct {
	Results string
	Error   string
}

// scanPort attempts to connect to a specific IP and port
// Returns a ScanResult indicating if the port is open
func scanPort(ip string, port int, timeout time.Duration) ScanResult {
	target := fmt.Sprintf("%s:%d", ip, port)

	// First try a quick connection to see if port is definitely closed
	conn, err := net.DialTimeout("tcp", target, 500*time.Millisecond)
	if err == nil {
		conn.Close()
		log.Printf("Found open port at %s\n", target)
		return ScanResult{IP: ip, Port: port, Open: true}
	}

	// If first attempt failed, try again with longer timeout
	conn, err = net.DialTimeout("tcp", target, timeout)
	if err != nil {
		if strings.Contains(err.Error(), "refused") {
			// Port is definitely closed
			return ScanResult{IP: ip, Port: port, Open: false}
		}
		// Log other types of errors for debugging
		log.Printf("Error scanning %s: %v\n", target, err)
		return ScanResult{IP: ip, Port: port, Open: false}
	}
	defer conn.Close()

	log.Printf("Found open port at %s\n", target)
	return ScanResult{IP: ip, Port: port, Open: true}
}

// worker function processes ports from a channel
// Uses a WaitGroup to signal when it's done
func worker(ip string, ports <-chan int, results chan<- ScanResult, wg *sync.WaitGroup) {
	defer wg.Done()
	for port := range ports {
		// Add small delay between scans to prevent overwhelming the network
		time.Sleep(10 * time.Millisecond)
		results <- scanPort(ip, port, 5*time.Second)
	}
}

// getPortsToScan returns the list of ports to scan based on the selected option
func getPortsToScan(portOption, portList string, startPort, endPort int) ([]int, error) {
	var ports []int

	switch portOption {
	case "common":
		// Common service ports
		commonPorts := []int{20, 21, 22, 23, 25, 53, 80, 110, 143, 443, 465, 587, 993, 995, 3306, 3389, 5432, 8080}
		return commonPorts, nil

	case "wellknown":
		// Well-known ports (1-1024)
		for i := 1; i <= 1024; i++ {
			ports = append(ports, i)
		}
		return ports, nil

	case "custom":
		// Validate custom range
		if startPort < 1 || endPort > 65535 || startPort > endPort {
			return nil, fmt.Errorf("invalid port range: must be between 1-65535 and start port must be less than end port")
		}
		for i := startPort; i <= endPort; i++ {
			ports = append(ports, i)
		}
		return ports, nil

	case "specific":
		// Parse comma-separated port list
		if portList == "" {
			return nil, fmt.Errorf("no ports specified")
		}
		portStrings := strings.Split(portList, ",")
		for _, portStr := range portStrings {
			port, err := strconv.Atoi(strings.TrimSpace(portStr))
			if err != nil {
				return nil, fmt.Errorf("invalid port number: %s", portStr)
			}
			if port < 1 || port > 65535 {
				return nil, fmt.Errorf("port %d is out of valid range (1-65535)", port)
			}
			ports = append(ports, port)
		}
		return ports, nil

	default:
		return nil, fmt.Errorf("invalid port option")
	}
}

// performScan executes the port scanning operation
// Returns scan results as a formatted string
func performScan(ip string, portOption string, portList string, startPort, endPort int) (string, error) {
	var buffer bytes.Buffer

	// Validate IP
	if net.ParseIP(ip) == nil {
		return "", fmt.Errorf("invalid IP address format")
	}

	// Get ports to scan
	ports, err := getPortsToScan(portOption, portList, startPort, endPort)
	if err != nil {
		return "", err
	}

	// Log scan start
	log.Printf("Starting scan for IP: %s, Port option: %s\n", ip, portOption)

	foundOpenPorts := false
	totalScanned := 0

	// Calculate optimal number of workers
	numWorkers := 10 // Base number of workers
	if len(ports) > 100 {
		numWorkers = 20
	}
	if len(ports) > 1000 {
		numWorkers = 30
	}

	buffer.WriteString(fmt.Sprintf("\nScanning IP: %s\n", ip))
	buffer.WriteString(fmt.Sprintf("Total ports to scan: %d\n", len(ports)))

	// Create channels for ports and results
	portsChan := make(chan int, numWorkers)
	results := make(chan ScanResult, numWorkers)

	// Create wait group for workers
	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(ip, portsChan, results, &wg)
	}

	// Send ports to workers
	go func() {
		for _, port := range ports {
			portsChan <- port
			totalScanned++
		}
		close(portsChan)
	}()

	// Start a goroutine to close results channel when all workers are done
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect and process results
	for result := range results {
		if result.Open {
			foundOpenPorts = true
			output := fmt.Sprintf("Port %d: Open\n", result.Port)
			buffer.WriteString(output)
			log.Printf("Found open port: %s:%d\n", result.IP, result.Port)
		}
	}

	// Add summary to the output
	summary := fmt.Sprintf("\nScan Summary:\n")
	summary += fmt.Sprintf("IP Address: %s\n", ip)
	summary += fmt.Sprintf("Total ports scanned: %d\n", totalScanned)
	if !foundOpenPorts {
		summary += "\nNo open ports were found. This could mean:\n"
		summary += "1. The target IP is not responding\n"
		summary += "2. No ports are open in the specified range\n"
		summary += "3. A firewall might be blocking the scan\n"
		summary += "\nTroubleshooting tips:\n"
		summary += "- Try scanning common ports first\n"
		summary += "- Check if Windows Defender or antivirus is blocking the scan\n"
		summary += "- Try scanning your local machine (127.0.0.1) first\n"
	}
	buffer.WriteString(summary)

	return buffer.String(), nil
}

// HTTP handler for the index page
func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, &PageData{})
}

// HTTP handler for the scan operation
func scanHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Get form values
	ip := r.FormValue("ip")
	portOption := r.FormValue("portOption")
	portList := r.FormValue("portList")

	var startPort, endPort int
	var err error

	// Parse custom port range if selected
	if portOption == "custom" {
		startPort, err = strconv.Atoi(r.FormValue("startPort"))
		if err != nil {
			startPort = 1
		}
		endPort, err = strconv.Atoi(r.FormValue("endPort"))
		if err != nil {
			endPort = 1024
		}
	}

	// Perform the scan
	results, err := performScan(ip, portOption, portList, startPort, endPort)

	// Prepare template data
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	data := &PageData{}

	if err != nil {
		data.Error = err.Error()
	} else {
		data.Results = results
	}

	tmpl.Execute(w, data)
}

func main() {
	// Setup logging
	logFile, err := os.OpenFile("port_scan.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Error opening log file: %v\n", err)
		return
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	// Setup HTTP routes
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/scan", scanHandler)

	// Start the web server
	fmt.Println("Server starting on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
