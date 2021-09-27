package engine

func (e *Engine) PublicationTypes() []string {
	return []string{
		"journal_article",
		"book",
		"book_chapter",
		"book_editor",
		"issue_editor",
		"conference",
		"dissertation",
		"miscellaneous",
		"report",
		"preprint",
	}
}
