package cmd

import (
	"a-dns/skill"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func NewSkillCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "skill",
		Short: "Gestion des skills pour agents IA",
	}

	cmd.AddCommand(NewSkillInstallCmd())
	return cmd
}

func NewSkillInstallCmd() *cobra.Command {
	var installAll bool

	cmd := &cobra.Command{
		Use:   "install [agent]",
		Short: "Installe la skill a-dns pour un agent IA",
		Long: `Installe la skill a-dns pour un agent IA spécifique.

Agents supportés: kilo, claude, cursor, copilot, aider, cline

Exemples:
  a-dns skill install kilo      # Installe pour Kilo
  a-dns skill install claude    # Installe pour Claude Code
  a-dns skill install --all     # Installe pour tous les agents`,
		Args: cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if installAll {
				installAllAgents()
				return
			}

			if len(args) == 0 {
				fmt.Fprintln(os.Stderr, "Erreur: spécifiez un agent ou utilisez --all")
				fmt.Fprintln(os.Stderr, "Agents: kilo, claude, cursor, copilot, aider, cline")
				os.Exit(1)
			}

			agentName := args[0]
			agent := skill.GetAgent(agentName)
			if agent == nil {
				fmt.Fprintf(os.Stderr, "Agent inconnu: %s\n", agentName)
				fmt.Fprintln(os.Stderr, "Agents: kilo, claude, cursor, copilot, aider, cline")
				os.Exit(1)
			}

			if err := installAgent(agent); err != nil {
				fmt.Fprintf(os.Stderr, "Erreur: %v\n", err)
				os.Exit(1)
			}
		},
	}

	cmd.Flags().BoolVarP(&installAll, "all", "a", false, "installer pour tous les agents")
	return cmd
}

func installAgent(agent *skill.Agent) error {
	targetPath := expandPath(agent.Path)
	targetFile := filepath.Join(targetPath, agent.Filename)

	if err := os.MkdirAll(targetPath, 0755); err != nil {
		return fmt.Errorf("création répertoire: %w", err)
	}

	if err := os.WriteFile(targetFile, []byte(skill.SkillContent), 0644); err != nil {
		return fmt.Errorf("écriture fichier: %w", err)
	}

	_ = outputResult(map[string]interface{}{
		"agent":  agent.Name,
		"path":   targetFile,
		"status": "installed",
	})
	return nil
}

func installAllAgents() {
	results := make([]map[string]interface{}, 0)
	installed := 0
	failed := 0

	for i := range skill.KnownAgents {
		agent := &skill.KnownAgents[i]
		targetPath := expandPath(agent.Path)

		if err := installAgentQuiet(agent); err != nil {
			failed++
			results = append(results, map[string]interface{}{
				"agent":  agent.Name,
				"status": "error",
				"error":  err.Error(),
			})
		} else {
			installed++
			results = append(results, map[string]interface{}{
				"agent":  agent.Name,
				"path":   filepath.Join(targetPath, agent.Filename),
				"status": "installed",
			})
		}
	}

	_ = outputResult(map[string]interface{}{
		"installed": installed,
		"failed":    failed,
		"agents":    results,
	})
}

func installAgentQuiet(agent *skill.Agent) error {
	targetPath := expandPath(agent.Path)
	targetFile := filepath.Join(targetPath, agent.Filename)

	if err := os.MkdirAll(targetPath, 0755); err != nil {
		return err
	}

	return os.WriteFile(targetFile, []byte(skill.SkillContent), 0644)
}

func expandPath(path string) string {
	if strings.HasPrefix(path, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			return path
		}
		return filepath.Join(home, path[1:])
	}
	return path
}
