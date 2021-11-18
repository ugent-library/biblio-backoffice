package engine

import "github.com/ugent-library/biblio-backend/internal/models"

// Add Author to publication
func (e *Engine) AddAuthorToPublication(pub *models.Publication, author *models.Contributor, delta int) {
	placeholder := models.Contributor{}
	authors := make([]models.Contributor, len(pub.Author))
	copy(authors, pub.Author)

	authors = append(authors, placeholder)
	copy(authors[delta+1:], authors[delta:])
	authors[delta] = *author

	pub.Author = authors
}

// Get Author from publication
func (e *Engine) GetAuthorFromPublication(pub *models.Publication, delta int) *models.Contributor {
	return &pub.Author[delta]
}

// // Update Author on publication
func (e *Engine) UpdateAuthorOnPublication(pub *models.Publication, contributor *models.Contributor, delta int) {
	authors := make([]models.Contributor, len(pub.Author))
	copy(authors, pub.Author)

	authors[delta] = *contributor

	pub.Author = authors
}

// // Remove Author from publication
func (e *Engine) RemoveAuthorFromPublication(pub *models.Publication, delta int) {
	authors := make([]models.Contributor, len(pub.Author))
	copy(authors, pub.Author)

	authors = append(authors[:delta], authors[delta+1:]...)

	pub.Author = authors
}
