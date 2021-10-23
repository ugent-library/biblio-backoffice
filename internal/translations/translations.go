package translations

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func init() {
	message.SetString(language.English, "credit_roles.first_author", "First author")
	message.SetString(language.English, "credit_roles.last_author", "Last author")
	message.SetString(language.English, "credit_roles.conceptualization", "Conceptualization")
	message.SetString(language.English, "credit_roles.data_curation", "Datacuration")
	message.SetString(language.English, "credit_roles.formal_analysis", "Formala nalysis")
	message.SetString(language.English, "credit_roles.funding_acquisition", "Funding acquisition")
	message.SetString(language.English, "credit_roles.investigation", "Investigation")
	message.SetString(language.English, "credit_roles.methodology", "Methodology")
	message.SetString(language.English, "credit_roles.project_administration", "Project administration")
	message.SetString(language.English, "credit_roles.resources", "Resources")
	message.SetString(language.English, "credit_roles.software", "Software")
	message.SetString(language.English, "credit_roles.supervision", "Supervision")
	message.SetString(language.English, "credit_roles.validation", "Validation")
	message.SetString(language.English, "credit_roles.visualization", "Visualization")
	message.SetString(language.English, "credit_roles.writing_original_draft", "Writing - original draft")
	message.SetString(language.English, "credit_roles.writing_review_editing", "Writing - review & editing")

	message.SetString(language.English, "publication_sorts.year.desc", "Year (newest first)")
	message.SetString(language.English, "publication_sorts.date_created.desc", "Added (newest first)")
	message.SetString(language.English, "publication_sorts.date_updated.desc", "Updated (newest first)")

	message.SetString(language.English, "publication_types.journal_article", "Journal article")
	message.SetString(language.English, "publication_types.book", "Journal article")
	message.SetString(language.English, "publication_types.book_chapter", "Book chapter")
	message.SetString(language.English, "publication_types.book_editor", "Book editor")
	message.SetString(language.English, "publication_types.issue_editor", "Issue editor")
	message.SetString(language.English, "publication_types.conference", "Conference")
	message.SetString(language.English, "publication_types.dissertation", "Dissertation")
	message.SetString(language.English, "publication_types.miscellaneous", "Miscellaneous")
	message.SetString(language.English, "publication_types.report", "Report")
	message.SetString(language.English, "publication_types.preprint", "Preprint")

	message.SetString(language.English, "journal_article_types.original", "Original")
	message.SetString(language.English, "journal_article_types.review", "Review")
	message.SetString(language.English, "journal_article_types.letter_note", "Letter note")
	message.SetString(language.English, "journal_article_types.proceedingsPaper", "Proceedings Paper")

	message.SetString(language.English, "conference_types.proceedingsPaper", "Preprint")
	message.SetString(language.English, "conference_types.abstract", "Abstract")
	message.SetString(language.English, "conference_types.poster", "Poster")
	message.SetString(language.English, "conference_types.other", "Other")

	message.SetString(language.English, "miscellaneous_types.artReview", "Art review")
	message.SetString(language.English, "miscellaneous_types.artisticWork", "Artistic Work")
	message.SetString(language.English, "miscellaneous_types.bibliography", "Bibliography")
	message.SetString(language.English, "miscellaneous_types.biography", "Biography")
	message.SetString(language.English, "miscellaneous_types.blogPost", "Blogpost")
	message.SetString(language.English, "miscellaneous_types.bookReview", "Book review")
	message.SetString(language.English, "miscellaneous_types.correction", "Correction")
	message.SetString(language.English, "miscellaneous_types.dictionaryEntry", "Dictionary entry")
	message.SetString(language.English, "miscellaneous_types.editorialMaterial", "Editorial material")
	message.SetString(language.English, "miscellaneous_types.encyclopediaEntry", "Encyclopedia entry")
	message.SetString(language.English, "miscellaneous_types.exhibitionReview", "Exhibition review")
	message.SetString(language.English, "miscellaneous_types.filmReview", "Film review")
	message.SetString(language.English, "miscellaneous_types.lectureSpeech", "Lecture speech")
	message.SetString(language.English, "miscellaneous_types.lemma", "Lemma")
	message.SetString(language.English, "miscellaneous_types.magazinePiece", "Magazine piece")
	message.SetString(language.English, "miscellaneous_types.manual", "Manual")
	message.SetString(language.English, "miscellaneous_types.musicEdition", "Music edition")
	message.SetString(language.English, "miscellaneous_types.musicReview", "Music review")
	message.SetString(language.English, "miscellaneous_types.newsArticle", "News article")
	message.SetString(language.English, "miscellaneous_types.newspaperPiece", "Newspaper piece")
	message.SetString(language.English, "miscellaneous_types.other", "Other")
	message.SetString(language.English, "miscellaneous_types.preprint", "Preprint")
	message.SetString(language.English, "miscellaneous_types.productReview", "Product review")
	message.SetString(language.English, "miscellaneous_types.report", "Report")
	message.SetString(language.English, "miscellaneous_types.technicalStandard", "Technical standard")
	message.SetString(language.English, "miscellaneous_types.textEdition", "Text edition")
	message.SetString(language.English, "miscellaneous_types.textTranslation", "Text translation")
	message.SetString(language.English, "miscellaneous_types.theatreReview", "Theatre review")
	message.SetString(language.English, "miscellaneous_types.workingPaper", "Working paper")

	message.SetString(language.English, "publication_classifications.A1", "A1 - article in the Web of Science")
	message.SetString(language.English, "publication_classifications.A2", "A2 - article published in an international peer reviewed journal")
	message.SetString(language.English, "publication_classifications.A3", "A3 - article published in a national peer reviewed journal")
	message.SetString(language.English, "publication_classifications.A4", "A4 - article published in a journal (not A1, A2 or A3)")
	message.SetString(language.English, "publication_classifications.B1", "B1 - author or coauthor of a book")
	message.SetString(language.English, "publication_classifications.B2", "B2 - author or coauthor of a book chapter")
	message.SetString(language.English, "publication_classifications.B3", "B3 - editor of a book or journal issue")
	message.SetString(language.English, "publication_classifications.C1", "C1 - conference paper (not WoS)")
	message.SetString(language.English, "publication_classifications.C3", "C3 - meeting abstract")
	message.SetString(language.English, "publication_classifications.D1", "D1 - doctoral thesis")
	message.SetString(language.English, "publication_classifications.D2", "D2 - student thesis")
	message.SetString(language.English, "publication_classifications.P1", "P1 - proceedings paper (WoS)")
	message.SetString(language.English, "publication_classifications.V", "V - miscellaneous")
	message.SetString(language.English, "publication_classifications.U", "U - unknown")

	message.SetString(language.English, "form_builder.classification", "Classification")
	message.SetString(language.English, "form_builder.doi", "DOI")
	message.SetString(language.English, "form_builder.isbn", "ISBN")
	message.SetString(language.English, "form_builder.issn", "ISSN")
	message.SetString(language.English, "form_builder.language", "Languages")
	message.SetString(language.English, "form_builder.title", "Title")
}
