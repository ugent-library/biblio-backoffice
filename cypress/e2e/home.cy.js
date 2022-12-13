describe('empty spec', () => {
  it('be able to load the home page anonymously', () => {
    cy.visit('/')

    cy.get('h2').should('have.text', 'Home')

    cy.contains('a', 'Log in').should('be.visible')

    cy.contains('h3', 'Publications').should('be.visible')
    cy.contains('a', 'Go to publications').should('be.visible')

    cy.contains('h3', 'Datasets').should('be.visible')
    cy.contains('a', 'Go to datasets').should('be.visible')

    cy.get('.u-scroll-wrapper__body').scrollTo('bottom', { duration: 250 })

    cy.contains('h3', 'Biblio Academic Bibliography').should('be.visible')
    cy.contains('a', 'Go to Biblio Academic Bibliography').should('be.visible')

    cy.contains('h3', 'Help').should('be.visible')
    cy.contains('a', 'How to register and deposit').should('be.visible')
  })
})
