package presets

func privilegedPresetsMap() map[string]bool {
	return map[string]bool{
		"base": false,
	}
}

func PrivilegedPreset(name string) bool { return privilegedPresetsMap()[name] }
func ListPrivilegedPresets() []string   { return sortedKeys(privilegedPresetsMap()) }

func initPresetsMap() map[string]bool {
	return map[string]bool{
		"base": true,
	}
}

func InitPreset(name string) bool { return initPresetsMap()[name] }
func ListInitPresets() []string   { return sortedKeys(initPresetsMap()) }

func capAddPresetsMap() map[string][]string {
	return map[string][]string{
		"base": {"SYS_PTRACE"},
		"net":  {"NET_ADMIN"},
	}
}

func CapAddPreset(name string) []string { return capAddPresetsMap()[name] }
func ListCapAddPresets() []string       { return sortedKeys(capAddPresetsMap()) }

func capDropPresetsMap() map[string][]string {
	return map[string][]string{
		"base": {"ALL"},
	}
}

func CapDropPreset(name string) []string { return capDropPresetsMap()[name] }
func ListCapDropPresets() []string       { return sortedKeys(capDropPresetsMap()) }

func securityOptPresetsMap() map[string][]string {
	return map[string][]string{
		"base":       {"seccomp=unconfined"},
		"unconfined": {"seccomp=unconfined", "apparmor=unconfined"},
	}
}

func SecurityOptPreset(name string) []string { return securityOptPresetsMap()[name] }
func ListSecurityOptPresets() []string       { return sortedKeys(securityOptPresetsMap()) }

func devicesPresetsMap() map[string][]string {
	return map[string][]string{
		"base": {"/dev/net/tun"},
		"fuse": {"/dev/fuse:/dev/fuse"},
	}
}

func DevicesPreset(name string) []string { return devicesPresetsMap()[name] }
func ListDevicesPresets() []string       { return sortedKeys(devicesPresetsMap()) }
