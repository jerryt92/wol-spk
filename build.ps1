[CmdletBinding()]
param()

$ErrorActionPreference = 'Stop'

$RootDir = $PSScriptRoot
Push-Location $RootDir
try {
    & go run ./tools/build
    if ($LASTEXITCODE -ne 0) {
        throw "Build failed with exit code $LASTEXITCODE"
    }
}
finally {
    Pop-Location
}
