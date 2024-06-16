package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/smtp"
	"strings"
	"sync"
	"time"
)

type Report struct {
	Status       string      `json:"status"`
	LastChecked  time.Time   `json:"last_checked"`
	StatusCode   int         `json:"status_code"`
	ResponseTime string      `json:"response_time"`
	Issues       string      `json:"issues"`
	Headers      http.Header `json:"headers"`
}

type Option struct {
	CheckSSL       bool        `json:"check_ssl"`
	RequestHeaders http.Header `json:"request_headers"`
}

type Service struct {
	Name        string   `json:"name"`
	URL         string   `json:"url"`
	Type        string   `json:"type"`
	Port        string   `json:"port"`
	Option      Option   `json:"option"`
	LastReports []Report `json:"last_reports"`
	Report      Report   `json:"report"`
}

var services = make(map[string]Service)
var mu sync.Mutex

func main() {
	http.HandleFunc("/service/add", addServiceHandler)
	http.HandleFunc("/service/check", checkServiceHandler)
	http.HandleFunc("/service/status", statusHandler)
	fmt.Println("Listening on http://localhost:8080")
	_ = http.ListenAndServe(":8080", nil)
}

func addServiceHandler(w http.ResponseWriter, r *http.Request) {
	var service Service
	if err := json.NewDecoder(r.Body).Decode(&service); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mu.Lock()
	services[service.Name] = service
	mu.Unlock()

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(service)
}

func checkServiceHandler(w http.ResponseWriter, r *http.Request) {
	serviceName := r.URL.Query().Get("name")
	mu.Lock()
	service, exists := services[serviceName]
	mu.Unlock()

	if !exists {
		http.Error(w, "Service not found", http.StatusNotFound)
		return
	}

	switch strings.ToLower(service.Type) {
	case "http", "https":
		checkHTTPService(&service)
	case "tcp":
		checkTCPService(&service)
	case "smtp":
		checkSMTPService(&service)
	// Add cases for database, cache, email smtp, etc.
	default:
		http.Error(w, "Unsupported service type", http.StatusBadRequest)
		return
	}

	mu.Lock()
	services[service.Name] = service
	mu.Unlock()

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(service)
}

func checkSMTPService(service *Service) {
	start := time.Now()
	service.Report.LastChecked = start

	// SMTP connectivity check
	conn, err := smtp.Dial(fmt.Sprintf("%s:%s", service.URL, service.Port))
	if err != nil {
		service.Report.Status = "Failed"
		service.Report.StatusCode = 0
		service.Report.ResponseTime = ""
		service.Report.Issues = fmt.Sprintf("SMTP connection failed: %v", err)
		service.LastReports = append(service.LastReports, service.Report)
		return
	}
	_ = conn.Close()
	service.Report.Status = "Success"
	service.Report.ResponseTime = fmt.Sprintf("%s", time.Since(start))
	service.Report.StatusCode = 200 // HTTP OK equivalent

	// DNS records check
	issues := checkDNSRecords(service.URL)
	if issues != "" {
		service.Report.Issues = issues
	} else {
		service.Report.Issues = "All DNS checks passed"
	}

	service.LastReports = append(service.LastReports, service.Report)
}

func checkDNSRecords(domain string) string {
	var issues strings.Builder

	// Check MX records
	mxRecords, err := net.LookupMX(domain)
	if err != nil || len(mxRecords) == 0 {
		issues.WriteString("Missing MX records; ")
	}

	// Check SPF records
	spfRecords, err := net.LookupTXT(domain)
	hasSPF := false
	if err == nil {
		for _, txt := range spfRecords {
			if strings.HasPrefix(txt, "v=spf1") {
				hasSPF = true
				break
			}
		}
	}
	if !hasSPF {
		issues.WriteString("Missing SPF records; ")
	}

	// Check DKIM records
	dkimSelector := "default._domainkey." + domain
	_, err = net.LookupTXT(dkimSelector)
	if err != nil {
		issues.WriteString("Missing DKIM records; ")
	}

	return issues.String()
}

func statusHandler(w http.ResponseWriter, _ *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(services)
}

func checkHTTPService(service *Service) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	client.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
	}

	start := time.Now()
	req, err := http.NewRequest("HEAD", service.URL, nil)
	if err == nil {
		for key, values := range service.Option.RequestHeaders {
			for _, value := range values {
				req.Header.Add(key, value)
			}
		}

		resp, err := client.Do(req)
		service.Report.LastChecked = start

		if err != nil || resp.StatusCode >= 400 {
			service.Report.Status = "Failed"
			service.Report.StatusCode = 0
			service.Report.ResponseTime = ""
			service.Report.Issues = fmt.Sprintf("SSL check failed: %v", err)

			client.Transport = &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			}
			resp, err = client.Do(req)
			if err != nil || resp.StatusCode >= 400 {
				service.Report.Status = "Failed"
				service.Report.StatusCode = 0
				service.Report.ResponseTime = ""
				service.Report.Issues += fmt.Sprintf("; Retry without SSL failed: %v", err)
				service.LastReports = append(service.LastReports, service.Report)
				return
			} else {
				service.Report.Issues += "; SSL issue detected, but passed without SSL verification"
			}
		}

		defer func() {
			_ = resp.Body.Close()
		}()
		service.Report.Status = "Success"
		service.Report.Headers = resp.Header
		service.Report.StatusCode = resp.StatusCode
		service.Report.ResponseTime = fmt.Sprintf("%s", time.Since(start))
		service.Report.Issues += checkSecurityHeaders(resp)
		service.LastReports = append(service.LastReports, service.Report)
	} else {
		service.Report.Status = "Failed"
		service.Report.StatusCode = 0
		service.Report.ResponseTime = ""
		service.Report.Issues = err.Error()
		service.LastReports = append(service.LastReports, service.Report)
	}
}

func checkTCPService(service *Service) {
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(service.URL, service.Port), 10*time.Second)
	service.Report.LastChecked = time.Now()

	if err != nil {
		service.Report.Status = "Failed"
		service.Report.StatusCode = 0
		service.Report.ResponseTime = ""
		service.Report.Issues = err.Error()
		service.LastReports = append(service.LastReports, service.Report)
		return
	}
	_ = conn.Close()
	service.Report.Status = "Success"
	service.Report.StatusCode = 200 // HTTP OK equivalent
	service.Report.ResponseTime = ""
	service.Report.Issues = ""
	service.LastReports = append(service.LastReports, service.Report)
}

func checkSecurityHeaders(resp *http.Response) string {
	securityHeaders := []string{
		"Strict-Transport-Security",
		"Content-Security-Policy",
		"X-Content-Type-Options",
		"X-Frame-Options",
		"X-XSS-Protection",
		"Referrer-Policy",
		"Permissions-Policy",
	}

	var issues strings.Builder
	for _, header := range securityHeaders {
		if value := resp.Header.Get(header); value == "" {
			issues.WriteString(fmt.Sprintf("%s: MISSING; ", header))
		}
	}
	return issues.String()
}
