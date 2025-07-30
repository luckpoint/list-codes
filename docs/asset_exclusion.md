# Asset File Exclusion

**list-codes** automatically excludes non-essential asset files to keep the output focused on source code for LLM analysis.

## Overview

Since **list-codes** is designed specifically for source code analysis, asset files like images, fonts, and media files are automatically excluded by default. This prevents output bloat while maintaining focus on analyzable code.

## Automatically Excluded Asset Types

**list-codes** automatically excludes these non-code file types by default:

### Images
- **Raster Images**: `.png`, `.jpg`, `.jpeg`, `.gif`, `.bmp`, `.ico`, `.tiff`, `.webp`, `.avif`, `.heic`
- **Vector Images**: `.svg` (excluded since they're typically decorative assets, not code)

### Fonts
- **Web Fonts**: `.woff`, `.woff2`, `.ttf`, `.otf`, `.eot`

### Media Files
- **Audio**: `.mp3`, `.wav`, `.ogg`, `.aac`, `.flac`, `.m4a`
- **Video**: `.mp4`, `.avi`, `.mov`, `.webm`, `.mkv`, `.wmv`

### Archives & Binaries
- **Archives**: `.zip`, `.tar`, `.gz`, `.rar`, `.7z`, `.bz2`
- **Documents**: `.pdf`, `.docx`, `.xlsx`, `.pptx`
- **Executables**: `.exe`, `.dmg`, `.deb`, `.rpm`

### Asset Directories
Common asset directories are also excluded:
- `assets/`, `static/`, `public/assets/`, `images/`, `media/`
- `dist/assets/`, `build/assets/`, `uploads/`

## Including Asset Files When Needed

If you need to include specific asset files for analysis (rare cases), use the `--include` flag:

```bash
# Include specific SVG files that contain meaningful code
list-codes --include "src/icons/*.svg"

# Include documentation images
list-codes --include "docs/images/*.png"

# Include a specific asset directory
list-codes --include "public/config-assets/**"
```

The `--include` flag overrides automatic asset exclusion for the specified patterns.

## Why Assets Are Excluded

1. **Size Efficiency**: Asset files are typically large and would bloat LLM input
2. **Relevance**: Binary assets contain no analyzable source code
3. **Focus**: Keeps analysis centered on actual code logic and structure
4. **Performance**: Reduces processing time and memory usage

## Implementation Details

Asset exclusion is implemented at the file system scanning level, similar to how dotfiles and test files are handled. The exclusion happens before file size limits are applied, making the tool more efficient.