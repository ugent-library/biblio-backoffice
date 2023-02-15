describe('Publication import', () => {
  it('should be a possible to import publications from Web of Science and save as draft', () => {
    cy.loginAsResearcher()

    cy.visit('/')

    cy.contains('Biblio Publications').click()

    cy.contains('Add Publication').click()

    // Add publication(s)
    cy.contains('Step 1').should('be.visible')
    cy.contains('.bc-toolbar-title', 'Start: add publication(s)').should('be.visible')
    cy.get('.c-stepper__step').as('steps').should('have.length', 3)
    cy.get('@steps').eq(0).should('have.class', 'c-stepper__step--active')
    cy.get('@steps').eq(1).should('not.have.class', 'c-stepper__step--active')
    cy.get('@steps').eq(2).should('not.have.class', 'c-stepper__step--active')

    cy.contains('Import from Web of Science').click()
    cy.contains('.btn', 'Add publication(s)').click()

    // Upload WoS file
    cy.get('.c-file-upload').should('contain.text', 'Drag and drop or')
    cy.contains('.btn', 'Upload file').get('.spinner-border').should('not.be.visible')
    cy.get('input[name=file]').selectFile('cypress/fixtures/import-from-wos.txt')
    cy.contains('.btn', 'Upload file').get('.spinner-border').should('be.visible')

    // Review and publish
    cy.contains('Step 2').should('be.visible')
    cy.contains('.bc-toolbar-title', 'Review and publish').should('be.visible')
    cy.get('@steps').eq(0).should('have.class', 'c-stepper__step--complete')
    cy.get('@steps').eq(1).should('have.class', 'c-stepper__step--active')
    cy.get('@steps').eq(2).should('not.have.class', 'c-stepper__step--active')

    cy.contains('Review and publish').should('be.visible')
    cy.contains('Imported publications Showing 3').should('be.visible')

    // Delete 2 publications
    deletePublication(
      'Enhancing bioflocculation in high-rate activated sludge improves effluent quality yet increases sensitivity to surface overflow rate'
    )
    cy.contains('Imported publications Showing 2').should('be.visible')

    deletePublication('Fusarium isolates from Belgium causing wilt in lettuce show genetic and pathogenic diversity')
    cy.contains('Imported publications Showing 1').should('be.visible')

    // Extract Biblio ID for remaining publication
    cy.get('.list-group-item-main')
      .should('have.length', 1)
      .contains('Biblio ID:')
      .find('.c-code')
      .invoke('text')
      .as('biblioID', { type: 'static' })

    // Try publishing remaining publication and verify validation error
    cy.ensureNoModal()

    cy.contains('.btn', 'Publish all to Biblio').click()

    cy.ensureModal('Unable to publish a publication due to the following errors').within(() => {
      cy.get('.alert.alert-danger').should('be.visible').should('contain.text', 'At least one UGent author is required')

      cy.contains('.btn', 'Close').click()
    })

    // Add internal UGent author
    cy.ensureNoModal()

    cy.contains('People & Affiliations').click()

    cy.ensureNoModal()

    cy.contains('.btn', 'Add author').click({ scrollBehavior: false })

    cy.ensureModal('Add author').within(function () {
      cy.intercept({
        pathname: `/publication/${this.biblioID}/contributors/author/suggestions`,
        query: {
          first_name: 'Dries',
          last_name: 'Moreels',
        },
      }).as('user-search')

      cy.contains('Search author').should('be.visible')
      cy.get('input[name=first_name]').type('Dries')
      cy.get('input[name=last_name]').type('Moreels')

      cy.wait('@user-search')

      cy.contains('.badge', 'Active UGent member')
        .closest('.list-group-item')
        // Make sure the right author is selected
        .should('contain.text', 'Dries Moreels')
        .should('contain.text', '802001088860')
        .contains('.btn', 'Add author')
        .click()

      cy.contains('Review author information').should('be.visible')

      cy.get('.list-group-item').should('have.length', 1).should('contain.text', 'Dries Moreels')

      // Using RegExp to match entire button text to make sure the "Save and add next" button is not picked instead
      cy.contains('.btn', /^Save$/).click()

      // Necessary to wait or next publication attempt may still be invalid
      cy.wait(2000)
    })

    cy.ensureNoModal()

    cy.contains('.btn', 'Back to "Review and publish" overview').click()

    // Verify publication is still draft
    cy.get('.list-group-item .badge')
      .should('have.class', 'badge-secondary')
      .find('.badge-text')
      .should('have.text', 'Biblio Draft')

    // Publish
    cy.intercept('POST', '/publication/add-multiple/*/publish').as('publish')
    cy.contains('.btn', 'Publish all to Biblio').click()
    cy.wait('@publish')

    // Finished
    cy.contains('.bc-toolbar-title', 'Congratulations!').should('be.visible')
    cy.get('@steps').eq(0).should('have.class', 'c-stepper__step--complete')
    cy.get('@steps').eq(1).should('have.class', 'c-stepper__step--complete')
    cy.get('@steps').eq(2).should('have.class', 'c-stepper__step--active')
    cy.contains('Imported publications Showing 1').should('be.visible')

    cy.contains('.btn', 'Continue to overview').click()
    cy.location('pathname').should('eq', '/publication')

    // Verify publication is published
    cy.get('@biblioID').then(biblioID => cy.visit(`/publication/${biblioID}`))

    cy.get('#summary .badge')
      .should('have.class', 'badge-default')
      .find('.badge-text')
      .should('have.text', 'Biblio Public')
  })

  // TODO: Not yet implemented
  // Example publication: "How can we possibly resolve the planet's nitrogen dilemma?" in import-from-wos.txt
  it('should report errors after import')
})

function deletePublication(title) {
  cy.ensureNoModal()

  cy.contains('.list-group-item-title', title)
    .closest('.list-group-item')
    .find('.c-button-toolbar')
    // The "..." dropdown toggle button
    .find('.dropdown .btn:has(i.if.if-more)')
    .click()
    .closest('.dropdown')
    .contains('button', 'Delete')
    .click()

  cy.ensureModal('Are you sure?').within(() => {
    cy.get('.modal-body > p').should('have.text', 'Are you sure you want to delete this publication?')

    cy.contains('.modal-footer .btn', 'Delete').click()
  })
}
