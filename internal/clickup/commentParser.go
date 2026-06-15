package clickup

/*
commentParser.go contains the (partly AI written) function BlocksToHTML to turn quill.js operation blocks into a usable HTML string that is usefull when sending an email to a customer
There might be some magic going on here, but it is mostly clear.
*/
import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

func (c *CodeBlock) UnmarshalJSON(data []byte) error {
	// Try string first
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		c.Lang = s
		return nil
	}

	// Try object: { "code-block": "plain" }
	var obj map[string]string
	if err := json.Unmarshal(data, &obj); err == nil {
		c.Lang = obj["code-block"]
		return nil
	}
	c.Lang = "plain"
	return nil

}

// ─── Renderer ────────────────────────────────────────────────────────────────

func BlocksToHTML(blocks []CommentBlock) string {
	var sb strings.Builder
	i := 0

	for i < len(blocks) {
		block := blocks[i]
		// --- Emoji ---
		if block.Type == "emoticon" && block.Emoticon != nil {
			cp, err := strconv.ParseInt(block.Emoticon.Code, 16, 32)
			if err == nil {
				sb.WriteRune(rune(cp))
			}
			i++
			continue
		}
		// --- @mention ---
		if block.Type == "tag" && block.User != nil {
			name := block.User.Username
			if name == "" {
				name = block.User.Email
			}
			if name == "" {
				name = fmt.Sprintf("%d", block.User.ID)
			}
			fmt.Fprintf(&sb, `<span class="cu-mention" data-user-id="%d">@%s</span>`,
				block.User.ID, escHTML(name))
			i++
			continue
		}
		// --- link_mention ---
		if block.Type == "link_mention" && block.LinkMention != nil {
			url := escHTML(block.LinkMention.URL)
			fmt.Fprintf(&sb, `<a href="%s" target="_blank" rel="noopener">%s</a>`, url, url)
			i++
			continue
		}
		// --- Image ---
		if block.Type == "image" && block.Image != nil {
			url := escHTML(block.Image.URL)
			name := escHTML(block.Image.Name)
			width := block.Attributes.Width
			if width == "" {
				width = "300"
			}
			fmt.Fprintf(&sb, `<img src="%s" alt="%s" width="%s">`, url, name, width)
			i++
			continue
		}
		// --- Bookmark ---
		if block.Type == "bookmark" && block.Bookmark != nil {
			url := escHTML(block.Bookmark.Url)
			if isImageURL(block.Bookmark.Url) {
				fmt.Fprintf(&sb, `<img src="%s" alt="image" width="300">`, url)
			} else {
				fmt.Fprintf(&sb, `<a class="cu-bookmark" href="%s" target="_blank" rel="noopener">🔗 %s</a>`, url, url)
			}
			i++
			continue
		}
		// --- Attachment ---
		if block.Type == "attachment" && block.Attachment != nil {
			url := escHTML(block.Attachment.URL)
			name := block.Attachment.Name
			if name == "" {
				name = url // fall back to URL if no display name
			}
			fmt.Fprintf(&sb, `<a class="cu-attachment" href="%s" target="_blank" rel="noopener">📎 %s</a>`,
				url, escHTML(name))
			i++
			continue
		}
		// --- Code block ---

		// Detect: next \n op has a code-block attribute

		if isCodeBlockOp(blocks, i+1) {
			lang := "plain"
			if blocks[i+1].Attributes.CodeBlock != nil {
				lang = blocks[i+1].Attributes.CodeBlock.Lang
			}
			if lang == "" {
				lang = "plain"
			}
			var code strings.Builder
			for i < len(blocks) {
				op := blocks[i]
				if op.Text == "\n" && op.Attributes.CodeBlock != nil {
					i++ // consume the closing line op
					break
				}
				code.WriteString(escHTML(op.Text))
				i++
			}
			fmt.Fprintf(&sb, "<pre><code class=\"language-%s\">%s</code></pre>",
				escHTML(lang), code.String())
			continue
		}

		// --- List item ---

		if isListOp(blocks, i+1) {
			lineOp := blocks[i+1]
			listType := lineOp.Attributes.List.List
			indent := lineOp.Attributes.Indent
			indentStyle := ""
			if indent > 0 {
				indentStyle = fmt.Sprintf(` style="margin-left:%.1fem"`, float64(indent)*1.5)
			}

			// Collect inline content up to (and consuming) the \n

			var inline strings.Builder
			for i < len(blocks) {
				op := blocks[i]
				if op.Text == "\n" {
					i++ // consume line op
					break
				}
				inline.WriteString(renderInline(op))
				i++
			}

			switch listType {
			case "bullet":
				fmt.Fprintf(&sb, "<ul><li%s>%s</li></ul>", indentStyle, inline.String())
			case "ordered":
				fmt.Fprintf(&sb, "<ol><li%s>%s</li></ol>", indentStyle, inline.String())
			case "checked":
				fmt.Fprintf(&sb, `<ul class="cu-checklist"><li%s><input type="checkbox" checked disabled> %s</li></ul>`,
					indentStyle, inline.String())
			case "unchecked":
				fmt.Fprintf(&sb, `<ul class="cu-checklist"><li%s><input type="checkbox" disabled> %s</li></ul>`,
					indentStyle, inline.String())
			case "toggled":
				fmt.Fprintf(&sb, "<details><summary>%s</summary>\n", inline.String())
				// Collect indented children
				for i < len(blocks) {
					op := blocks[i]
					if op.Text == "\n" {
						if op.Attributes.Indent == 0 {
							break
						}
						i++
						continue
					}
					if op.Attributes.Indent == 0 && op.Type == "" {
						break
					}
					fmt.Fprintf(&sb, "<p>%s</p>", renderInline(op))
					i++
				}
				sb.WriteString("</details>")
			case "none":
				fmt.Fprintf(&sb, "<p%s>%s</p>", indentStyle, inline.String())
			default:
				fmt.Fprintf(&sb, "<li%s>%s</li>", indentStyle, inline.String())
			}
			continue
		}

		// --- Paragraph / inline run ---
		var inline strings.Builder
		for i < len(blocks) {
			op := blocks[i]
			// Special types break out to be handled at top of loop
			if op.Type == "emoticon" || op.Type == "tag" || op.Type == "image" {
				break
			}
			if op.Text == "\n" {
				attrs := op.Attributes
				// Plain \n with no block attrs = line break
				if !attrs.Bold && !attrs.Italic && !attrs.Code &&
					attrs.List == nil && attrs.CodeBlock == nil {
					//inline.WriteString("<br>") TODO add this back if it breaks
				}
				i++
				break
			}
			inline.WriteString(renderInline(op))
			i++
		}
		if inline.Len() > 0 {
			fmt.Fprintf(&sb, "<p>%s</p>", inline.String())
		}
	}
	return sb.String()

}

// ─── Inline renderer ─────────────────────────────────────────────────────────

func renderInline(op CommentBlock) string {
	if op.Text == "" {
		return ""
	}
	out := escHTML(op.Text)
	a := op.Attributes
	if a.Code {
		out = "<code>" + out + "</code>"
	}
	if a.Bold {
		out = "<strong>" + out + "</strong>"
	}
	if a.Italic {
		out = "<em>" + out + "</em>"
	}
	if a.Strike {
		out = "<s>" + out + "</s>"
	}
	if a.Underline {
		out = "<u>" + out + "</u>"
	}
	if a.Link != "" {
		out = fmt.Sprintf(`<a href="%s" target="_blank" rel="noopener">%s</a>`, escHTML(a.Link), out)
	}
	if a.Color != "" {
		out = fmt.Sprintf(`<span style="color:%s">%s</span>`, escHTML(a.Color), out)
	}
	if a.Background != "" {
		out = fmt.Sprintf(`<span style="background-color:%s">%s</span>`, escHTML(a.Background), out)
	}
	return out
}

// ─── Helpers ─────────────────────────────────────────────────────────────────

func isCodeBlockOp(blocks []CommentBlock, idx int) bool {
	if idx >= len(blocks) {
		return false
	}
	return blocks[idx].Text == "\n" && blocks[idx].Attributes.CodeBlock != nil
}

func isListOp(blocks []CommentBlock, idx int) bool {
	if idx >= len(blocks) {
		return false
	}
	return blocks[idx].Text == "\n" && blocks[idx].Attributes.List != nil
}

func escHTML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, `"`, "&quot;")
	s = strings.ReplaceAll(s, "'", "&#39;") // add this
	return s
}
func isImageURL(url string) bool {
	url = strings.ToLower(strings.Split(url, "?")[0]) // strip query params
	return strings.HasSuffix(url, ".png") ||
		strings.HasSuffix(url, ".jpg") ||
		strings.HasSuffix(url, ".jpeg") ||
		strings.HasSuffix(url, ".gif") ||
		strings.HasSuffix(url, ".webp")
}
