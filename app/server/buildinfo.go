package main

import (
	"runtime/debug"
	"time"
)

type BuildInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// ReadBuildInfo возвращает [BuildInfo], заполненную информацей из исполняемого файла.
func ReadBuildInfo() BuildInfo {
	debugbi, _ := debug.ReadBuildInfo()
	return parseBuildInfo(debugbi)
}

// parseBuildInfo возвращает [BuildInfo].
//
// Дополнительно поле [BuildInfo.Version] обогащается хэшем коммита,
// т.к. до go1.24 для main пакета здесь будет "(devel)", см. [golang/go#50603].
//
// [golang/go#50603]: https://github.com/golang/go/issues/50603#issuecomment-2181188811
func parseBuildInfo(debugbi *debug.BuildInfo) BuildInfo {
	var bi BuildInfo
	bi.Name, bi.Version = debugbi.Main.Path, debugbi.Main.Version

	// Начиная с go1.24 это ненужно,
	// см. https://github.com/golang/go/issues/50603#issuecomment-2181188811
	if bi.Version != "(devel)" {
		return bi
	}

	var vtime, vhash string
	var modified bool

	const pvlen = 12 // Длина хеша в псевдо-версии

	for _, setting := range debugbi.Settings {
		switch setting.Key {
		case "vcs.revision":
			vhash = setting.Value
			if len(vhash) > pvlen {
				vhash = vhash[:pvlen]
			}
			vhash = "-" + vhash
		case "vcs.time":
			t, _ := time.Parse(time.RFC3339, setting.Value)
			vtime = "-" + t.Format(`20060102150405`)
		case "vcs.modified":
			modified = setting.Value == "true"
		}
	}

	if vtime == "" || vhash == "" {
		return bi
	}

	// https://go.dev/ref/mod#pseudo-versions
	bi.Version = "v0.0.0" + vtime + vhash
	if modified {
		bi.Version += "+dirty"
	}

	return bi
}
