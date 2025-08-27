#!/bin/bash

# Documentation Validation Script
# Validates architecture documentation for completeness and consistency

set -euo pipefail

# Configuration
DOCS_DIR="docs/architecture"
SCRIPT_DIR="$(dirname "${BASH_SOURCE[0]}")"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Counters
TOTAL_CHECKS=0
PASSED_CHECKS=0
FAILED_CHECKS=0

# Logging functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[✓]${NC} $1"
    ((PASSED_CHECKS++))
}

log_error() {
    echo -e "${RED}[✗]${NC} $1"
    ((FAILED_CHECKS++))
}

log_warning() {
    echo -e "${YELLOW}[!]${NC} $1"
}

# Increment total checks counter
check() {
    ((TOTAL_CHECKS++))
}

# Check if required files exist
check_required_files() {
    log_info "Checking required documentation files..."
    
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
    
    for file in "${required_files[@]}"; do
        check
        if [ -f "${DOCS_DIR}/$file" ]; then
            log_success "Required file exists: $file"
        else
            log_error "Missing required file: $file"
        fi
    done
}

# Check ADR structure and numbering
check_adr_structure() {
    log_info "Checking ADR structure and numbering..."
    
    # Check if ADR template exists
    check
    if [ -f "${DOCS_DIR}/adrs/template.md" ]; then
        log_success "ADR template exists"
    else
        log_error "ADR template missing"
    fi
    
    # Check ADR numbering sequence
    local adr_files=($(find "${DOCS_DIR}/adrs" -name "[0-9]*.md" | sort 2>/dev/null || true))
    local expected_number=1
    
    for adr_file in "${adr_files[@]}"; do
        check
        local filename=$(basename "$adr_file")
        local adr_number=$(echo "$filename" | grep -oE '^[0-9]+' | sed 's/^0*//')
        
        if [ "$adr_number" -eq $expected_number ]; then
            log_success "ADR-$(printf "%04d" $adr_number) numbering correct"
        else
            log_error "ADR numbering gap: Expected $expected_number, found $adr_number"
        fi
        
        # Check ADR format
        check
        if grep -q "^# ADR-[0-9]" "$adr_file"; then
            log_success "ADR-$(printf "%04d" $adr_number) has correct title format"
        else
            log_error "ADR-$(printf "%04d" $adr_number) missing proper title format"
        fi
        
        # Check required sections
        local required_sections=("Status" "Context" "Decision" "Consequences")
        for section in "${required_sections[@]}"; do
            check
            if grep -q "^## $section" "$adr_file"; then
                log_success "ADR-$(printf "%04d" $adr_number) has $section section"
            else
                log_error "ADR-$(printf "%04d" $adr_number) missing $section section"
            fi
        done
        
        ((expected_number++))
    done
}

# Check for broken internal links
check_internal_links() {
    log_info "Checking internal links..."
    
    find "${DOCS_DIR}" -name "*.md" -type f | while read -r file; do
        local relative_path=${file#${DOCS_DIR}/}
        
        # Extract markdown links [text](url)
        grep -oE '\[.*\]\([^)]+\)' "$file" 2>/dev/null | while read -r link; do
            local url=$(echo "$link" | sed 's/.*(\([^)]*\)).*/\1/')
            
            # Check internal links (relative paths ending in .md or starting with ./)
            if [[ "$url" =~ \.md$ ]] && [[ ! "$url" =~ ^https?:// ]]; then
                check
                local target_file
                
                if [[ "$url" =~ ^\./ ]]; then
                    # Relative to current directory
                    target_file="$(dirname "$file")/$url"
                else
                    # Relative to docs root
                    target_file="${DOCS_DIR}/$url"
                fi
                
                if [ -f "$target_file" ]; then
                    log_success "Valid link in $relative_path: $url"
                else
                    log_error "Broken link in $relative_path: $url -> $(realpath --relative-to="$DOCS_DIR" "$target_file" 2>/dev/null || echo "$target_file")"
                fi
            fi
        done
    done
}

# Check diagram syntax
check_mermaid_diagrams() {
    log_info "Checking Mermaid diagram syntax..."
    
    # Create temporary files for diagram validation
    local temp_dir=$(mktemp -d)
    local diagram_count=0
    
    find "${DOCS_DIR}" -name "*.md" -type f | while read -r file; do
        local relative_path=${file#${DOCS_DIR}/}
        local current_diagram=0
        
        # Extract Mermaid diagrams
        awk '
        /^```mermaid$/ { in_diagram = 1; current_diagram++; next }
        /^```$/ && in_diagram { 
            close(output_file)
            in_diagram = 0
            next 
        }
        in_diagram {
            if (!output_file) {
                output_file = "'"$temp_dir"'/diagram_" current_diagram ".mmd"
            }
            print > output_file
        }
        ' "$file"
        
        # Count diagrams in this file
        local file_diagram_count=$(grep -c '^```mermaid$' "$file" 2>/dev/null || echo 0)
        diagram_count=$((diagram_count + file_diagram_count))
        
        if [ $file_diagram_count -gt 0 ]; then
            log_info "Found $file_diagram_count Mermaid diagram(s) in $relative_path"
        fi
    done
    
    # Basic syntax validation for extracted diagrams
    if [ $diagram_count -gt 0 ]; then
        for diagram_file in "$temp_dir"/*.mmd; do
            if [ -f "$diagram_file" ]; then
                check
                local first_line=$(head -n1 "$diagram_file")
                
                # Check for valid Mermaid diagram types
                if [[ "$first_line" =~ ^(graph|flowchart|sequenceDiagram|classDiagram|stateDiagram|erDiagram|journey|gantt|pie|gitGraph) ]]; then
                    log_success "Valid Mermaid diagram type: $first_line"
                else
                    log_error "Invalid or missing Mermaid diagram type in $(basename "$diagram_file"): $first_line"
                fi
            fi
        done
    else
        log_info "No Mermaid diagrams found"
    fi
    
    # Cleanup
    rm -rf "$temp_dir"
}

# Check table of contents consistency
check_toc_consistency() {
    log_info "Checking table of contents consistency..."
    
    local readme_file="${DOCS_DIR}/README.md"
    
    if [ -f "$readme_file" ]; then
        # Check if all documented files are linked in README
        local doc_files=(
            "01-system-context.md"
            "02-container-architecture.md"
            "03-component-architecture.md"
            "04-data-architecture.md"
            "05-security-architecture.md"
            "06-quality-attributes.md"
        )
        
        for doc_file in "${doc_files[@]}"; do
            check
            if grep -q "$doc_file" "$readme_file"; then
                log_success "README links to $doc_file"
            else
                log_error "README missing link to $doc_file"
            fi
        done
        
        # Check ADRs directory link
        check
        if grep -q "adrs/" "$readme_file"; then
            log_success "README links to ADRs directory"
        else
            log_error "README missing link to ADRs directory"
        fi
    else
        log_error "README.md not found in docs directory"
    fi
}

# Check content quality
check_content_quality() {
    log_info "Checking content quality..."
    
    find "${DOCS_DIR}" -name "*.md" -type f | while read -r file; do
        local relative_path=${file#${DOCS_DIR}/}
        
        # Skip template files
        if [[ "$relative_path" == "adrs/template.md" ]]; then
            continue
        fi
        
        # Check minimum content length
        local word_count=$(wc -w < "$file")
        local line_count=$(wc -l < "$file")
        
        check
        if [ $word_count -gt 100 ]; then
            log_success "$relative_path has substantial content ($word_count words)"
        else
            log_error "$relative_path has insufficient content ($word_count words, minimum 100)"
        fi
        
        # Check for TODO markers
        check
        local todo_count=$(grep -c -i "TODO\|FIXME\|XXX" "$file" 2>/dev/null || echo 0)
        if [ $todo_count -eq 0 ]; then
            log_success "$relative_path has no TODO markers"
        else
            log_warning "$relative_path has $todo_count TODO marker(s)"
        fi
        
        # Check for proper headings structure
        check
        if grep -q "^# " "$file"; then
            log_success "$relative_path has main heading"
        else
            log_error "$relative_path missing main heading (# title)"
        fi
    done
}

# Generate validation report
generate_report() {
    local report_file="${DOCS_DIR}/validation-report.md"
    local timestamp=$(date -u +"%Y-%m-%d %H:%M:%S UTC")
    
    cat > "$report_file" << EOF
# Documentation Validation Report

**Generated:** $timestamp  
**Total Checks:** $TOTAL_CHECKS  
**Passed:** $PASSED_CHECKS  
**Failed:** $FAILED_CHECKS  
**Success Rate:** $(echo "scale=1; $PASSED_CHECKS * 100 / $TOTAL_CHECKS" | bc -l 2>/dev/null || echo "N/A")%

## Summary

EOF
    
    if [ $FAILED_CHECKS -eq 0 ]; then
        echo "✅ **All validation checks passed!**" >> "$report_file"
        echo "" >> "$report_file"
        echo "The architecture documentation is complete and consistent." >> "$report_file"
    else
        echo "❌ **$FAILED_CHECKS validation check(s) failed.**" >> "$report_file"
        echo "" >> "$report_file"
        echo "Please review the issues above and update the documentation accordingly." >> "$report_file"
    fi
    
    echo "" >> "$report_file"
    echo "## Validation Categories" >> "$report_file"
    echo "" >> "$report_file"
    echo "- ✅ Required Files" >> "$report_file"
    echo "- ✅ ADR Structure" >> "$report_file"
    echo "- ✅ Internal Links" >> "$report_file"
    echo "- ✅ Mermaid Diagrams" >> "$report_file"
    echo "- ✅ Table of Contents" >> "$report_file"
    echo "- ✅ Content Quality" >> "$report_file"
    
    log_info "Validation report saved to: $report_file"
}

# Print final summary
print_summary() {
    echo
    echo "========================================="
    echo "Documentation Validation Summary"
    echo "========================================="
    echo "Total Checks: $TOTAL_CHECKS"
    echo "Passed: $PASSED_CHECKS"
    echo "Failed: $FAILED_CHECKS"
    
    if [ $FAILED_CHECKS -eq 0 ]; then
        echo -e "Status: ${GREEN}✅ ALL CHECKS PASSED${NC}"
        exit 0
    else
        echo -e "Status: ${RED}❌ $FAILED_CHECKS CHECKS FAILED${NC}"
        exit 1
    fi
}

# Main execution
main() {
    log_info "Starting documentation validation..."
    echo
    
    check_required_files
    echo
    check_adr_structure
    echo
    check_internal_links
    echo
    check_mermaid_diagrams
    echo
    check_toc_consistency
    echo
    check_content_quality
    echo
    
    generate_report
    print_summary
}

# Run if called directly
if [ "${BASH_SOURCE[0]}" = "${0}" ]; then
    main "$@"
fi