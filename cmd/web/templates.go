package main

import "snippetbox.ab.net/internal/models"

// 添加一个 Snippets 字段
type templateData struct {
	Snippet  *models.Snippet
	Snippets []*models.Snippet
}
