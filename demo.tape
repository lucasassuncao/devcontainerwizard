Output demo.gif

Set FontSize 20
Set Width 1920
Set Height 1080
Set Padding 20
Set Margin 20
Set MarginFill "#1e1e2e"
Set BorderRadius 8
Set WindowBar Colorful
Set Theme "Catppuccin Mocha"
Set PlaybackSpeed 1

Env TERM "dumb"

Type "# How to use devcontainerwizard..."
Enter
Sleep 3s

Type "# Let's start by creating a config.yaml from a template"
Enter
Sleep 3s

Type "devcontainerwizard init -t dockerfile"
Enter
Sleep 3s

Type "# Showing config.yaml content"
Enter
Sleep 3s

Type "head -n 20 config.yaml"
Enter
Sleep 5s

Type "clear"
Enter
Sleep 3s

Type "# Converting config.yaml to devcontainer.json"
Enter
Sleep 3s

Type "devcontainerwizard"
Enter
Sleep 3s

Type "# Showing devcontainer.json content"
Enter
Sleep 3s

Type "head -n 20 .devcontainer/devcontainer.json"
Enter
Sleep 5s

Type "clear"
Enter
Sleep 3s

Type "# That's it!"
Enter
Sleep 3s