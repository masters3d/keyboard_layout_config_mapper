package cli

import (
	"fmt"
	"strings"
)

// DiffOptions controls how the diff is displayed
type DiffOptions struct {
	ContextLines int    // Number of context lines around changes (default 3)
	Color        bool   // Whether to use color output
	Unified      bool   // Use unified diff format (git-style)
	ShowHeader   bool   // Show file headers
	MaxWidth     int    // Maximum line width before truncation
}

// DefaultDiffOptions returns sensible defaults for diff display
func DefaultDiffOptions() DiffOptions {
	return DiffOptions{
		ContextLines: 3,
		Color:        true,
		Unified:      true,
		ShowHeader:   true,
		MaxWidth:     120,
	}
}

// DiffChunk represents a section of changes in a diff
type DiffChunk struct {
	LocalStart  int
	LocalCount  int
	RemoteStart int
	RemoteCount int
	Lines       []DiffLine
}

// DiffLine represents a single line in a diff
type DiffLine struct {
	Type    LineType
	Content string
	LocalNo int  // Line number in local file (0 if not applicable)
	RemoteNo int // Line number in remote file (0 if not applicable)
}

// LineType represents the type of a diff line
type LineType int

const (
	LineContext LineType = iota // unchanged line (context)
	LineAdded                   // line added in remote (green, +)
	LineRemoved                 // line removed from local (red, -)
)

// UnifiedDiff generates a git-style unified diff between two texts
func UnifiedDiff(localPath, remotePath string, localContent, remoteContent string, opts DiffOptions) string {
	localLines := strings.Split(strings.TrimSpace(localContent), "\n")
	remoteLines := strings.Split(strings.TrimSpace(remoteContent), "\n")
	
	// Handle empty files
	if len(localLines) == 1 && localLines[0] == "" {
		localLines = []string{}
	}
	if len(remoteLines) == 1 && remoteLines[0] == "" {
		remoteLines = []string{}
	}
	
	// If files are identical, return empty diff
	if equalLines(localLines, remoteLines) {
		return ""
	}
	
	chunks := generateDiffChunks(localLines, remoteLines, opts.ContextLines)
	
	var result strings.Builder
	
	// File headers
	if opts.ShowHeader {
		if opts.Color {
			result.WriteString(colorize("--- ", colorRed))
			result.WriteString(colorize(localPath, colorBold))
			result.WriteString(colorize(" (local)\n", colorRed))
			result.WriteString(colorize("+++ ", colorGreen))
			result.WriteString(colorize(remotePath, colorBold))
			result.WriteString(colorize(" (remote)\n", colorGreen))
		} else {
			result.WriteString(fmt.Sprintf("--- %s (local)\n", localPath))
			result.WriteString(fmt.Sprintf("+++ %s (remote)\n", remotePath))
		}
	}
	
	// Generate chunks
	for _, chunk := range chunks {
		// Chunk header
		header := fmt.Sprintf("@@ -%d,%d +%d,%d @@",
			chunk.LocalStart, chunk.LocalCount,
			chunk.RemoteStart, chunk.RemoteCount)
		
		if opts.Color {
			result.WriteString(colorize(header, colorCyan))
		} else {
			result.WriteString(header)
		}
		result.WriteString("\n")
		
		// Chunk lines
		for _, line := range chunk.Lines {
			content := line.Content
			if len(content) > opts.MaxWidth {
				content = content[:opts.MaxWidth-3] + "..."
			}
			
			switch line.Type {
			case LineContext:
				if opts.Color {
					result.WriteString(colorize(" "+content, colorReset))
				} else {
					result.WriteString(" " + content)
				}
			case LineRemoved:
				if opts.Color {
					result.WriteString(colorize("-"+content, colorRed))
				} else {
					result.WriteString("-" + content)
				}
			case LineAdded:
				if opts.Color {
					result.WriteString(colorize("+"+content, colorGreen))
				} else {
					result.WriteString("+" + content)
				}
			}
			result.WriteString("\n")
		}
	}
	
	return result.String()
}

// SimpleDiffSummary provides a concise summary of differences
func SimpleDiffSummary(localLines, remoteLines []string) string {
	if equalLines(localLines, remoteLines) {
		return "No differences"
	}
	
	stats := getDiffStats(localLines, remoteLines)
	
	var parts []string
	if stats.Added > 0 {
		parts = append(parts, fmt.Sprintf("+%d", stats.Added))
	}
	if stats.Removed > 0 {
		parts = append(parts, fmt.Sprintf("-%d", stats.Removed))
	}
	if stats.Modified > 0 {
		parts = append(parts, fmt.Sprintf("~%d", stats.Modified))
	}
	
	summary := strings.Join(parts, " ")
	if summary == "" {
		return "Minor differences"
	}
	
	return fmt.Sprintf("%s lines", summary)
}

// DiffStats holds statistics about differences
type DiffStats struct {
	Added    int
	Removed  int
	Modified int
	Total    int
}

// Color constants for terminal output
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
	colorWhite  = "\033[37m"
	colorBold   = "\033[1m"
)

// Helper functions

func equalLines(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func colorize(text, color string) string {
	return color + text + colorReset
}

func getDiffStats(localLines, remoteLines []string) DiffStats {
	// Simple algorithm to estimate diff stats
	maxLen := len(localLines)
	if len(remoteLines) > maxLen {
		maxLen = len(remoteLines)
	}
	
	var added, removed, modified int
	
	for i := 0; i < maxLen; i++ {
		var localLine, remoteLine string
		
		hasLocal := i < len(localLines)
		hasRemote := i < len(remoteLines)
		
		if hasLocal {
			localLine = strings.TrimSpace(localLines[i])
		}
		if hasRemote {
			remoteLine = strings.TrimSpace(remoteLines[i])
		}
		
		if !hasLocal && hasRemote {
			added++
		} else if hasLocal && !hasRemote {
			removed++
		} else if localLine != remoteLine {
			modified++
		}
	}
	
	return DiffStats{
		Added:    added,
		Removed:  removed,
		Modified: modified,
		Total:    added + removed + modified,
	}
}

func generateDiffChunks(localLines, remoteLines []string, contextLines int) []DiffChunk {
	// This is a simplified diff algorithm. For production use, consider using
	// a proper LCS (Longest Common Subsequence) algorithm like Myers' algorithm.
	
	var chunks []DiffChunk
	maxLines := len(localLines)
	if len(remoteLines) > maxLines {
		maxLines = len(remoteLines)
	}
	
	if maxLines == 0 {
		return chunks
	}
	
	// Find all differences
	var diffs []int
	for i := 0; i < maxLines; i++ {
		var localLine, remoteLine string
		
		if i < len(localLines) {
			localLine = strings.TrimSpace(localLines[i])
		}
		if i < len(remoteLines) {
			remoteLine = strings.TrimSpace(remoteLines[i])
		}
		
		if localLine != remoteLine {
			diffs = append(diffs, i)
		}
	}
	
	if len(diffs) == 0 {
		return chunks
	}
	
	// Group nearby differences into chunks
	chunkStart := 0
	for i := 1; i <= len(diffs); i++ {
		isLast := i == len(diffs)
		isDistant := !isLast && diffs[i]-diffs[i-1] > contextLines*2
		
		if isLast || isDistant {
			// Create chunk from chunkStart to i-1
			chunk := createChunk(localLines, remoteLines, diffs[chunkStart], diffs[i-1], contextLines)
			chunks = append(chunks, chunk)
			chunkStart = i
		}
	}
	
	return chunks
}

func createChunk(localLines, remoteLines []string, startDiff, endDiff, contextLines int) DiffChunk {
	// Calculate chunk boundaries with context
	start := startDiff - contextLines
	if start < 0 {
		start = 0
	}
	
	end := endDiff + contextLines
	maxLen := len(localLines)
	if len(remoteLines) > maxLen {
		maxLen = len(remoteLines)
	}
	if end >= maxLen {
		end = maxLen - 1
	}
	
	var lines []DiffLine
	localIdx := start
	remoteIdx := start
	
	for i := start; i <= end; i++ {
		var localLine, remoteLine string
		hasLocal := i < len(localLines)
		hasRemote := i < len(remoteLines)
		
		if hasLocal {
			localLine = localLines[i]
		}
		if hasRemote {
			remoteLine = remoteLines[i]
		}
		
		if hasLocal && hasRemote && strings.TrimSpace(localLine) == strings.TrimSpace(remoteLine) {
			// Context line
			lines = append(lines, DiffLine{
				Type:     LineContext,
				Content:  localLine,
				LocalNo:  localIdx + 1,
				RemoteNo: remoteIdx + 1,
			})
			localIdx++
			remoteIdx++
		} else {
			// Difference
			if hasLocal && (!hasRemote || strings.TrimSpace(localLine) != strings.TrimSpace(remoteLine)) {
				lines = append(lines, DiffLine{
					Type:    LineRemoved,
					Content: localLine,
					LocalNo: localIdx + 1,
				})
				localIdx++
			}
			if hasRemote && (!hasLocal || strings.TrimSpace(localLine) != strings.TrimSpace(remoteLine)) {
				lines = append(lines, DiffLine{
					Type:     LineAdded,
					Content:  remoteLine,
					RemoteNo: remoteIdx + 1,
				})
				remoteIdx++
			}
		}
	}
	
	return DiffChunk{
		LocalStart:  start + 1,
		LocalCount:  localIdx - start,
		RemoteStart: start + 1,
		RemoteCount: remoteIdx - start,
		Lines:       lines,
	}
}