/*
Package version provides full version information for an application.
All the version information are populated at build time.

Examples of populating version info at build time:

	version=$( git describe --tags --dirty --abbrev=14 | sed -E 's/-([0-9]+)-g/.\1+/' )
	revision=$( git rev-parse --short HEAD 2> /dev/null || echo 'unknown' )
	branch=$( git rev-parse --abbrev-ref HEAD 2> /dev/null || echo 'unknown' )

	BUILD_USER=${BUILD_USER:-"${USER}@${HOSTNAME}"}
	BUILD_DATE=${BUILD_DATE:-$( date +%Y%m%d-%H:%M:%S )}

	go build -ldflags '-X github.com/funkygao/golib/version.Version=${version} -X github.com/funkygao/golib/version.Revision=${revision}'
*/
package version

// Build information. Populated at build-time.
var (
	Version   = "unknown"
	Revision  = "unknown"
	Branch    = "unknown"
	BuildUser = "unknown"
	BuildDate = "unknown"
        GoVersion = "unknown"
)

// Info provides the iterable version information.
var Info = map[string]string{
	"version":   Version,
	"revision":  Revision,
	"branch":    Branch,
	"buildUser": BuildUser,
	"buildDate": BuildDate,
        "goVersion": GoVersion,
}
