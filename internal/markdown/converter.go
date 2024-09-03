package markdown

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/sirupsen/logrus"

	"github.com/radurobot/go-markdown-crawler/internal/storage"
)

var converter = md.NewConverter("", true, nil)

func ConvertToMarkdown(site string, body []byte, outputFolder string, hashStore storage.HashStore) error {
	markdownContent, err := converter.ConvertString(string(body))
	if err != nil {
		logrus.Errorf("Failed to convert HTML to Markdown: %s", err)
		return err
	}

	uniqueSections := synthesizeContent(markdownContent, hashStore)

	finalContent := strings.Join(uniqueSections, "\n\n")

	sanitizedFilename := sanitizeFilename(site)
	markdownFile := filepath.Join(outputFolder, sanitizedFilename+".md")

	if err := os.WriteFile(markdownFile, bytes.NewBufferString(finalContent).Bytes(), 0644); err != nil {
		logrus.Errorf("Failed to write Markdown file %s: %s", markdownFile, err)
		return err
	}

	logrus.Infof("Successfully wrote %s", markdownFile)
	return nil
}

func sanitizeFilename(urlStr string) string {
	u, err := url.Parse(urlStr)
	if err != nil {
		return "invalid_url"
	}
	return strings.ReplaceAll(u.Host+u.Path, "/", "_")
}

func synthesizeContent(content string, hashStore storage.HashStore) []string {
	sections := strings.Split(content, "\n\n")
	var uniqueSections []string

	for _, section := range sections {
		sectionHash := fmt.Sprintf("%x", md5.Sum([]byte(section)))

		if !hashStore.Exists(sectionHash) {
			hashStore.Add(sectionHash)
			uniqueSections = append(uniqueSections, section)
		}
	}

	return uniqueSections
}
