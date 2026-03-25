# parser

This package contains all specialized optimizers (parsers) that compress terminal output for specific CLIs and tools.

```mermaid
graph TD
    Interface["Parser Interface [pkg/parser/interface.go]"]
    Git["GitParser [pkg/parser/git.go]"]
    Thneed["ThneedParser [pkg/parser/thneed.go]"]
    NPM["NPMParser [pkg/parser/npm.go]"]
    FS["FS/LS Parser [pkg/parser/fs.go]"]
    
    Interface <|-- Git
    Interface <|-- Thneed
    Interface <|-- NPM
    Interface <|-- FS
```
