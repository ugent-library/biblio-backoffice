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

    cy.contains('.dropdown-item', 'Delete').click()

    cy.ensureModal('Are you sure?').closeModal('Delete')

    cy.ensureNoModal()

    assertToast('Publication was successfully deleted.')
  })

  it('should be possible to dismiss the publish publication toast', () => {
    setUpPublication()

    cy.contains('.btn', 'Save as draft').click()

    cy.visitPublication()

    cy.contains('.btn', 'Publish to Biblio').click()

    cy.ensureModal('Are you sure?').closeModal('Publish')

    cy.ensureNoModal()

    assertToast('Publication was successfully published.')
  })

  it('should be possible to dismiss the withdraw publication toast', () => {
    setUpPublication()

    cy.contains('.btn', 'Publish to Biblio').click()

    cy.visitPublication()

    cy.contains('.btn', 'Withdraw').click()

    cy.ensureModal('Are you sure?').closeModal('Withdraw')

    cy.ensureNoModal()

    assertToast('Publication was successfully withdrawn.')
  })

  it('should be possible to dismiss the republish publication toast', () => {
    setUpPublication()

    cy.contains('.btn', 'Publish to Biblio').click()

    cy.visitPublication()

    cy.contains('.btn', 'Withdraw').click()

    cy.ensureModal('Are you sure?').closeModal('Withdraw')

    cy.ensureNoModal()

    cy.ensureToast('Publication was successfully withdrawn.').closeToast()

    // Make sure withdraw-toast is gone first
    cy.ensureNoToast()

    cy.contains('.btn', 'Republish').click()

    cy.ensureModal('Are you sure?').closeModal('Republish')

    cy.ensureNoModal()

    assertToast('Publication was successfully republished.')
  })

  it('should be possible to dismiss the locked publication toast', () => {
    setUpPublication()

    cy.contains('.btn', 'Publish to Biblio').click()

    cy.loginAsLibrarian()
    cy.switchMode('Librarian')

    cy.visitPublication()

    cy.contains('.btn', 'Lock record').click()

    assertToast('Publication was successfully locked.')
  })

  it('should be possible to dismiss the unlocked publication toast', () => {
    setUpPublication()

    cy.contains('.btn', 'Publish to Biblio').click()

    cy.loginAsLibrarian()
    cy.switchMode('Librarian')

    cy.visitPublication()

    cy.contains('.btn', 'Lock record').click()

    cy.ensureToast('Publication was successfully locked.').closeToast()

    // Make sure lock-toast is gone first
    cy.ensureNoToast()

    cy.contains('.btn', 'Unlock record').click()

    assertToast('Publication was successfully unlocked.')
  })

  function setUpPublication() {
    cy.visit('/publication/add')
    cy.contains('Enter a publication manually').click()
    cy.contains('.btn', 'Add publication(s)').click()

    cy.contains('Miscellaneous').click()
    cy.contains('.btn', 'Add publication(s)').click()

    cy.contains('Publication details').closest('.card-header').contains('.btn', 'Edit').click()

    cy.ensureModal('Edit publication details')
      .within(() => {
        cy.setFieldByLabel('Title', 'Issue 1246 test [CYPRESSTEST]')
        cy.setFieldByLabel('Publication year', new Date().getFullYear().toString())
      })
      .closeModal(true)

    cy.contains('People & Affiliations').click()

    cy.contains('.btn', 'Add author').click()

    cy.ensureModal('Add author').within(() => {
      cy.intercept('/publication/*/contributors/author/suggestions?*').as('suggestions')

      cy.setFieldByLabel('First name', 'Griet')
      cy.setFieldByLabel('Last name', 'Alleman')

      cy.wait('@suggestions')

      cy.intercept('/publication/*/contributors/author/confirm-create?*').as('confirmCreate')

      cy.contains('.btn', 'Add author').click()
    })

    cy.wait('@confirmCreate')

    cy.ensureModal('Add author').closeModal(/^Save$/)

    cy.ensureNoModal()

    cy.contains('.btn', 'Complete Description').click()

    cy.extractBiblioId()
  }

  function assertToast(toastMessage: string) {
    cy.ensureToast(toastMessage).closeToast()

    // Reduced assertion timeout here so the test still works if someone decides to reduce the
    // toast dismissal timeout in the future.
    cy.ensureNoToast({ timeout: 500 })
  }
})
