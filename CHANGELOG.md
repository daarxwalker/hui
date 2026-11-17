# Changelog

Write the changes of the next release under `## [Unreleased]`. `task release VERSION=vX.Y.Z`
puts them into the git tag and the GitHub release, then empties the section again.

## [Unreleased]
- Fixed: `-dir` pointing at the parent of a generated directory moved it there silently.
  Taking over a directory from v0.1.2 now looks only at the files directly in it, and a
  subdirectory with its own `.hui-assets` counts as foreign.
- `-check` tells you to move or delete a foreign file instead of running `go generate`,
  which would refuse to write. `CheckAssets` returns `[]hui.AssetChange` (`Name`, `Kind`,
  `Module`, `Version`, `String()`) instead of strings.
- An interrupted run no longer leaves its files foreign: `.hui-assets` lists everything the
  run is about to write or delete before it touches a file, and the next run finishes it.
