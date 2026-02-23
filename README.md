<div align="center">
  <img src=".github/assets/logo.svg" alt="Ohara Logo" width="200">
  <h1>Ohara</h1>

  <a href="https://github.com/tanishqrupaal/ohara/actions/workflows/release.yaml"><img alt="Build Workflow" src="https://github.com/tanishqrupaal/ohara/actions/workflows/release.yaml/badge.svg"></a>&nbsp;<a href="https://hub.docker.com/r/tanishqrupaal/ohara"><img alt="Docker Pulls" src="https://img.shields.io/docker/pulls/tanishqrupaal/ohara"></a><br>
  <a href="https://github.com/tanishqrupaal/ohara/releases"><img alt="GitHub Release" src="https://img.shields.io/github/v/release/tanishqrupaal/ohara"></a><br><br>

  <a href="#screenshots">Screenshots</a> &bull; <a href="#installation-and-usage">Install & Use</a> &bull; <a href="#tips-and-notes">Tips & Notes</a>
</div>

---

A self-hosted web application for tracking professional achievements, glue work, and feedback as "Touchpoints." Optimized for both manual entry through its web GUI and automated data entry via REST API using AI assistants or `curl`. The name is based on the island of Ohara from One Piece.

## Features

- Track touchpoints with categories, tags, people involved, and URLs
- 12-month timeline chart with category diversity overlay
- Generate and view Markdown reports with syntax highlighting and Mermaid diagrams
- Filter by date range, category, and tags
- Self-contained single binary with embedded frontend — no external CDN dependencies

## Screenshots

<details>
<summary>Click to expand screenshots</summary>

*Screenshots coming soon*

</details>

## Installation and Usage

### Docker (Recommended)

```bash
docker run -d -p 8080:8080 -v ohara-data:/data tanishqrupaal/ohara
```

### Binary

Download from [releases](https://github.com/tanishqrupaal/ohara/releases) and run:

```bash
./ohara --data-dir ./data --port 8080
```

### Build from Source

```bash
git clone https://github.com/tanishqrupaal/ohara
cd ohara
make build
./ohara --data-dir ./data --port 8080
```

## Tips and Notes

- All dates are stored in UTC and displayed in the browser's local timezone
- The `--debug` flag enables verbose zerolog output for troubleshooting
- Data is stored as flat JSON files in the data directory — no database required
- Reports are Markdown files stored in `<data-dir>/reports/` and support code blocks, Mermaid diagrams, and GFM tables
- Categories and tags are validated against `metadata.json` — add new ones via the API before using them
