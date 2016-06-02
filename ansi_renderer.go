package main

import (
	"bytes"
	"strings"

	md "gopkg.in/russross/blackfriday.v2"
)

type ansiRenderer struct {
	buf        bytes.Buffer
	inLinkNode bool
}

func (r *ansiRenderer) Render(ast *md.Node) []byte {
	ast.Walk(func(node *md.Node, entering bool) md.WalkStatus {
		switch node.Type {
		case md.Document:
		case md.BlockQuote:
			lines := bytes.SplitAfter(node.Literal, []byte("\n"))
			for _, line := range lines {
				if !bytes.Equal(line, []byte("")) {
					r.buf.WriteString("┃   ")
					r.buf.Write(line)
				}
			}
			r.buf.WriteString("\n")
		case md.Paragraph:
			if entering {
			} else {
				r.buf.WriteString("\n")
				r.buf.WriteString("\n")
			}
		case md.Header:
			if entering {
				r.buf.WriteString("\n")
				r.buf.WriteString("\033[1m")
				r.buf.WriteString(strings.Repeat("#", node.HeaderData.Level))
				r.buf.WriteString(" ")
				r.buf.Write(node.Literal)
			} else {
				r.buf.WriteString("\033[0m")
				r.buf.WriteString("\n")
				r.buf.WriteString("\n")
			}
		case md.HorizontalRule:
			r.buf.WriteString("\033[2m")
			r.buf.WriteString("---")
			r.buf.WriteString("\033[0m")
			r.buf.WriteString("\n")
		case md.Emph:
			if entering {
				r.buf.WriteString("\033[3m")
			} else {
				r.buf.WriteString("\033[0m")
			}
		case md.Strong:
			if entering {
				r.buf.WriteString("\033[1m")
			} else {
				r.buf.WriteString("\033[0m")
			}
		case md.Del:
			if entering {
				r.buf.WriteString("\033[2m")
				r.Write([]byte("~~"))
			} else {
				r.Write([]byte("~~"))
				r.buf.WriteString("\033[0m")
			}
		case md.Link:
			if entering {
				r.inLinkNode = true
				r.buf.WriteString("\033[4m")
			} else if r.inLinkNode {
				r.inLinkNode = false
				r.Write(node.LinkData.Destination)
				r.buf.WriteString("\033[0m")
			}
		case md.Image:
			if entering {
				r.buf.WriteString("\033[4m")
				r.Write(node.LinkData.Destination)
			} else {
				r.inLinkNode = false
				r.buf.WriteString("\033[0m")
			}
		case md.Text:
			if !r.inLinkNode {
				r.Write(node.Literal)
			}
		case md.CodeBlock:
			lines := bytes.SplitAfter(node.Literal, []byte("\n"))
			r.buf.WriteString("\033[2m")
			for _, line := range lines {
				if !bytes.Equal(line, []byte("")) {
					r.buf.WriteString("┃   ")
					r.buf.Write(line)
				}
			}
			r.buf.WriteString("\033[0m")
			r.buf.WriteString("\n")
		case md.Code:
			r.buf.WriteString("\033[2m")
			r.Write(node.Literal)
			r.buf.WriteString("\033[0m")
		case md.List, md.Item,
			md.HTMLBlock, md.HTMLSpan,
			md.Softbreak, md.Hardbreak,
			md.Table, md.TableCell, md.TableHead, md.TableBody, md.TableRow:
			// TODO(uwe): Implement
			r.buf.WriteString("\033[31m")
			r.Write(node.Literal)
			r.buf.WriteString("\033[0m")
		}
		return md.GoToNext
	})
	return r.buf.Bytes()
}

func (r *ansiRenderer) Write(b []byte) {
	r.buf.Write(bytes.Replace(b, []byte("\n"), []byte(" "), -1))
}
