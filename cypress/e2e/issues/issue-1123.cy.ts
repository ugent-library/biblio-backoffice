describe('Issue #1123:  WoS import cuts keywords up because of newlines in import', () => {
  it('should not split up keywords by newlines', () => {
    cy.loginAsResearcher()

    cy.visit('/publication/add')

    cy.contains('Import from Web of Science').click()
    cy.contains('.btn', 'Add publication(s)').click()

    cy.get('input[name=file]').selectFile('cypress/fixtures/wos-000963572100001.txt')

    cy.contains('a.list-group-item', 'Description').click()

    cy.contains('label', 'Keywords')
      .next('div')
      .find('ul > li > span')
      .then($spans => Cypress._.map($spans, 'innerText'))
      .should('include.members', ['testosterone treatment', 'in vitro maturation'])
      .should('not.include', 'testosterone')
      .should('not.include', 'treatment')
      .should('not.include', 'in vitro')
      .should('not.include', 'maturation')
  })
})
