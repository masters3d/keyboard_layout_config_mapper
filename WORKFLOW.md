# 🚀 KLCM Keyboard Configuration Workflow

This guide walks you through the complete process of making, validating, and contributing keyboard configuration changes using KLCM.

## 📋 Quick Reference

```bash
# 1. Update to latest configurations
klcm pull --preview                    # Preview remote changes
klcm pull                              # Apply updates

# 2. Make your changes
# Edit configs/zmk_adv360/adv360.keymap or configs/zmk_glove80/glove80.keymap

# 3. Validate your changes
klcm validate adv360                   # Validate specific keyboard
klcm validate                          # Validate all keyboards

# 4. Sync changes between keyboards (optional)
klcm sync adv360 glove80 --preview     # Preview sync
klcm sync adv360 glove80               # Apply sync

# 5. Create pull requests
klcm pr create --dry-run               # Preview PR creation
klcm pr create                         # Create actual PRs
```

## 🔄 Complete Workflow Process

### **Phase 1: Preparation** 🏗️

#### **Step 1.1: Update to Latest Remote Configurations**
```bash
# Always start with the latest configurations
klcm pull --preview
```
**What this does:**
- Downloads latest configurations from remote repositories
- Shows you exactly what has changed since your last pull
- Helps prevent merge conflicts

**Decision Point:** If there are remote changes, apply them:
```bash
klcm pull  # Apply the updates
```

#### **Step 1.2: Validate Current State**
```bash
klcm validate
```
**What this does:**
- Ensures all configurations are syntactically correct
- Identifies any existing issues before you make changes
- Provides baseline validation

**Expected Output:** ✅ All configurations should validate successfully

---

### **Phase 2: Making Changes** ✏️

#### **Step 2.1: Edit Configuration Files**
Edit the keyboard configuration files directly:
- **Advantage360**: `configs/zmk_adv360/adv360.keymap`
- **Glove80**: `configs/zmk_glove80/glove80.keymap`

**Best Practices:**
- Make small, focused changes
- Test one change at a time
- Keep a note of what you're changing and why

#### **Step 2.2: Immediate Validation**
After each change:
```bash
klcm validate adv360    # If you changed Advantage360
klcm validate glove80   # If you changed Glove80
```

**What to look for:**
- ✅ **Success**: Configuration is valid
- ❌ **Error**: Fix the syntax error before continuing

---

### **Phase 3: Cross-Keyboard Synchronization** 🔄

#### **Step 3.1: Preview Sync (if you want same changes on both keyboards)**
```bash
klcm sync adv360 glove80 --preview
```

**What this shows:**
- All differences between the two keyboard configurations
- What changes would be copied from source to target
- Line-by-line diff of modifications

#### **Step 3.2: Apply Sync (if desired)**
```bash
klcm sync adv360 glove80
```

**When to sync:**
- ✅ You want the same layout on both keyboards
- ✅ You made a useful change that applies to both
- ❌ The keyboards have different physical layouts
- ❌ You want keyboard-specific customizations

#### **Step 3.3: Validate After Sync**
```bash
klcm validate
```

---

### **Phase 4: Final Validation & Review** ✅

#### **Step 4.1: Comprehensive Validation**
```bash
klcm validate --verbose
```

#### **Step 4.2: Review Your Changes**
```bash
git status                    # See which files changed
git diff                      # See exactly what changed
```

#### **Step 4.3: Compare with Remote**
```bash
klcm compare-remote adv360    # Check differences from remote
klcm compare-remote glove80   # Check differences from remote
```

---

### **Phase 5: Contributing Back** 🎁

#### **Step 5.1: Preview PR Creation**
```bash
klcm pr create --dry-run
```

**What this shows:**
- Which repositories will get PRs
- What branch names will be used
- Summary of changes for each PR

#### **Step 5.2: Create Pull Requests**
```bash
klcm pr create
```

**What this does:**
- Creates feature branches in upstream repositories
- Commits your changes with descriptive messages
- Opens pull requests with detailed descriptions
- Provides PR URLs for tracking

#### **Step 5.3: Monitor PR Status**
```bash
klcm pr status               # Check status of your PRs
```

---

## 🛡️ Error Recovery & Troubleshooting

### **If Validation Fails**
```bash
# See detailed error information
klcm validate --verbose adv360

# Common fixes:
# 1. Check for missing semicolons
# 2. Verify key binding syntax: &kp KEY_NAME
# 3. Check for unclosed braces { }
# 4. Ensure proper indentation
```

### **If Sync Creates Unwanted Changes**
```bash
# Preview first to see what would change
klcm sync adv360 glove80 --preview

# If you don't want all changes, make them manually instead
```

### **If Pull Conflicts with Local Changes**
```bash
# See what conflicts exist
klcm pull --preview

# Options:
# 1. Backup your changes: cp configs/zmk_adv360/adv360.keymap /tmp/my-backup.keymap
# 2. Pull updates: klcm pull
# 3. Reapply your changes manually
```

---

## 🎯 Common Workflows

### **Workflow A: Single Keyboard Change**
```bash
klcm pull                           # Update to latest
# Edit configs/zmk_adv360/adv360.keymap
klcm validate adv360                # Validate
klcm pr create                      # Create PR
```

### **Workflow B: Sync Change to Both Keyboards**
```bash
klcm pull                           # Update to latest
# Edit configs/zmk_adv360/adv360.keymap
klcm validate adv360                # Validate first keyboard
klcm sync adv360 glove80 --preview  # Preview sync
klcm sync adv360 glove80            # Apply sync
klcm validate                       # Validate all
klcm pr create                      # Create PRs for both
```

### **Workflow C: Keyboard-Specific Customizations**
```bash
klcm pull                           # Update to latest
# Edit configs/zmk_adv360/adv360.keymap
# Edit configs/zmk_glove80/glove80.keymap (different changes)
klcm validate                       # Validate all
klcm pr create                      # Create separate PRs
```

---

## 🔧 Advanced Tips

### **Branch Management**
```bash
# Use custom branch names
klcm pr create --branch "feature/my-awesome-layout"

# Check PR status
klcm pr status
```

### **Validation Best Practices**
- Always validate after each change
- Use `--verbose` flag for detailed error information
- Test configurations on actual hardware when possible

### **Sync Best Practices**
- Always preview sync before applying
- Consider keyboard-specific differences (key counts, layouts)
- Sync from your "primary" keyboard to others

---

## 📞 Getting Help

### **CLI Help**
```bash
klcm --help                    # General help
klcm pull --help              # Command-specific help
klcm validate --help          # Validation options
klcm pr --help                # PR management help
```

### **Verbose Output**
Add `--verbose` to any command for detailed information:
```bash
klcm validate --verbose
klcm pull --verbose --preview
```

This workflow ensures your changes are always valid, properly tested, and ready for contribution back to the community! 🎉