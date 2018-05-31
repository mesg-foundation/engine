#!/usr/bin/env node

const fs = require("fs")
const exec = require('child_process').execSync

const version = process.argv[2]

const changelog = fs.readFileSync("./CHANGELOG.md").toString()
const date = (new Date()).toJSON().slice(0, 10)


exec(`git tag ${version}`)

const template = `## [Unreleased]

#### Changed
#### Added
#### Removed
#### Fixed

## [${version}] - ${date}`

fs.writeFileSync("./CHANGELOG.md", changelog.replace("## [Unreleased]", template))

exec(`git add CHANGELOG.md`)
exec(`git commit -m "Update version ${version}"`)
exec(`git push && git push origin refs/tags/${version}`)