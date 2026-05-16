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
	// filter state
	filter    string
	filtering bool
	fCursor   int // cursor within filtered results
	fOffset   int // scroll offset within filtered results
}

// IsFiltering reports whether the list is in filter mode.
func (lm ListModel) IsFiltering() bool { return lm.filtering }

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

// SetHeight updates the visible row count and reclamps the scroll offset.
func (lm *ListModel) SetHeight(h int) {
	lm.height = h
	lm.clampScroll()
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

// filteredItems returns all non-separator items that match the current filter.
func (lm ListModel) filteredItems() []ListItem {
	f := strings.ToLower(lm.filter)
	var out []ListItem
	for _, it := range lm.items {
		if it.Separator {
			continue
		}
		if f == "" || strings.Contains(strings.ToLower(it.Key), f) {
			out = append(out, it)
		}
	}
	return out
}

// SelectedItem returns the currently highlighted item (nil if separator or empty).
func (lm ListModel) SelectedItem() *ListItem {
	if lm.filtering {
		items := lm.filteredItems()
		if lm.fCursor >= len(items) {
			return nil
		}
		it := items[lm.fCursor]
		return &it
	}
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
	key, ok := msg.(tea.KeyMsg)
	if !ok {
		return lm, nil
	}
	if lm.filtering {
		return lm.updateFilter(key)
	}
	switch key.String() {
	case "/":
		lm.filtering = true
		lm.filter = ""
		lm.fCursor = 0
		lm.fOffset = 0
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
			k := it.Key
			return lm, func() tea.Msg { return DeleteItemMsg{Key: k} }
		}
	}
	return lm, nil
}

func (lm ListModel) updateFilter(key tea.KeyMsg) (ListModel, tea.Cmd) {
	switch key.String() {
	case "esc":
		lm.filtering = false
		lm.filter = ""
		lm.fCursor = 0
		lm.fOffset = 0
	case "enter":
		items := lm.filteredItems()
		if lm.fCursor < len(items) {
			sel := items[lm.fCursor].Key
			for i, it := range lm.items {
				if it.Key == sel {
					lm.cursor = i
					lm.clampScroll()
					break
				}
			}
		}
		lm.filtering = false
	case "backspace", "ctrl+h":
		if len(lm.filter) > 0 {
			lm.filter = lm.filter[:len(lm.filter)-1]
			lm.fCursor = 0
			lm.fOffset = 0
		}
	case "up", "k":
		lm.moveFCursor(-1)
	case "down", "j":
		lm.moveFCursor(1)
	default:
		if r := key.Runes; len(r) == 1 && r[0] >= 32 {
			lm.filter += string(r)
			lm.fCursor = 0
			lm.fOffset = 0
		}
	}
	return lm, nil
}

func (lm *ListModel) moveFCursor(delta int) {
	items := lm.filteredItems()
	n := len(items)
	if n == 0 {
		return
	}
	lm.fCursor = (lm.fCursor + delta + n) % n
	lm.clampFScroll()
}

func (lm *ListModel) clampFScroll() {
	visH := lm.height - 1
	if visH <= 0 {
		return
	}
	if lm.fCursor < lm.fOffset {
		lm.fOffset = lm.fCursor
	}
	if lm.fCursor >= lm.fOffset+visH {
		lm.fOffset = lm.fCursor - visH + 1
	}
}

// moveCursor advances by delta, wrapping at the list edges and skipping
// separator rows. The +n term keeps the modulus positive for negative deltas.
// The outer bound (i < n) prevents an infinite loop if every row is a separator.
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

// renderListItem renders one non-separator row, with or without selection highlight.
func renderListItem(it ListItem, selected bool) string {
	if selected {
		mark := "+"
		if it.Existing {
			mark = "●"
		}
		return selectedItemStyle.Render("▶ " + mark + "  " + it.Key)
	}
	if it.Existing {
		return existingItemStyle.Render("  ●  " + it.Key)
	}
	return availableItemStyle.Render("  +  " + it.Key)
}

// View renders the visible slice of the list.
func (lm ListModel) View() string {
	if lm.filtering {
		return lm.viewFilter()
	}
	var sb strings.Builder
	end := lm.offset + lm.height
	if end > len(lm.items) {
		end = len(lm.items)
	}
	for i := lm.offset; i < end; i++ {
		if i > lm.offset {
			sb.WriteByte('\n')
		}
		it := lm.items[i]
		if it.Separator {
			sb.WriteString(sectionLabelStyle.Render(it.Key))
		} else {
			sb.WriteString(renderListItem(it, i == lm.cursor))
		}
	}
	return sb.String()
}

// viewFilter renders the list in filter mode: filtered items + prompt on the last line.
func (lm ListModel) viewFilter() string {
	items := lm.filteredItems()
	visH := lm.height - 1 // reserve last row for the filter prompt
	end := lm.fOffset + visH
	if end > len(items) {
		end = len(items)
	}

	lines := make([]string, 0, lm.height)
	for i := lm.fOffset; i < end; i++ {
		lines = append(lines, renderListItem(items[i], i == lm.fCursor))
	}
	for len(lines) < visH {
		lines = append(lines, "")
	}
	lines = append(lines, filterPromptStyle.Render("/"+lm.filter+"▋"))
	return strings.Join(lines, "\n")
}
