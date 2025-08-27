#!/bin/bash

# Architecture Documentation Generator
# Automates generation and validation of architecture documentation

set -euo pipefail

# Configuration
DOCS_DIR="docs/architecture"
OUTPUT_DIR="${DOCS_DIR}/generated"
MERMAID_CONFIG="${DOCS_DIR}/automation/mermaid-config.json"
PLANTUML_JAR="${DOCS_DIR}/automation/plantuml.jar"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Logging functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check dependencies
check_dependencies() {
    log_info "Checking dependencies..."
    
    local missing_deps=()
    
    # Check for Node.js (for Mermaid CLI)
    if ! command -v node &> /dev/null; then
        missing_deps+=("node")
    fi
    
    # Check for Mermaid CLI
    if ! command -v mmdc &> /dev/null; then
        missing_deps+=("@mermaid-js/mermaid-cli (npm install -g @mermaid-js/mermaid-cli)")
    fi
    
    # Check for Java (for PlantUML)
    if ! command -v java &> /dev/null; then
        missing_deps+=("java")
    fi
    
    # Check for Go (for code analysis)
    if ! command -v go &> /dev/null; then
        missing_deps+=("go")
    fi
    
    if [ ${#missing_deps[@]} -ne 0 ]; then
        log_error "Missing dependencies:"
        for dep in "${missing_deps[@]}"; do
            echo "  - $dep"
        done
        exit 1
    fi
    
    log_success "All dependencies available"
}

# Setup output directories
setup_directories() {
    log_info "Setting up output directories..."
    
    mkdir -p "${OUTPUT_DIR}/diagrams"
    mkdir -p "${OUTPUT_DIR}/metrics"
    mkdir -p "${OUTPUT_DIR}/reports"
    
    log_success "Directories created"
}

# Generate Mermaid diagrams
generate_mermaid_diagrams() {
    log_info "Generating Mermaid diagrams..."
    
    # Find all Mermaid code blocks in documentation
    find "${DOCS_DIR}" -name "*.md" -type f | while read -r file; do
        local basename=$(basename "$file" .md)
        local diagram_count=0
        
        # Extract Mermaid diagrams from markdown files
        awk '
        /^```mermaid$/ { in_diagram = 1; diagram_count++; next }
        /^```$/ && in_diagram { 
            close(output_file)
            in_diagram = 0
            next 
        }
        in_diagram {
            if (!output_file) {
                output_file = "'"${OUTPUT_DIR}/diagrams/${basename}"'-" diagram_count ".mmd"
                print "Extracting diagram to " output_file > "/dev/stderr"
            }
            print > output_file
        }
        END { print diagram_count > "'"${OUTPUT_DIR}/diagrams/${basename}"'.count" }
        ' "$file"
    done
    
    # Convert Mermaid files to SVG
    if [ -d "${OUTPUT_DIR}/diagrams" ] && [ "$(ls -A ${OUTPUT_DIR}/diagrams/*.mmd 2>/dev/null)" ]; then
        for mermaid_file in "${OUTPUT_DIR}/diagrams"/*.mmd; do
            if [ -f "$mermaid_file" ]; then
                local svg_file="${mermaid_file%.mmd}.svg"
                local png_file="${mermaid_file%.mmd}.png"
                
                log_info "Converting $(basename "$mermaid_file") to SVG/PNG..."
                
                # Generate SVG
                if mmdc -i "$mermaid_file" -o "$svg_file" -c "$MERMAID_CONFIG" 2>/dev/null; then
                    log_success "Generated $(basename "$svg_file")"
                else
                    log_warning "Failed to generate $(basename "$svg_file")"
                fi
                
                # Generate PNG for compatibility
                if mmdc -i "$mermaid_file" -o "$png_file" -c "$MERMAID_CONFIG" 2>/dev/null; then
                    log_success "Generated $(basename "$png_file")"
                else
                    log_warning "Failed to generate $(basename "$png_file")"
                fi
            fi
        done
    else
        log_info "No Mermaid diagrams found"
    fi
}

# Generate architecture metrics
generate_metrics() {
    log_info "Generating architecture metrics..."
    
    local metrics_file="${OUTPUT_DIR}/metrics/architecture-metrics.json"
    
    # Analyze Go codebase
    local total_lines=$(find . -name "*.go" -not -path "./vendor/*" -not -path "./.git/*" | xargs wc -l | tail -1 | awk '{print $1}')
    local go_files=$(find . -name "*.go" -not -path "./vendor/*" -not -path "./.git/*" | wc -l)
    local test_files=$(find . -name "*_test.go" -not -path "./vendor/*" | wc -l)
    local service_files=$(find ./services -name "*.go" 2>/dev/null | wc -l || echo 0)
    local template_files=$(find ./templates -name "*.templ" 2>/dev/null | wc -l || echo 0)
    
    # Generate metrics JSON
    cat > "$metrics_file" << EOF
{
  "timestamp": "$(date -u +"%Y-%m-%dT%H:%M:%SZ")",
  "codebase": {
    "total_lines": $total_lines,
    "go_files": $go_files,
    "test_files": $test_files,
    "test_coverage_ratio": $(echo "scale=2; $test_files / ($go_files - $test_files)" | bc -l 2>/dev/null || echo "0.00"),
    "service_files": $service_files,
    "template_files": $template_files
  },
  "architecture": {
    "layers": ["presentation", "business", "data"],
    "services": ["task", "project", "calendar", "oauth", "home"],
    "external_integrations": ["google_calendar", "google_oauth"],
    "database": "boltdb",
    "frontend": "htmx_alpine_tailwind"
  },
  "documentation": {
    "total_docs": $(find "${DOCS_DIR}" -name "*.md" | wc -l),
    "adr_count": $(find "${DOCS_DIR}/adrs" -name "*.md" -not -name "template.md" 2>/dev/null | wc -l || echo 0),
    "diagram_count": $(find "${OUTPUT_DIR}/diagrams" -name "*.svg" 2>/dev/null | wc -l || echo 0)
  }
}
EOF
    
    log_success "Generated architecture metrics: $(basename "$metrics_file")"
}

# Validate documentation
validate_documentation() {
    log_info "Validating documentation..."
    
    local validation_report="${OUTPUT_DIR}/reports/validation-report.md"
    local issues_found=0
    
    cat > "$validation_report" << EOF
# Documentation Validation Report

Generated: $(date -u +"%Y-%m-%d %H:%M:%S UTC")

## Validation Results

EOF
    
    # Check for required documentation files
    local required_files=(
        "README.md"
        "01-system-context.md"
        "02-container-architecture.md"
        "03-component-architecture.md"
        "04-data-architecture.md"
        "05-security-architecture.md"
        "06-quality-attributes.md"
        "adrs/template.md"
    )
    
    echo "### Required Files" >> "$validation_report"
    echo "" >> "$validation_report"
    
    for file in "${required_files[@]}"; do
        if [ -f "${DOCS_DIR}/$file" ]; then
            echo "- ✅ $file" >> "$validation_report"
        else
            echo "- ❌ $file (MISSING)" >> "$validation_report"
            ((issues_found++))
        fi
    done
    
    echo "" >> "$validation_report"
    
    # Check for broken internal links
    echo "### Link Validation" >> "$validation_report"
    echo "" >> "$validation_report"
    
    find "${DOCS_DIR}" -name "*.md" -type f | while read -r file; do
        local broken_links=0
        
        # Extract markdown links
        grep -oE '\[.*\]\([^)]+\)' "$file" | while read -r link; do
            local url=$(echo "$link" | sed 's/.*(\([^)]*\)).*/\1/')
            
            # Check internal links (relative paths)
            if [[ "$url" =~ ^[^:]+\.md$ ]] || [[ "$url" =~ ^\./ ]]; then
                local target_file
                if [[ "$url" =~ ^\./ ]]; then
                    target_file="$(dirname "$file")/$url"
                else
                    target_file="${DOCS_DIR}/$url"
                fi
                
                if [ ! -f "$target_file" ]; then
                    echo "- ❌ Broken link in $(basename "$file"): $url" >> "$validation_report"
                    ((broken_links++))
                fi
            fi
        done
        
        if [ $broken_links -eq 0 ]; then
            echo "- ✅ $(basename "$file"): No broken links" >> "$validation_report"
        fi
    done
    
    echo "" >> "$validation_report"
    
    # Check ADR numbering
    echo "### ADR Validation" >> "$validation_report"
    echo "" >> "$validation_report"
    
    local adr_files=($(find "${DOCS_DIR}/adrs" -name "[0-9]*.md" | sort))
    local expected_number=1
    
    for adr_file in "${adr_files[@]}"; do
        local filename=$(basename "$adr_file")
        local adr_number=$(echo "$filename" | grep -oE '^[0-9]+')
        
        if [ "$adr_number" -ne $expected_number ]; then
            echo "- ❌ ADR numbering gap: Expected $expected_number, found $adr_number" >> "$validation_report"
            ((issues_found++))
        else
            echo "- ✅ ADR-$(printf "%04d" $adr_number): Correct numbering" >> "$validation_report"
        fi
        
        ((expected_number++))
    done
    
    echo "" >> "$validation_report"
    echo "## Summary" >> "$validation_report"
    echo "" >> "$validation_report"
    
    if [ $issues_found -eq 0 ]; then
        echo "✅ **All validation checks passed**" >> "$validation_report"
        log_success "Documentation validation passed"
    else
        echo "❌ **$issues_found issue(s) found**" >> "$validation_report"
        log_warning "Documentation validation found $issues_found issue(s)"
    fi
    
    log_success "Generated validation report: $(basename "$validation_report")"
}

# Generate summary report
generate_summary() {
    log_info "Generating summary report..."
    
    local summary_file="${OUTPUT_DIR}/architecture-summary.md"
    
    cat > "$summary_file" << EOF
# Architecture Documentation Summary

Generated: $(date -u +"%Y-%m-%d %H:%M:%S UTC")

## Documentation Status

### Core Architecture Documents
- [x] System Context
- [x] Container Architecture  
- [x] Component Architecture
- [x] Data Architecture
- [x] Security Architecture
- [x] Quality Attributes

### Architecture Decision Records (ADRs)
$(find "${DOCS_DIR}/adrs" -name "[0-9]*.md" | sort | while read -r adr; do
    local title=$(grep "^# ADR-" "$adr" | head -1 | sed 's/^# //')
    local status=$(grep "^Accepted\|^Proposed\|^Rejected" "$adr" | head -1)
    echo "- [x] $title ($status)"
done)

### Generated Artifacts
- Diagrams: $(find "${OUTPUT_DIR}/diagrams" -name "*.svg" 2>/dev/null | wc -l || echo 0) SVG files
- Metrics: Architecture metrics available
- Validation: $(if [ -f "${OUTPUT_DIR}/reports/validation-report.md" ]; then echo "Completed"; else echo "Not run"; fi)

### Codebase Overview
$(if [ -f "${OUTPUT_DIR}/metrics/architecture-metrics.json" ]; then
    echo "- Total Lines of Code: $(jq -r '.codebase.total_lines' "${OUTPUT_DIR}/metrics/architecture-metrics.json")"
    echo "- Go Files: $(jq -r '.codebase.go_files' "${OUTPUT_DIR}/metrics/architecture-metrics.json")"
    echo "- Test Files: $(jq -r '.codebase.test_files' "${OUTPUT_DIR}/metrics/architecture-metrics.json")"
    echo "- Services: $(jq -r '.codebase.service_files' "${OUTPUT_DIR}/metrics/architecture-metrics.json")"
    echo "- Templates: $(jq -r '.codebase.template_files' "${OUTPUT_DIR}/metrics/architecture-metrics.json")"
else
    echo "- Metrics not available"
fi)

## Quick Links
- [README](README.md)
- [System Context](01-system-context.md)
- [Container Architecture](02-container-architecture.md)
- [Component Architecture](03-component-architecture.md)
- [Data Architecture](04-data-architecture.md)
- [Security Architecture](05-security-architecture.md)
- [Quality Attributes](06-quality-attributes.md)
- [ADR Directory](adrs/)
- [Generated Artifacts](generated/)

EOF
    
    log_success "Generated summary: $(basename "$summary_file")"
}

# Main execution
main() {
    log_info "Starting architecture documentation generation..."
    echo
    
    check_dependencies
    setup_directories
    generate_mermaid_diagrams
    generate_metrics
    validate_documentation
    generate_summary
    
    echo
    log_success "Architecture documentation generation completed!"
    log_info "Generated files available in: $OUTPUT_DIR"
}

# Run if called directly
if [ "${BASH_SOURCE[0]}" = "${0}" ]; then
    main "$@"
fi