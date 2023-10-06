// https://github.com/ugent-library/biblio-backoffice/issues/1246

describe('Issue #1246: Close button on toast does not work', () => {
  beforeEach(() => {
    cy.loginAsResearcher()
  })

  it('should be possible to dismiss the delete publication toast', () => {
    setUpPublication()

    cy.contains('.btn', 'Save as draft').click()

    cy.visitPublication()

    cy.get('.bc-toolbar')
      // The "..." dropdown toggle button
      .find('.dropdown .btn:has(i.if.if-more)')
      .click()

    cy.contains('Delete').click()

    cy.ensureModal('Are you sure?').within(() => {
      cy.contains('.btn-danger', 'Delete').click()
    })

    cy.ensureNoModal()

    assertToast('Publication was successfully deleted.')
  })

  it('should be possible to dismiss the publish publication toast', () => {
    setUpPublication()

    cy.contains('.btn', 'Save as draft').click()

    cy.visitPublication()

    cy.contains('Publish to Biblio').click()

    cy.ensureModal('Are you sure?').within(() => {
      cy.contains('.btn-success', 'Publish').click()
    })

    cy.ensureNoModal()

    assertToast('Publication was successfully published.')
  })

  it('should be possible to dismiss the withdraw publication toast', () => {
    setUpPublication()

    cy.contains('.btn', 'Publish to Biblio').click()

    cy.visitPublication()

    cy.contains('Withdraw').click()

    cy.ensureModal('Are you sure?').within(() => {
      cy.contains('.btn-danger', 'Withdraw').click()
    })

    cy.ensureNoModal()

    assertToast('Publication was successfully withdrawn.')
  })

  it('should be possible to dismiss the republish publication toast', () => {
    setUpPublication()

    cy.contains('.btn', 'Publish to Biblio').click()

    cy.visitPublication()

    cy.contains('Withdraw').click()

    cy.ensureModal('Are you sure?').within(() => {
      cy.contains('.btn-danger', 'Withdraw').click()
    })

    cy.ensureNoModal()

    cy.contains('Republish').click()

    cy.ensureModal('Are you sure?').within(() => {
      cy.contains('.btn-success', 'Republish').click()
    })

    cy.ensureNoModal()

    assertToast('Publication was successfully republished.')
  })

  it('should be possible to dismiss the locked publication toast', () => {
    setUpPublication()

    cy.contains('.btn', 'Publish to Biblio').click()

    cy.loginAsLibrarian()
    cy.switchMode('Librarian')

    cy.visitPublication()

    cy.contains('Lock record').click()

    assertToast('Publication was successfully locked.')
  })

  it('should be possible to dismiss the locked publication toast', () => {
    setUpPublication()

    cy.contains('.btn', 'Publish to Biblio').click()

    cy.loginAsLibrarian()
    cy.switchMode('Librarian')

    cy.visitPublication()

    cy.contains('Lock record').click()

    // Make sure lock-toast is gone first
    cy.get('.toast', { timeout: 6000 }).should('not.exist')

    cy.contains('Unlock record').click()

    assertToast('Publication was successfully unlocked.')
  })

  function setUpPublication() {
    cy.visit('/publication/add')
    cy.contains('Enter a publication manually').click()
    cy.contains('.btn', 'Add publication(s)').click()

    cy.contains('Miscellaneous').click()
    cy.contains('.btn', 'Add publication(s)').click()

    cy.contains('Publication details').closest('.card-header').contains('.btn', 'Edit').click()

    cy.ensureModal('Edit publication details').within(() => {
      cy.get('input[type=text][name=title]').type('Issue 1246 test [CYPRESSTEST]')
      cy.get('input[type=text][name=year]').type(new Date().getFullYear().toString())

      cy.contains('.btn', 'Save').click()
    })

    cy.contains('People & Affiliations').click()

    cy.contains('.btn', 'Add author').click()

    cy.ensureModal('Add author').within(() => {
      cy.get('input[type=text][name=first_name]').type('Dries')
      cy.get('input[type=text][name=last_name]').type('Moreels')

      cy.contains('.btn', 'Add author').click()
    })

    cy.ensureModal('Add author').within(() => {
      cy.contains('.btn', /^Save$/).click()
    })

    cy.ensureNoModal()

    cy.contains('.btn', 'Complete Description').click()

    cy.extractBiblioId()
  }

  function assertToast(toastMessage: string) {
    cy.contains('.toast', toastMessage)
      .should('be.visible')
      .within(() => {
        cy.get('.btn-close').click()
      })

    // Reduced assertion timeout here so the test still works if someone decides to reduce the
    // toast dismissal timeout in the future.
    cy.contains('.toast', toastMessage, { timeout: 1000 }).should('not.exist')
  }
})
