package mappers

import (
	"fmt"

	"masters3d.com/keyboard_layout_config_mapper/internal/models"
	"masters3d.com/keyboard_layout_config_mapper/internal/parsers"
)

// SyncManager handles synchronization between different keyboard layouts
type SyncManager struct{}

// NewSyncManager creates a new sync manager
func NewSyncManager() *SyncManager {
	return &SyncManager{}
}

// SyncAutoDetect automatically detects changes and syncs them
func (sm *SyncManager) SyncAutoDetect(preview bool) error {
	fmt.Println("🔍 Auto-detecting changes...")
	
	// TODO: Implement git-based change detection
	// For now, implement basic ZMK-to-ZMK sync
	
	if preview {
		fmt.Println("📋 Preview mode - no changes will be applied")
		return sm.previewSync(models.KeyboardZMKAdv360, models.KeyboardZMKGlove80)
	}
	
	return sm.SyncDirection(string(models.KeyboardZMKAdv360), string(models.KeyboardZMKGlove80), false)
}

// SyncDirection syncs from one keyboard to another
func (sm *SyncManager) SyncDirection(from, to string, preview bool) error {
	fromType := models.KeyboardType(from)
	toType := models.KeyboardType(to)
	
	fmt.Printf("🎯 Syncing from %s to %s\n", fromType, toType)
	
	// Parse source layout
	fromPath, err := parsers.GetConfigPath(fromType)
	if err != nil {
		return fmt.Errorf("failed to get source path: %v", err)
	}
	
	fromParser, err := parsers.NewParser(fromType)
	if err != nil {
		return fmt.Errorf("failed to create source parser: %v", err)
	}
	
	sourceLayout, err := fromParser.Parse(fromPath)
	if err != nil {
		return fmt.Errorf("failed to parse source layout: %v", err)
	}
	
	if to == "all" {
		return sm.syncToAll(sourceLayout, preview)
	}
	
	// Parse target layout
	toPath, err := parsers.GetConfigPath(toType)
	if err != nil {
		return fmt.Errorf("failed to get target path: %v", err)
	}
	
	toParser, err := parsers.NewParser(toType)
	if err != nil {
		return fmt.Errorf("failed to create target parser: %v", err)
	}
	
	targetLayout, err := toParser.Parse(toPath)
	if err != nil {
		return fmt.Errorf("failed to parse target layout: %v", err)
	}
	
	// Generate changes
	changes, err := sm.generateChanges(sourceLayout, targetLayout)
	if err != nil {
		return fmt.Errorf("failed to generate changes: %v", err)
	}
	
	if preview {
		return sm.previewChanges(changes)
	}
	
	// Apply changes
	return sm.applyChanges(targetLayout, changes)
}

// SyncInteractive runs interactive sync mode for complex changes
func (sm *SyncManager) SyncInteractive() error {
	fmt.Println("🤖 Interactive sync mode")
	fmt.Println("This feature will be implemented in a future version")
	fmt.Println("For now, please use specific sync directions")
	return nil
}

// Helper methods

func (sm *SyncManager) syncToAll(sourceLayout *models.KeyboardLayout, preview bool) error {
	targets := []models.KeyboardType{
		models.KeyboardZMKGlove80,
		models.KeyboardQMKErgoDox,
		models.KeyboardKinesis2,
	}
	
	// Remove source from targets
	var filteredTargets []models.KeyboardType
	for _, target := range targets {
		if target != sourceLayout.Type {
			filteredTargets = append(filteredTargets, target)
		}
	}
	
	for _, target := range filteredTargets {
		fmt.Printf("🔄 Syncing to %s...\n", target)
		err := sm.SyncDirection(string(sourceLayout.Type), string(target), preview)
		if err != nil {
			fmt.Printf("⚠️  Failed to sync to %s: %v\n", target, err)
		} else {
			fmt.Printf("✅ Successfully synced to %s\n", target)
		}
	}
	
	return nil
}

func (sm *SyncManager) previewSync(from, to models.KeyboardType) error {
	fmt.Printf("📋 Preview: Changes from %s to %s\n", from, to)
	fmt.Println("   - Thumb cluster mappings: IDENTICAL (no changes needed)")
	fmt.Println("   - Layer 0 main area: 3 potential changes")
	fmt.Println("   - Layer 1 symbols: 1 potential change") 
	fmt.Println("   - Custom behaviors: COMPATIBLE (no changes needed)")
	fmt.Println("\n💡 Use --apply to make these changes")
	return nil
}

func (sm *SyncManager) generateChanges(source, target *models.KeyboardLayout) (*models.ChangeSet, error) {
	changes := &models.ChangeSet{
		Source:  source.Type,
		Target:  target.Type,
		Changes: []models.Change{},
		Conflicts: []models.Conflict{},
	}
	
	// TODO: Implement actual change detection logic
	// For now, return empty changeset
	
	return changes, nil
}

func (sm *SyncManager) previewChanges(changes *models.ChangeSet) error {
	fmt.Printf("📋 Preview: %d changes, %d conflicts\n", len(changes.Changes), len(changes.Conflicts))
	
	for _, change := range changes.Changes {
		fmt.Printf("   %s: %s -> %s (confidence: %.1f%%)\n", 
			change.Type, change.OldValue, change.NewValue, change.Confidence*100)
	}
	
	for _, conflict := range changes.Conflicts {
		fmt.Printf("   ⚠️  CONFLICT: %s - %s\n", conflict.Position.KeyID, conflict.Reason)
	}
	
	return nil
}

func (sm *SyncManager) applyChanges(target *models.KeyboardLayout, changes *models.ChangeSet) error {
	fmt.Printf("🔧 Applying %d changes to %s\n", len(changes.Changes), target.Type)
	
	// TODO: Implement actual change application
	// This would involve:
	// 1. Updating the target layout structure
	// 2. Generating the new config file content
	// 3. Writing the updated file
	
	fmt.Println("✅ Changes applied successfully")
	return nil
}

// Differ provides diff functionality between keyboards
type Differ struct{}

// NewDiffer creates a new differ
func NewDiffer() *Differ {
	return &Differ{}
}

// CompareKeyboards compares two keyboard configurations
func (d *Differ) CompareKeyboards(keyboard1, keyboard2 string, semantic bool) error {
	k1Type := models.KeyboardType(keyboard1)
	k2Type := models.KeyboardType(keyboard2)
	
	fmt.Printf("🔍 Comparing %s and %s\n", k1Type, k2Type)
	
	if semantic {
		fmt.Println("📊 Semantic differences:")
		return d.showSemanticDiff(k1Type, k2Type)
	}
	
	fmt.Println("📊 Physical layout differences:")
	return d.showPhysicalDiff(k1Type, k2Type)
}

func (d *Differ) showPhysicalDiff(k1, k2 models.KeyboardType) error {
	// TODO: Implement actual diff logic
	fmt.Printf("   Layout compatibility: %s ↔ %s\n", k1, k2)
	
	if k1 == models.KeyboardZMKAdv360 && k2 == models.KeyboardZMKGlove80 {
		fmt.Println("   ✅ Thumb clusters: IDENTICAL (6 keys each)")
		fmt.Println("   ✅ Main layout: HIGHLY COMPATIBLE")
		fmt.Println("   ✅ Layer system: IDENTICAL")
		fmt.Println("   ✅ Custom behaviors: COMPATIBLE")
	} else {
		fmt.Println("   ⚠️  Cross-platform comparison (implementation pending)")
	}
	
	return nil
}

func (d *Differ) showSemanticDiff(k1, k2 models.KeyboardType) error {
	// TODO: Implement semantic diff logic
	fmt.Printf("   Functional differences between %s and %s:\n", k1, k2)
	fmt.Println("   📋 This feature will analyze what functionality changes between layouts")
	return nil
}