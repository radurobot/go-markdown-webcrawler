# Overview

go-markdown-webcrawler is a tool for web crawling and HTML-to-Markdown conversion. It's nothing more than a wrapper. It leverages the crawling capabilities of [gospider](https://github.com/jaeles-project/gospider) and enhances it by converting crawled HTML content to Markdown format using [html-to-markdown](https://github.com/JohannesKaufmann/html-to-markdown). The tool ensures that only unique content is stored by generating and storing hashes for each paragraph.

## Features

- **Web Crawling**: Uses gospider's crawling logic to fetch content from websites.
- **HTML to Markdown Conversion**: Converts fetched HTML content to Markdown format.
- **Unique Content Storage**: Each paragraph is hashed, and only unique hashes are stored to avoid duplication.
- **Storage Options**: Supports in-memory storage or temporary persistent storage using SQLite for hash management when crawling large websites.
- **Configurable via CLI**: A wide range of options are available to customize the crawling and conversion process - thanks to [gospider](https://github.com/jaeles-project/gospider).

## GO Install

```sh
go install github.com/radurobot/go-markdown-crawler@latest
```

## Usage

### Basic Usage

Crawl a single site and convert the content to Markdown:

```sh
go-markdown-webcrawler --site=http://example.com --output=output_dir
```

### Crawling Multiple Sites

You can specify a list of sites to crawl:

```sh
go-markdown-webcrawler --sites=sites.txt --output=output_dir
```

### Custom User-Agent and Proxy

Specify a custom user-agent and proxy:

```sh
go-markdown-webcrawler --site=http://example.com --user-agent="MyCustomAgent" --proxy="http://127.0.0.1:8080" --output=output_dir
```

### Storage Options

By default, hashes are stored in memory. You can persist them using SQLite:

```sh
go-markdown-webcrawler --site=http://example.com --output=output_dir --in-memory=false
```

This will store the hashes in a SQLite database file (`hashes.db` by default).

### Full List of Options

```sh
go-markdown-webcrawler --help
```

## License

go-markdown-webcrawler is licensed under the MIT License. See the [LICENSE](LICENSE) file for more details.

## Credits

- [gospider](https://github.com/jaeles-project/gospider) for the crawling logic.
- [html-to-markdown](https://github.com/JohannesKaufmann/html-to-markdown) for the conversion library.
