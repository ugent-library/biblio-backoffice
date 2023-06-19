describe('Issue #1124:  Add friendlier consistent confirmation toaster when locking or unlocking a record', () => {
  it('should not split up keywords by newlines', () => {
    cy.loginAsResearcher()

    cy.visit('/publication/add')

    cy.contains('Import from Web of Science').click()
    cy.contains('.btn', 'Add publication(s)').click()

    cy.get('input[name=file]').selectFile('cypress/fixtures/wos-000863180800005.txt')

    cy.contains('Biblio ID:')
      .find('code')
      .invoke('text')
      .as('biblioId')
      .then(biblioId => {
        // Publication does not have lock icon in the list view
        cy.visit('/publication')
        cy.get('input[name=q][placeholder="Search..."]').type(biblioId)
        cy.contains('.btn', 'Search').click()
        cy.contains('Publications Showing 1').should('be.visible')
        cy.get('.list-group-item .if.if-lock').should('not.exist')
        cy.contains('.list-group-item', 'Locked').should('not.exist')

        // Publication does not have lock icon in the detail view
        cy.visit(`/publication/${biblioId}`)
        cy.get('.c-subline .if.if-lock').should('not.exist')
        cy.contains('.c-subline', 'Locked').should('not.exist')

        // Edit button is still available for researchers for unlocked publications
        cy.contains('.btn', 'Edit').should('be.visible').click()
        cy.ensureModal('Edit publication details')

        // Publication cannot be locked by researcher
        cy.contains('.btn', 'Lock record').should('not.exist')

        // Publication can be locked by librarian
        cy.loginAsLibrarian()
        cy.switchMode('Librarian')

        cy.visit('/publication')
        cy.get('input[name=q][placeholder="Search..."]').type(biblioId)
        cy.contains('.btn', 'Search').click()
        cy.contains('Publications Showing 1').should('be.visible')
        cy.get('.list-group-item .if.if-lock').should('not.exist')
        cy.contains('.list-group-item', 'Locked').should('not.exist')

        cy.visit(`/publication/${biblioId}`)
        cy.contains('.btn', 'Lock record').should('be.visible').click()

        // Confirmation toast is displayed upon locking (and hidden after 5s)
        cy.contains('.toast', 'Publication was successfully locked.').as('lockedToast').should('be.visible')
        cy.wait(5000)
        cy.get('@lockedToast').should('not.exist')

        // Publication now has lock icon in the detail view
        cy.get('.c-subline .if.if-lock').should('be.visible')
        cy.contains('.c-subline', 'Locked').should('be.visible')

        // Edit button is still available for librarians for locked publications
        cy.contains('.btn', 'Edit').should('be.visible').click()
        cy.ensureModal('Edit publication details')

        // Publication does have lock icon in the list view for librarians
        cy.visit('/publication', { qs: { q: biblioId } })
        cy.contains('Publications Showing 1').should('be.visible')
        cy.get('.list-group-item .if.if-lock').should('be.visible')
        cy.contains('.list-group-item', 'Locked').should('be.visible')

        // Publication does have lock icon in the list view for researcher
        cy.loginAsResearcher()
        cy.visit('/publication', { qs: { q: biblioId, 'f[scope]': 'all' } })
        cy.contains('Publications Showing 1').should('be.visible')
        cy.get('.list-group-item .if.if-lock').should('be.visible')
        cy.contains('.list-group-item', 'Locked').should('be.visible')

        // In the detail view, as a researcher...
        cy.visit(`/publication/${biblioId}`)

        // Publication now has a lock icon
        cy.get('.c-subline .if.if-lock').should('be.visible')
        cy.contains('.c-subline', 'Locked').should('be.visible')

        // Edit button is no longer available
        cy.contains('.btn', 'Edit').should('not.exist')

        // Now go back as librarian to unlock the publication
        cy.loginAsLibrarian()
        cy.switchMode('Librarian')

        cy.visit(`/publication/${biblioId}`)
        cy.contains('.btn', 'Unlock record').should('be.visible').click()

        // Confirmation toast is displayed upon locking (and hidden after 5s)
        cy.contains('.toast', 'Publication was successfully unlocked.').as('unlockedToast').should('be.visible')
        cy.wait(5000)
        cy.get('@unlockedToast').should('not.exist')

        // Lock icon is removed again
        cy.get('.c-subline .if.if-lock').should('not.exist')
        cy.contains('.c-subline', 'Locked').should('not.exist')
      })
  })
})
