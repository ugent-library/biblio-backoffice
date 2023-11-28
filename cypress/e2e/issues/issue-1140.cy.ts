// https://github.com/ugent-library/biblio-backoffice/issues/1140

describe('Issue #1140: External contributor info is empty in the suggest box', () => {
  it('should display the external contributor name in the suggestions', () => {
    cy.loginAsResearcher()

    cy.setUpPublication('Book')

    cy.visitPublication()

    cy.updateFields(
      'Authors',
      () => {
        cy.intercept({
          pathname: '/publication/*/contributors/author/suggestions',
          query: {
            first_name: 'John',
            last_name: 'Doe',
          },
        }).as('suggestions')

        cy.setFieldByLabel('First name', 'John')
        cy.setFieldByLabel('Last name', 'Doe')

        cy.wait('@suggestions')

        cy.contains('#person-suggestions .list-group-item', 'John Doe').contains('.btn', 'Add external author').click()

        cy.setFieldByLabel('Roles', 'Validation')
      },
      /^Save$/
    )

    cy.contains('table#contributors-author-table tr', 'John Doe').find('.if.if-edit').click()

    cy.get('#person-suggestions')
      .find('.list-group-item')
      .should('have.length', 1)
      .find('.bc-avatar-text')
      .should('contain', 'Current selection')
      .should('contain', 'External, non-UGent')
      .should('contain', 'John Doe')
  })
})
