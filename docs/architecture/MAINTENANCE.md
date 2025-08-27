# Architecture Documentation Maintenance

## Overview

This document outlines the procedures and best practices for maintaining the architecture documentation of the Camel-Do project. It ensures documentation stays current, accurate, and useful as the system evolves.

## Documentation Structure

```
docs/architecture/
├── README.md                    # Main documentation index
├── 01-system-context.md         # System boundaries and external integrations
├── 02-container-architecture.md # Service architecture and deployment
├── 03-component-architecture.md # Internal module structure
├── 04-data-architecture.md      # Data models and persistence
├── 05-security-architecture.md  # Security patterns and threat model
├── 06-quality-attributes.md     # Non-functional requirements
├── adrs/                        # Architecture Decision Records
│   ├── template.md              # ADR template
│   ├── 0001-*.md               # Individual ADRs (numbered sequentially)
│   └── ...
├── automation/                  # Documentation tooling
│   ├── docs-generate.sh         # Generate diagrams and metrics
│   ├── validate-docs.sh         # Validate documentation quality
│   ├── mermaid-config.json      # Diagram styling configuration
│   └── ...
├── generated/                   # Auto-generated content (gitignored)
│   ├── diagrams/               # SVG/PNG exports of Mermaid diagrams
│   ├── metrics/                # Codebase metrics and statistics
│   └── reports/                # Validation and summary reports
└── MAINTENANCE.md              # This file
```

## Maintenance Responsibilities

### When to Update Documentation

#### Code Changes
- **New Features**: Update relevant architecture documents when adding new services, components, or major features
- **Refactoring**: Update component and data architecture when restructuring code
- **Security Changes**: Update security architecture for auth, encryption, or security-related changes
- **External Integrations**: Update system context and container architecture for new external dependencies

#### Architecture Decisions
- **Create ADRs**: Document all significant architectural decisions using the ADR template
- **Update Existing Docs**: Reflect ADR decisions in relevant architecture documents
- **Cross-Reference**: Ensure ADRs are referenced in appropriate architecture documents

#### Quality Attribute Changes
- **Performance Requirements**: Update quality attributes document for performance targets
- **Security Requirements**: Document new security requirements or compliance needs
- **Scalability Changes**: Update for changes in user base or system scale

### Regular Maintenance Tasks

#### Monthly Reviews
- [ ] Run documentation validation (`./docs/architecture/automation/validate-docs.sh`)
- [ ] Check for broken links and missing references
- [ ] Update generated diagrams and metrics
- [ ] Review and update any TODO markers in documentation

#### Quarterly Updates
- [ ] Review and update system context for external system changes
- [ ] Validate architecture documentation against current codebase
- [ ] Update quality attribute measurements and targets
- [ ] Review ADR status (mark superseded ADRs if applicable)

#### Annual Reviews
- [ ] Comprehensive architecture documentation review
- [ ] Update documentation structure if needed
- [ ] Archive or reorganize outdated documentation
- [ ] Review and update maintenance procedures

## Automation Tools

### Documentation Generation
```bash
# Generate all documentation artifacts
./docs/architecture/automation/docs-generate.sh

# This script:
# - Extracts and converts Mermaid diagrams to SVG/PNG
# - Generates architecture metrics from codebase
# - Creates summary reports
# - Validates documentation completeness
```

### Documentation Validation
```bash
# Validate documentation quality and consistency
./docs/architecture/automation/validate-docs.sh

# This script checks:
# - Required files exist
# - ADR structure and numbering
# - Internal link validity
# - Mermaid diagram syntax
# - Table of contents consistency
# - Content quality metrics
```

### Continuous Integration
Add to CI/CD pipeline:
```yaml
# Example GitHub Actions workflow
- name: Validate Architecture Documentation
  run: |
    cd docs/architecture/automation
    ./validate-docs.sh
    
- name: Generate Documentation Artifacts
  run: |
    cd docs/architecture/automation
    ./docs-generate.sh
```

## ADR Management

### Creating New ADRs

1. **Use Sequential Numbering**: Next available number (check existing ADRs)
2. **Copy Template**: Use `adrs/template.md` as starting point
3. **Follow Format**: Include all required sections (Status, Context, Decision, Consequences)
4. **Get Review**: Have ADRs reviewed before marking as "Accepted"

```bash
# Create new ADR
cp docs/architecture/adrs/template.md docs/architecture/adrs/0005-new-decision.md

# Edit the new ADR file
# Update the README.md to link to the new ADR
```

### ADR Lifecycle
- **Proposed**: Initial draft, under discussion
- **Accepted**: Decision approved and being implemented
- **Rejected**: Decision rejected, keep for historical context
- **Deprecated**: No longer applicable, replaced by newer ADR
- **Superseded**: Replaced by specific newer ADR (reference the new ADR)

### Updating Existing ADRs
- **Never modify accepted ADRs**: Create new ADR that supersedes the old one
- **Update status only**: Change status from "Proposed" to "Accepted/Rejected"
- **Add references**: Link to related ADRs or implementation details

## Diagram Management

### Mermaid Diagrams
- **Embed in Markdown**: Use ```mermaid code blocks in documentation
- **Consistent Styling**: Use the provided mermaid-config.json for consistent appearance
- **Automatic Generation**: Run automation scripts to generate SVG/PNG exports
- **Version Control**: Keep diagram source in markdown, generated images are gitignored

### Diagram Types and Usage
- **System Context**: High-level system interactions (flowchart/graph)
- **Container Architecture**: Service relationships and data flow (flowchart)
- **Component Architecture**: Internal component relationships (flowchart/class)
- **Sequence Diagrams**: Process flows and interactions (sequenceDiagram)
- **Data Models**: Entity relationships (erDiagram)

## Content Guidelines

### Writing Style
- **Clear and Concise**: Use simple, direct language
- **Present Tense**: Describe current state, not future plans
- **Active Voice**: Prefer active over passive voice
- **Technical Accuracy**: Ensure technical details are correct and current

### Structure Guidelines
- **Logical Flow**: Organize content from high-level to detailed
- **Cross-References**: Link related concepts across documents
- **Examples**: Include code examples and concrete implementations
- **Visual Aids**: Use diagrams to illustrate complex concepts

### Review Checklist
Before committing documentation updates:
- [ ] Run validation script and fix any issues
- [ ] Check spelling and grammar
- [ ] Verify all links work correctly
- [ ] Ensure diagrams render correctly
- [ ] Update table of contents if needed
- [ ] Add appropriate cross-references

## Tool Dependencies

### Required Tools
- **Node.js**: For Mermaid CLI diagram generation
- **Mermaid CLI**: `npm install -g @mermaid-js/mermaid-cli`
- **Go**: For codebase analysis and metrics
- **Bash**: For automation scripts (Linux/macOS/WSL)
- **jq**: For JSON processing in scripts
- **bc**: For calculations in scripts

### Installation
```bash
# Install Mermaid CLI
npm install -g @mermaid-js/mermaid-cli

# Verify installation
mmdc --version

# Install other dependencies (Ubuntu/Debian)
sudo apt-get install jq bc

# Install other dependencies (macOS)
brew install jq bc
```

## Troubleshooting

### Common Issues

#### Validation Failures
```bash
# Run validation with detailed output
./docs/architecture/automation/validate-docs.sh

# Common fixes:
# - Fix broken internal links
# - Add missing ADR sections
# - Update table of contents
# - Fix Mermaid diagram syntax
```

#### Diagram Generation Issues
```bash
# Test Mermaid diagram generation manually
mmdc -i diagram.mmd -o diagram.svg

# Common fixes:
# - Check mermaid-config.json syntax
# - Verify Mermaid CLI installation
# - Fix diagram syntax errors
```

#### Missing Dependencies
```bash
# Check which tools are missing
which node mmdc go jq bc

# Install missing tools as needed
# Refer to installation section above
```

## Best Practices

### Documentation Updates
1. **Update Documentation First**: Update docs before or with code changes
2. **Small, Frequent Updates**: Prefer small, regular updates over large batch updates
3. **Review Changes**: Have documentation changes reviewed like code changes
4. **Test Automation**: Ensure automation scripts work after documentation changes

### ADR Best Practices
1. **Be Specific**: Include concrete details and rationale
2. **Consider Alternatives**: Document why alternatives were not chosen
3. **Update Status**: Keep ADR status current (Proposed → Accepted/Rejected)
4. **Reference Implementation**: Link to code that implements the decision

### Diagram Best Practices
1. **Keep Simple**: Avoid overly complex diagrams
2. **Consistent Style**: Use the same styling across all diagrams
3. **Update Together**: Update diagrams when architecture changes
4. **Version Control**: Keep diagram source in version control

## Contact and Support

For questions about architecture documentation maintenance:
- Create an issue in the project repository
- Discuss in team architecture reviews
- Reference this maintenance guide for procedures

Regular maintenance of architecture documentation ensures it remains a valuable resource for understanding and evolving the Camel-Do system.