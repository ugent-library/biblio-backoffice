package engine

// func (e *Engine) PublicationTypes() []string {
// 	return []string{
// 		"journal_article",
// 		"book",
// 		"book_chapter",
// 		"book_editor",
// 		"issue_editor",
// 		"conference",
// 		"dissertation",
// 		"miscellaneous",
// 		"report",
// 		"preprint",
// 	}
// }

// func (e *Engine) PublicationStatuses() []string {
// 	return []string{
// 		"new",
// 		"private",
// 		"submitted",
// 		"returned",
// 		"public",
// 		"deleted",
// 	}
// }

func (e *Engine) PublicationSorts() []string {
	return []string{
		"date_created.desc",
		"date_updated.desc",
		"year.desc",
	}
}
