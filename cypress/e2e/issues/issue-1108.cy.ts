// https://github.com/ugent-library/biblio-backoffice/issues/1108

describe('Issue #1108: Cannot add author without first name', () => {
  beforeEach(() => {
    cy.loginAsResearcher()

    cy.visit('/publication/add')

    cy.contains('Import from Web of Science').click()
    cy.contains('.btn', 'Add publication(s)').click()

    cy.get('input[name=file]').selectFile('cypress/fixtures/wos-000963572100001.txt')

    cy.contains('People & Affiliations').click()

    cy.contains('.btn', 'Add author').click()
  })

  it('should be possible to add author without first name', () => {
    cy.ensureModal('Add author').within(() => {
      cy.setFieldByLabel('Last name', 'Doe')

      cy.contains('Doe External, non-UGent').closest('.list-group-item').contains('.btn', 'Add external author').click()
    })

    cy.ensureModal('Add author').within(() => {
      cy.contains('.btn', /^Save$/).click()
    })

    cy.ensureNoModal()

    cy.get('.card#authors')
      .find('#contributors-author-table tr')
      .last()
      .find('td')
      .first()
      .find('.bc-avatar-text')
      .should('contain', '[missing] Doe')
  })

  it('should be possible to add author without last name', () => {
    cy.ensureModal('Add author').within(() => {
      cy.setFieldByLabel('First name', 'John')

      cy.contains('John External, non-UGent')
        .closest('.list-group-item')
        .contains('.btn', 'Add external author')
        .click()
    })

    cy.ensureModal('Add author').within(() => {
      cy.contains('.btn', /^Save$/).click()
    })

    cy.ensureNoModal()

    cy.get('.card#authors')
      .find('#contributors-author-table tr')
      .last()
      .find('td')
      .first()
      .find('.bc-avatar-text')
      .should('contain', 'John [missing]')
  })
})
