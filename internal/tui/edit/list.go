package edit

import (
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// allKnownKeys lists every top-level key recognised by the DevContainer model,
// in a sensible display order.
var allKnownKeys = []string{
	"name", "image", "build", "dockerComposeFile", "service", "runServices",
	"workspaceFolder", "workspaceMount",
	"remoteUser", "containerUser", "updateRemoteUserUID", "userEnvProbe",
	"containerEnv", "remoteEnv", "localEnv",
	"forwardPorts", "appPort", "portsAttributes", "otherPortsAttributes", "mounts",
	"runArgs", "startupCommand", "overrideCommand", "command", "entrypoint",
	"init", "privileged", "capAdd", "capDrop", "securityOpt", "devices",
	"hostRequirements",
	"features", "overrideFeatureInstallOrder",
	"initializeCommand", "onCreateCommand", "updateContentCommand",
	"postCreateCommand", "postStartCommand", "postAttachCommand", "waitFor",
	"watch", "customizations", "secrets", "shutdownAction",
}

// ListItem represents one row in the left panel.
type ListItem struct {
	Key       string
	Existing  bool
	Separator bool // visual divider row, not selectable
}

// SpaceOnItemMsg is sent to the root model when the user presses Space on a list item.
type SpaceOnItemMsg struct {
	Item ListItem
}

// DeleteItemMsg is sent when the user presses d on an existing item.
type DeleteItemMsg struct{ Key string }

// ListModel is the left-panel list.
type ListModel struct {
	items  []ListItem
	cursor int
	height int // visible rows (excluding borders)
	offset int // scroll offset
}

// BuildListItems constructs the merged item list from the currently existing blocks.
// Only keys present in allKnownKeys are shown; unknown keys are silently ignored.
func BuildListItems(existing []Block) []ListItem {
	knownSet := make(map[string]bool, len(allKnownKeys))
	for _, k := range allKnownKeys {
		knownSet[k] = true
	}

	existingSet := make(map[string]bool, len(existing))
	for _, b := range existing {
		if knownSet[b.Key] {
			existingSet[b.Key] = true
		}
	}

	items := make([]ListItem, 0, len(allKnownKeys)+4)

	// Existing known keys in file order.
	existingItems := make([]ListItem, 0)
	for _, b := range existing {
		if knownSet[b.Key] {
			existingItems = append(existingItems, ListItem{Key: b.Key, Existing: true})
		}
	}
	if len(existingItems) > 0 {
		items = append(items, ListItem{Separator: true, Key: "ADDED"})
		items = append(items, existingItems...)
	}

	// Available keys alphabetically (skip already-existing).
	available := make([]string, 0)
	for _, k := range allKnownKeys {
		if !existingSet[k] {
			available = append(available, k)
		}
	}
	sort.Strings(available)

	if len(available) > 0 {
		items = append(items, ListItem{Separator: true, Key: ""}) // blank line
		items = append(items, ListItem{Separator: true, Key: "AVAILABLE"})
		for _, k := range available {
			items = append(items, ListItem{Key: k, Existing: false})
		}
	}

	return items
}

// NewListModel creates the list model.
func NewListModel(existing []Block, height int) ListModel {
	items := BuildListItems(existing)
	cursor := 0
	for i, it := range items {
		if !it.Separator {
			cursor = i
			break
		}
	}
	return ListModel{items: items, cursor: cursor, height: height}
}

// Rebuild refreshes the list after blocks change without losing cursor position.
func (lm *ListModel) Rebuild(existing []Block) {
	prevKey := ""
	if lm.cursor < len(lm.items) && !lm.items[lm.cursor].Separator {
		prevKey = lm.items[lm.cursor].Key
	}
	lm.items = BuildListItems(existing)
	if prevKey != "" {
		for i, it := range lm.items {
			if it.Key == prevKey {
				lm.cursor = i
				lm.clampScroll()
				return
			}
		}
	}
	// Cursor was on a separator, key not found, or list is empty — find first real item.
	lm.cursor = 0
	for i, it := range lm.items {
		if !it.Separator {
			lm.cursor = i
			break
		}
	}
	lm.clampScroll()
}

// AddedCount returns how many recognised top-level keys are present.
func (lm ListModel) AddedCount() int {
	n := 0
	for _, it := range lm.items {
		if it.Existing {
			n++
		}
	}
	return n
}

// SelectedItem returns the currently highlighted item (nil if separator).
func (lm ListModel) SelectedItem() *ListItem {
	if lm.cursor >= len(lm.items) {
		return nil
	}
	it := lm.items[lm.cursor]
	if it.Separator {
		return nil
	}
	return &it
}

func (lm ListModel) Update(msg tea.Msg) (ListModel, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		switch msg.String() {
		case "up", "k":
			lm.moveCursor(-1)
		case "down", "j":
			lm.moveCursor(1)
		case " ":
			if it := lm.SelectedItem(); it != nil {
				item := *it
				return lm, func() tea.Msg { return SpaceOnItemMsg{Item: item} }
			}
		case "d":
			if it := lm.SelectedItem(); it != nil && it.Existing {
				key := it.Key
				return lm, func() tea.Msg { return DeleteItemMsg{Key: key} }
			}
		}
	}
	return lm, nil
}

func (lm *ListModel) moveCursor(delta int) {
	n := len(lm.items)
	if n == 0 {
		return
	}
	for i := 0; i < n; i++ {
		lm.cursor = (lm.cursor + delta + n) % n
		if !lm.items[lm.cursor].Separator {
			break
		}
	}
	lm.clampScroll()
}

func (lm *ListModel) clampScroll() {
	if lm.height <= 0 {
		return
	}
	if lm.cursor < lm.offset {
		lm.offset = lm.cursor
	}
	if lm.cursor >= lm.offset+lm.height {
		lm.offset = lm.cursor - lm.height + 1
	}
}

// View renders the visible slice of the list.
func (lm ListModel) View() string {
	var sb strings.Builder
	end := lm.offset + lm.height
	if end > len(lm.items) {
		end = len(lm.items)
	}
	for i := lm.offset; i < end; i++ {
		it := lm.items[i]
		switch {
		case it.Separator:
			sb.WriteString(sectionLabelStyle.Render(it.Key))
		case i == lm.cursor:
			mark := "+"
			if it.Existing {
				mark = "●"
			}
			sb.WriteString(selectedItemStyle.Render("▶ " + mark + "  " + it.Key))
		case it.Existing:
			sb.WriteString(existingItemStyle.Render("  ●  " + it.Key))
		default:
			sb.WriteString(availableItemStyle.Render("  +  " + it.Key))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}
