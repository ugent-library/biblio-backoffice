// https://github.com/ugent-library/biblio-backoffice/issues/1140

describe('Issue #1140: External contributor info is empty in the suggest box', () => {
  it('should display the external contributor name in the suggestions', () => {
    cy.loginAsResearcher()

    cy.visit('/publication/add')

    cy.contains('Import from Web of Science').click()
    cy.contains('.btn', 'Add publication(s)').click()

    cy.get('input[name=file]').selectFile('cypress/fixtures/wos-000963572100001.txt')

    cy.contains('People & Affiliations').click()

    cy.contains('.btn', 'Add author').click()

    cy.ensureModal('Add author').within(() => {
      cy.setFieldByLabel('First name', 'John')
      cy.setFieldByLabel('Last name', 'Doe')

      cy.contains('.btn', 'Add external author').click()
    })

    cy.ensureModal('Add author')
      .within(() => {
        cy.setFieldByLabel('Roles', 'Validation')
      })
      .closeModal(/^Save$/)

    cy.ensureNoModal()

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
