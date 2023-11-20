// https://github.com/ugent-library/biblio-backoffice/issues/1247

describe('Issue #1247: User menu popup hidden behind publication details', () => {
  const testCases = {
    '/': 'home page',
    '/publication': 'publications page',
    '/publication/add': 'add publication page',
    '/publication/add?method=wos': 'add Web of Science publication page',
    '/publication/add?method=identifier': 'add publication from identifier page',
    '/publication/add?method=manual': 'add manual publication page',
    '/publication/add?method=bibtex': 'add BibTeX publication page',
    '/dataset': 'datasets page',
    '/dataset/add': 'add dataset page',
    '/dashboard/publications/faculties': 'publications - faculties dashboard page',
    '/dashboard/publications/socs': 'publications - SOCs dashboard page',
    '/dashboard/datasets/faculties': 'datasets - faculties dashboard page',
    '/dashboard/datasets/socs': 'datasets - SOCs dashboard page',
  }

  beforeEach(() => {
    cy.loginAsLibrarian()

    cy.switchMode('Librarian')
  })

  Object.entries(testCases).forEach(([path, name]) => {
    it(`should fully display the user menu on the ${name}`, () => {
      cy.visit(path)

      assertUserMenuWorks()
    })
  })

  it('should fully display the user menu on all pages during manual publication set-up', () => {
    cy.visit('/publication/add?method=manual')

    assertUserMenuWorks()

    cy.contains('Miscellaneous').click()
    cy.contains('.btn', 'Add publication(s)').click()

    cy.location('pathname').should('eq', '/publication/add-single/import/confirm')

    assertUserMenuWorks()

    cy.updateFields(
      'Publication details',
      () => {
        cy.setFieldByLabel('Title', 'Test publication [CYPRESSTEST]')
        cy.setFieldByLabel('Publication year', new Date().getFullYear().toString())
      },
      true
    )

    cy.updateFields(
      'Authors',
      () => {
        cy.setFieldByLabel('First name', 'Dries')
        cy.setFieldByLabel('Last name', 'Moreels')

        cy.contains('.btn', 'Add author').click()
      },
      /^Save$/
    )

    cy.contains('.btn', 'Complete Description').click()

    cy.location('pathname').should('match', new RegExp('/publication/\\w+/add/confirm'))

    cy.extractBiblioId()

    assertUserMenuWorks()

    cy.contains('.btn', 'Publish to Biblio').click()

    cy.location('pathname').should('match', new RegExp('/publication/\\w+/add/finish'))

    assertUserMenuWorks()

    cy.contains('.btn', 'Continue to overview').click()

    cy.location('pathname').should('eq', '/publication')

    assertUserMenuWorks()

    cy.visitPublication()

    cy.get('@biblioId').then(biblioId => {
      cy.location('pathname').should('eq', `/publication/${biblioId}`)
    })

    assertUserMenuWorks()
  })

  it('should fully display the user menu on all pages during manual dataset set-up', () => {
    cy.visit('/dataset/add')

    assertUserMenuWorks()

    cy.contains('Register a dataset manually').click('left')
    cy.contains('.btn', 'Add dataset').click()

    cy.location('pathname').should('eq', '/dataset/add')

    assertUserMenuWorks()

    cy.updateFields(
      'Dataset details',
      () => {
        cy.setFieldByLabel('Title', `Test dataset [CYPRESSTEST]`)
        cy.setFieldByLabel('Persistent identifier type', 'DOI')
        cy.setFieldByLabel('Identifier', '10.5072/test/t')

        cy.setFieldByLabel('Publication year', new Date().getFullYear().toString())
        cy.setFieldByLabel('Data format', 'text/csv').next('.autocomplete-hits').contains('.badge', 'text/csv').click()
        cy.setFieldByLabel('Publisher', 'UGent')
        cy.setFieldByLabel('License', 'CC0 (1.0)')
        cy.setFieldByLabel('Access level', 'Open access')
      },
      true
    )

    cy.updateFields(
      'Creators',
      () => {
        cy.setFieldByLabel('First name', 'Dries')
        cy.setFieldByLabel('Last name', 'Moreels')

        cy.contains('.btn', 'Add creator').click()
      },
      /^Save$/
    )

    cy.contains('.btn', 'Complete Description').click()

    cy.location('pathname').should('match', new RegExp('/dataset/\\w+/add/confirm'))

    cy.extractBiblioId()

    assertUserMenuWorks()

    cy.contains('.btn', 'Publish to Biblio').click()

    cy.location('pathname').should('match', new RegExp('/dataset/\\w+/add/finish'))

    assertUserMenuWorks()

    cy.contains('.btn', 'Continue to overview').click()

    cy.location('pathname').should('eq', '/dataset')

    assertUserMenuWorks()

    cy.visitDataset()

    cy.get('@biblioId').then(biblioId => {
      cy.location('pathname').should('eq', `/dataset/${biblioId}`)
    })

    assertUserMenuWorks()
  })

  function assertUserMenuWorks() {
    cy.get('.nav-main .dropdown-menu').as('userMenu').should('not.be.visible')

    cy.get('.nav-main .bc-avatar .if-user:visible').as('userName', { type: 'static' }).click()

    cy.get('@userMenu')
      .should('be.visible')
      .within(() => {
        cy.get('.bc-avatar-and-text').should('be.visible')
        cy.get('.dropdown-divider').should('be.visible')
        cy.contains('.dropdown-item', 'Logout').should('be.visible')
      })

    cy.get('@userName').click()

    cy.get('@userMenu').should('not.be.visible')
  }
})
