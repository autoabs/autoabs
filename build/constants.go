package build

const (
	BuildingRank  = 100
	PendingRank   = 200
	FailedRank    = 300
	CompletedRank = 400
	ArchivedRank  = 500
)

const testPkgBuild = `pkgname=(%s)
pkgver=%s
pkgrel=%s
pkgdesc='AutoABS Test'
arch=('%s')
url='https://github.com/autoabs/autoabs'
license=('custom')

%s`

const testPkgBuildPackage = `package_%s() {
  echo 'AutoABS test %s package'
}
`
