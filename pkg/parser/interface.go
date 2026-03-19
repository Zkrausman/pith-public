package parser

type Parser interface {
	Name() string
	CanParse(cmd string, args []string) bool
	Parse(output string) string
}

func GetAllParsers() []Parser {
	parsers := []Parser{
		// Git
		&CompositeGitParser{},
		&GitStatusParser{},
		&GitLogParser{},
		&GitDiffParser{},
		&GitBranchParser{},
		// FS
		&LsParser{},
		&FindParser{},
		&TreeParser{},
		&DuParser{},
		// Text
		&GrepParser{},
		&MinifyParser{},
		&SourceParser{},
		// Infra
		&EnvParser{},
		&DockerPsParser{},
		&GitHubParser{},
		&GitHubReleaseParser{},
		&DependencyParser{},
		&TestParser{},
		// New Tools
		&WebParser{},
		&PithParser{},
		&PowerShellParser{},
		&GoParser{},
	}
	return append(parsers, &ChainParser{})
}
