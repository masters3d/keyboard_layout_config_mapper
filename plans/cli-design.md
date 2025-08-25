# CLI Design Plan: Keyboard Layout Configuration Mapper (KLCM)

## Vision
Create an intuitive, user-friendly CLI that makes keyboard layout synchronization and management effortless. Focus on making the common case (ZMK ↔ ZMK) simple while providing power user features for complex scenarios.

## CLI Architecture

### **Command Structure**
```bash
klcm <command> [subcommand] [options] [arguments]
```

### **Core Commands**

#### 1. **`klcm status`** - Show current state
```bash
klcm status                    # Show all keyboards status
klcm status --keyboard adv360  # Show specific keyboard
klcm status --changes          # Show uncommitted changes
```

#### 2. **`klcm sync`** - Synchronize configurations
```bash
# Primary workflow (from main keyboard to others)
klcm sync                      # Auto-detect changes and sync
klcm sync --from adv360        # Sync from Advantage360 to all others
klcm sync --from adv360 --to glove80  # Specific sync
klcm sync --dry-run           # Preview changes without applying

# Interactive mode for complex changes
klcm sync --interactive       # Guide through conflicts/decisions
klcm sync --force            # Skip validation and force sync
```

#### 3. **`klcm edit`** - Smart editing assistance
```bash
klcm edit                     # Interactive editor selection
klcm edit --keyboard adv360   # Edit specific keyboard
klcm edit --layer 0           # Edit specific layer
klcm edit --behavior homerow  # Edit specific behavior
```

#### 4. **`klcm validate`** - Configuration validation
```bash
klcm validate                 # Validate all configurations
klcm validate --keyboard glove80  # Validate specific keyboard
klcm validate --compile       # Test compilation
```

#### 5. **`klcm diff`** - Show differences
```bash
klcm diff adv360 glove80      # Show differences between keyboards
klcm diff --layer 0           # Show differences in specific layer
klcm diff --semantic          # Show functional differences
```

#### 6. **`klcm backup`** - Backup management
```bash
klcm backup create            # Create backup of current state
klcm backup list              # List available backups
klcm backup restore <name>    # Restore from backup
```

#### 7. **`klcm update`** - Fetch latest configurations
```bash
klcm update                   # Update all from upstream
klcm update --keyboard adv360 # Update specific keyboard
klcm update --check           # Check for updates without applying
```

#### 8. **`klcm pr`** - GitHub integration
```bash
klcm pr create               # Create PRs for all changed keyboards
klcm pr create --keyboard adv360  # Create PR for specific keyboard
klcm pr status              # Show PR status
klcm pr merge               # Auto-merge ready PRs
```

## User Experience Design

### **Progressive Disclosure**
- **Beginner**: Simple commands with sensible defaults
- **Intermediate**: Additional options and interactive modes  
- **Advanced**: Full control with detailed options

### **Smart Defaults**
- Auto-detect which keyboard has changes
- Intelligently choose sync direction
- Suggest best practices

### **Helpful Feedback**
```bash
$ klcm sync
✓ Detected changes in adv360 configuration
✓ Analyzing layer 0 changes...
⚠ Found homerow modifier timing change
➜ Syncing to glove80... ✓
➜ Syncing to qmk_ergodx... ⚠ Manual review needed
➜ Syncing to kinesis2... ⚠ Limited compatibility

Summary:
- ✅ glove80: Synced successfully
- ⚠️  qmk_ergodx: 3 items need manual review
- ⚠️  kinesis2: 8 items couldn't be mapped

Run 'klcm review' to address pending items
```

## Technical Implementation

### **Technology Stack**
- **Language**: Go (fast, single binary, excellent CLI libraries)
- **CLI Framework**: Cobra + Viper (configuration)
- **Parsing**: Custom parsers for each format
- **Git Integration**: go-git library
- **GitHub API**: Native GitHub API client
- **Configuration**: YAML for KLCM settings

### **Project Structure**
```
klcm/
├── cmd/                    # CLI commands
│   ├── root.go
│   ├── sync.go
│   ├── edit.go
│   └── ...
├── internal/
│   ├── parsers/            # Configuration parsers
│   │   ├── zmk.go
│   │   ├── qmk.go
│   │   └── kinesis.go
│   ├── mappers/            # Cross-platform mapping
│   │   ├── layout.go
│   │   ├── behavior.go
│   │   └── layer.go
│   ├── generators/         # Configuration generators
│   │   ├── zmk.go
│   │   ├── qmk.go
│   │   └── kinesis.go
│   ├── git/               # Git operations
│   ├── github/            # GitHub API integration
│   └── config/            # KLCM configuration
├── pkg/                   # Public packages
│   ├── keyboard/          # Keyboard definitions
│   └── layout/            # Layout representations
└── configs/               # Default configurations
```

### **Core Data Structures**

#### **Keyboard Definition**
```go
type Keyboard struct {
    Name         string
    Type         KeyboardType // ZMK, QMK, Kinesis2
    ConfigPath   string
    Repository   GitRepository
    LayoutMatrix [][]string    // Physical layout mapping
    Capabilities []Capability  // What features it supports
}
```

#### **Layout Representation**
```go
type Layout struct {
    Layers     []Layer
    Behaviors  []Behavior
    Macros     []Macro
    Combos     []Combo
    Metadata   LayoutMetadata
}
```

#### **Sync Operation**
```go
type SyncOperation struct {
    Source      Keyboard
    Targets     []Keyboard
    Changes     []Change
    Conflicts   []Conflict
    Strategy    SyncStrategy
}
```

## User Workflow Examples

### **Daily Use Case** (Primary: ZMK ↔ ZMK)
```bash
# User makes changes to Advantage360
$ klcm sync
✓ Auto-detected changes in adv360
✓ Synced to glove80 successfully
✓ Created backup: backup-2024-01-15-143022

# User wants to test changes
$ klcm validate --compile
✓ adv360: Compilation successful
✓ glove80: Compilation successful

# User is happy, wants to push changes
$ klcm pr create
✓ Created PR for adv360: "Update homerow modifier timing"
✓ Created PR for glove80: "Sync from adv360: homerow timing"
```

### **Complex Change Workflow**
```bash
# User makes significant behavior changes
$ klcm sync --interactive
? Changes detected in layer 0. How should we handle the new combo?
  > Skip for incompatible keyboards
    Try to adapt for all keyboards
    Review each keyboard individually

? QMK ErgoDox: Custom ZMK behavior can't be directly translated. Options:
  > Use closest QMK equivalent
    Create custom QMK function
    Skip this change
    Let me handle manually

✓ Sync completed with 2 manual review items
? Would you like to review them now? (Y/n) y

Opening review mode...
```

### **Emergency Recovery**
```bash
$ klcm backup restore
? Multiple backups found. Which would you like to restore?
  > backup-2024-01-15-143022 (2 hours ago) - "Before homerow changes"
    backup-2024-01-15-120000 (5 hours ago) - "Weekly backup"
    backup-2024-01-14-180000 (yesterday) - "Before layer restructure"

✓ Restored configuration successfully
? Would you like to create a PR to revert upstream changes? (y/N)
```

## Configuration File
```yaml
# ~/.klcm/config.yaml
keyboards:
  adv360:
    type: zmk
    path: ./configs/zmk_adv360/adv360.keymap
    repository:
      url: https://github.com/masters3d/Adv360-Pro-ZMK
      branch: cheyo
      path: config/adv360.keymap
    primary: true  # This is the user's primary keyboard
  
  glove80:
    type: zmk
    path: ./configs/zmk_glove80/glove80.keymap
    repository:
      url: https://github.com/masters3d/glove80-zmk-config
      branch: cheyo
      path: config/glove80.keymap

sync:
  auto_backup: true
  default_strategy: "conservative"  # conservative, aggressive, interactive
  primary_keyboard: "adv360"

github:
  auto_pr: false
  pr_title_template: "KLCM: Sync from {{.Source}} - {{.Summary}}"
  review_required: true
```

## Success Metrics
- **Time to sync**: < 10 seconds for common changes
- **User confidence**: Clear feedback on what will happen
- **Error recovery**: Always possible to undo/restore
- **Learning curve**: New users productive within 5 minutes
- **Expert efficiency**: Power users can script complex workflows

## Implementation Phases

### **Phase 1: Foundation** (Week 1-2)
- [ ] Basic CLI structure with Cobra
- [ ] ZMK parser for both keyboards
- [ ] Simple sync logic (adv360 → glove80)
- [ ] Configuration validation

### **Phase 2: Core Features** (Week 3-4)
- [ ] All 4 keyboard parsers
- [ ] Cross-platform mapping logic
- [ ] Interactive sync mode
- [ ] Backup/restore functionality

### **Phase 3: Advanced Features** (Week 5-6)
- [ ] GitHub integration
- [ ] Smart diff and conflict resolution
- [ ] Advanced validation and testing
- [ ] Performance optimization

### **Phase 4: Polish** (Week 7-8)
- [ ] Comprehensive error handling
- [ ] Documentation and help system
- [ ] User testing and feedback
- [ ] Release preparation

## Next Steps
1. Set up Go project structure
2. Implement basic ZMK parser
3. Create simple adv360 → glove80 sync
4. Build CLI framework around it
5. Iterate based on real usage