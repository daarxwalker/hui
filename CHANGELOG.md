# Changelog

Write the changes of the next release under `## [Unreleased]`. `task release VERSION=vX.Y.Z`
puts them into the git tag and the GitHub release, then empties the section again.

## [Unreleased]

- `go tool hui assets [-dir DIR] [-check]` writes the Alpine plugins as `hui/hui.js` into the
  app's frontend sources, at exactly the hui version in `go.mod`. Run it through
  `//go:generate`, commit the file and import it from the app's bundler entry.
- `hui.Assets()`, `hui.AssetHeader()`, `hui.WriteAssets()` and `hui.CheckAssets()` let a
  library on top of hui write its own assets together with hui's.
- The npm package `@daar/hui-alpine` is no longer published. Replace
  `import hui from '@daar/hui-alpine'` with the generated `./hui/hui.js`.
