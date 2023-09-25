describe('Issue #1140: External contributor info is empty in the suggest box', () => {
  it('should display the external contributor name in the suggestions', () => {
    cy.loginAsResearcher()

    cy.visit('/publication/add')

    cy.contains('Import from Web of Science').click()
    cy.contains('.btn', 'Add publication(s)').click()

    cy.get('input[name=file]').selectFile('cypress/fixtures/wos-000963572100001.txt')

    cy.contains('People & Affiliations').click()

    cy.contains('.btn', 'Add author').click({ scrollBehavior: false })

    cy.ensureModal('Add author').within(() => {
      cy.get('input[name=first_name]').type('John')
      cy.get('input[name=last_name]').type('Doe')

      cy.contains('.btn', 'Add external author').click()
    })

    cy.ensureModal('Add author').within(() => {
      cy.contains('label', 'Roles').next().find('select').select('Validation')

      cy.contains('.btn', /^Save$/).click()
    })

    cy.ensureNoModal()

    cy.contains('table#contributors-author-table tr', 'John Doe')
      .find('.if.if-edit')
      .click({ scrollBehavior: 'nearest' })

    cy.get('#person-suggestions')
      .find('.list-group-item')
      .should('have.length', 1)
      .find('.bc-avatar-text')
      .should('contain', 'Current selection')
      .should('contain', 'External, non-UGent')
      .should('contain', 'John Doe')
  })
})
