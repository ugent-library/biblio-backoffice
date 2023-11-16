// https://github.com/ugent-library/biblio-backoffice/issues/1309

describe('Issue #1309: Cannot return to profile after using "view as"', () => {
  it('should be possible for a researcher to return to their own profile when viewing as another user', () => {
    const LIBRARIAN_NAME = Cypress.env('LIBRARIAN_NAME')

    cy.loginAsLibrarian()

    cy.visit('/')

    cy.contains('.bc-avatar-and-text:visible', LIBRARIAN_NAME).click()

    cy.contains('View as').click()

    cy.ensureModal('View as other user').within(() => {
      cy.setFieldByLabel('First name', 'Dries')
      cy.setFieldByLabel('Last name', 'Moreels')

      cy.contains('802001088860')
        .should('be.visible')
        .closest('.list-group-item')
        .contains('.btn', 'Change user')
        .click()
    })

    cy.ensureNoModal()

    cy.contains('Viewing the perspective of Dries Moreels.').should('be.visible')
    cy.get('.bc-avatar-and-text:visible').should('contain.text', 'Dries Moreels')

    cy.contains('.btn', `return to ${LIBRARIAN_NAME}`).click()

    cy.contains('Viewing the perspective of').should('not.exist')
    cy.get('.bc-avatar-and-text:visible').should('contain.text', LIBRARIAN_NAME)
  })
})
