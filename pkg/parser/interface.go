package parser

type Parser interface {
	Name() string
	CanParse(cmd string, args []string) bool
	Parse(output string) string
}

func GetAllParsers() []Parser {
	return []Parser{
		// Git
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
		// Infra
		&EnvParser{},
		&DockerPsParser{},
		&DependencyParser{},
		&TestParser{},
	}
}
