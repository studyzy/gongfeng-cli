package cmd

// Version 是当前程序版本，构建时通过 -ldflags 注入
var Version = "dev"

func init() {
	rootCmd.Version = Version
	rootCmd.SetVersionTemplate("gongfeng version {{.Version}}\n")
}
