package main

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/yuin/goldmark/util"
)

// DiffLineType classifies a line in a diff-annotated code block.
type DiffLineType int

const (
	DiffContext  DiffLineType = iota
	DiffAddition              // line prefixed with '+'
	DiffDeletion              // line prefixed with '-'
)

var (
	filenameIcons = map[string]string{
		"makefile":          "fas fa-cogs",
		"dockerfile":        "fab fa-docker",
		"readme.md":         "fas fa-book",
		"license":           "fas fa-certificate",
		"license.txt":       "fas fa-certificate",
		"license.md":        "fas fa-certificate",
		"go.mod":            "fab fa-golang",
		"go.sum":            "fab fa-golang",
		"package.json":      "fab fa-npm",
		"package-lock.json": "fab fa-npm",
		".gitignore":        "fab fa-git-alt",
		".gitmodules":       "fab fa-git-alt",
		".gitattributes":    "fab fa-git-alt",
	}

	extensionIcons = map[string]string{
		".qml":       "fas fa-file-code",
		".json":      "fas fa-file-code",
		".go":        "fab fa-golang",
		".js":        "fab fa-js",
		".ts":        "fab fa-js-square",
		".html":      "fab fa-html5",
		".htm":       "fab fa-html5",
		".css":       "fab fa-css3-alt",
		".py":        "fab fa-python",
		".rb":        "fas fa-gem",
		".php":       "fab fa-php",
		".java":      "fab fa-java",
		".rs":        "fas fa-gears",
		".c":         "fas fa-code",
		".cpp":       "fas fa-code",
		".h":         "fas fa-code",
		".hpp":       "fas fa-code",
		".cs":        "fas fa-code",
		".swift":     "fas fa-code",
		".kt":        "fas fa-k",
		".kts":       "fas fa-k",
		".yml":       "fas fa-file-code",
		".yaml":      "fas fa-file-code",
		".toml":      "fas fa-cog",
		".ini":       "fas fa-sliders-h",
		".conf":      "fas fa-sliders-h",
		".config":    "fas fa-sliders-h",
		".env":       "fas fa-key",
		".xml":       "fas fa-code",
		".csv":       "fas fa-file-csv",
		".sql":       "fas fa-database",
		".txt":       "fas fa-file-lines",
		".pdf":       "fas fa-file-pdf",
		".doc":       "fas fa-file-word",
		".docx":      "fas fa-file-word",
		".xls":       "fas fa-file-excel",
		".xlsx":      "fas fa-file-excel",
		".ppt":       "fas fa-file-powerpoint",
		".pptx":      "fas fa-file-powerpoint",
		".svg":       "fas fa-bezier-curve",
		".mp3":       "fas fa-file-audio",
		".wav":       "fas fa-file-audio",
		".ogg":       "fas fa-file-audio",
		".flac":      "fas fa-file-audio",
		".mp4":       "fas fa-file-video",
		".mov":       "fas fa-file-video",
		".avi":       "fas fa-file-video",
		".mkv":       "fas fa-file-video",
		".webm":      "fas fa-file-video",
		".vue":       "fab fa-vuejs",
		".svelte":    "fas fa-fire",
		".gradle":    "fab fa-android",
		".xcodeproj": "fab fa-apple",
	}
)

// getFileIcon returns a Font Awesome CSS class for the given filename,
// matching first by exact basename then by file extension.
func getFileIcon(filename string) string {
	basename := strings.ToLower(filepath.Base(filename))
	if icon, exists := filenameIcons[basename]; exists {
		return icon
	}

	ext := strings.ToLower(filepath.Ext(filename))
	if icon, exists := extensionIcons[ext]; exists {
		return icon
	}

	return ""
}

// DirectoryNode represents a file or folder in a parsed directory tree
// used for rendering directory-structure code blocks.
type DirectoryNode struct {
	Name     string
	IsFolder bool
	Indent   int
	Children []*DirectoryNode
}

// buildDirectoryTree parses indented text into a tree of DirectoryNode,
// using indentation levels to determine parent-child relationships.
func buildDirectoryTree(content string) *DirectoryNode {
	var lines []string
	for _, line := range strings.Split(content, "\n") {
		if len(strings.TrimSpace(line)) > 0 {
			lines = append(lines, line)
		}
	}

	root := &DirectoryNode{Children: []*DirectoryNode{}}
	if len(lines) == 0 {
		return root
	}

	type stackEntry struct {
		node   *DirectoryNode
		indent int
	}

	currentParent := root
	stack := []stackEntry{{node: root, indent: -1}}

	for i, line := range lines {
		trimmedLine := strings.TrimLeft(line, " \t")
		currentIndent := len(line) - len(trimmedLine)

		entry := &DirectoryNode{
			Name:     trimmedLine,
			IsFolder: strings.HasSuffix(trimmedLine, "/"),
			Indent:   currentIndent,
			Children: []*DirectoryNode{},
		}

		if i > 0 {
			if currentIndent > stack[len(stack)-1].indent {
				lastChild := currentParent.Children[len(currentParent.Children)-1]
				stack = append(stack, stackEntry{node: currentParent, indent: currentParent.Indent})
				currentParent = lastChild
			} else if currentIndent < stack[len(stack)-1].indent {
				for len(stack) > 1 && stack[len(stack)-1].indent >= currentIndent {
					stack = stack[:len(stack)-1]
					currentParent = stack[len(stack)-1].node
				}
			}
		}

		currentParent.Children = append(currentParent.Children, entry)
	}

	return root
}

// renderDirectoryTreeRecursive writes HTML for a DirectoryNode and its
// children, nesting folders within dir-children divs.
func renderDirectoryTreeRecursive(w util.BufWriter, entry *DirectoryNode) {
	for _, child := range entry.Children {
		if child.IsFolder {
			w.WriteString(fmt.Sprintf("<div class=\"dir-entry dir-folder\"><i class=\"fas fa-folder-open\"></i> %s</div>", child.Name))
		} else {
			fileIcon := getFileIcon(child.Name)
			w.WriteString(fmt.Sprintf("<div class=\"dir-entry dir-file\"><i class=\"%s\"></i> %s</div>", fileIcon, child.Name))
		}

		if len(child.Children) > 0 {
			w.WriteString("<div class=\"dir-children\">")
			renderDirectoryTreeRecursive(w, child)
			w.WriteString("</div>")
		}
	}
}

// renderDirectoryStructure parses indented text into a directory tree
// and writes the complete HTML representation to w.
func renderDirectoryStructure(w util.BufWriter, content string) {
	root := buildDirectoryTree(content)
	w.WriteString("<div class=\"directory-tree\">")
	renderDirectoryTreeRecursive(w, root)
	w.WriteString("</div>")
}

// processDiffContent strips leading +/- markers from diff-formatted code,
// returning the cleaned content and a per-line type classification.
func processDiffContent(content string) (string, []DiffLineType) {
	lines := strings.Split(content, "\n")
	var cleanedLines []string
	var lineTypes []DiffLineType

	for _, line := range lines {
		if len(line) == 0 {
			cleanedLines = append(cleanedLines, line)
			lineTypes = append(lineTypes, DiffContext)
			continue
		}

		if line[0] == '+' {
			cleanedLines = append(cleanedLines, line[1:])
			lineTypes = append(lineTypes, DiffAddition)
		} else if line[0] == '-' {
			cleanedLines = append(cleanedLines, line[1:])
			lineTypes = append(lineTypes, DiffDeletion)
		} else {
			cleanedLines = append(cleanedLines, line)
			lineTypes = append(lineTypes, DiffContext)
		}
	}

	return strings.Join(cleanedLines, "\n"), lineTypes
}

// postProcessDiffHTML wraps Chroma-highlighted HTML lines with diff CSS
// classes (gi for additions, gd for deletions) based on line types.
func postProcessDiffHTML(html string, lineTypes []DiffLineType) string {
	lines := strings.Split(html, "\n")
	var processedLines []string

	for i, line := range lines {
		if i < len(lineTypes) {
			switch lineTypes[i] {
			case DiffAddition:
				if strings.Contains(line, "<span class=\"line\">") {
					processedLines = append(processedLines, strings.Replace(line, "<span class=\"line\">", "<span class=\"line gi\">", 1))
				} else {
					processedLines = append(processedLines, fmt.Sprintf("<span class=\"gi\">%s</span>", line))
				}
			case DiffDeletion:
				if strings.Contains(line, "<span class=\"line\">") {
					processedLines = append(processedLines, strings.Replace(line, "<span class=\"line\">", "<span class=\"line gd\">", 1))
				} else {
					processedLines = append(processedLines, fmt.Sprintf("<span class=\"gd\">%s</span>", line))
				}
			default:
				processedLines = append(processedLines, line)
			}
		} else {
			processedLines = append(processedLines, line)
		}
	}

	return strings.Join(processedLines, "\n")
}
