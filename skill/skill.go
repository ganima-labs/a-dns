package skill

import _ "embed"

//go:embed SKILL.md
var SkillContent string

type Agent struct {
	Name        string
	Description string
	Path        string
	Filename    string
}

var KnownAgents = []Agent{
	{
		Name:        "kilo",
		Description: "Kilo (kilocode)",
		Path:        "~/.kilocode/skills/a-dns",
		Filename:    "SKILL.md",
	},
	{
		Name:        "claude",
		Description: "Claude Code",
		Path:        "~/.claude/skills/a-dns",
		Filename:    "SKILL.md",
	},
	{
		Name:        "cursor",
		Description: "Cursor IDE",
		Path:        "~/.cursor/skills/a-dns",
		Filename:    "SKILL.md",
	},
	{
		Name:        "copilot",
		Description: "GitHub Copilot",
		Path:        "~/.copilot/skills/a-dns",
		Filename:    "SKILL.md",
	},
	{
		Name:        "aider",
		Description: "Aider",
		Path:        "~/.aider/skills/a-dns",
		Filename:    "SKILL.md",
	},
	{
		Name:        "cline",
		Description: "Cline (VS Code)",
		Path:        "~/.cline/skills/a-dns",
		Filename:    "SKILL.md",
	},
}

func GetAgent(name string) *Agent {
	for _, a := range KnownAgents {
		if a.Name == name {
			return &a
		}
	}
	return nil
}
