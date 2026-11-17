# Changelog

Write the changes of the next release under `## [Unreleased]`. `task release VERSION=vX.Y.Z`
puts them into the git tag and the GitHub release, then empties the section again.

## [Unreleased]

- Fixed: a wrong `-dir` (for example the folder with the app's sources) made
  `go tool hui assets` delete every file that was not in the new set. The command now owns
  only the files listed in `.hui-assets` (path, module, version per line) and stops with an
  error naming any other file, without writing or deleting anything. A directory without
  the list is written only when it is empty; a directory from v0.1.2 that holds just the
  generated `hui.js` is taken over and gets the list. Commit `.hui-assets` with the rest.
- Fixed: hidden files such as `.DS_Store` or `.gitkeep` made `-check` fail and were deleted.
  Files and directories starting with a dot are now ignored.
- `go tool hui assets` also writes `hui.d.ts`, so `import hui from './hui/hui.js'` compiles in
  TypeScript with `strict` without `@types/alpinejs`.
- `-check` reports each file as `missing`, `differs` (with the version it came from),
  `stale` or `foreign`, plus a missing or outdated `.hui-assets`.
- Breaking for libraries on top of hui: `Assets()` returns `map[string]hui.Asset` (`Module`,
  `Version`, `Content`) and `WriteAssets`/`CheckAssets` take that map. `hui.CreateAsset(module,
  content)` builds an asset of another module, binary files included.
