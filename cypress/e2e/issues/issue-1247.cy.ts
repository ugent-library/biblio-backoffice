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

  // TODO: add similar test for datasets
  it(`should fully display the user menu on all pages during manual publication set-up`, () => {
    cy.visit('/publication/add?method=manual')

    assertUserMenuWorks()

    cy.contains('Miscellaneous').click()
    cy.contains('.btn', 'Add publication(s)').click()

    cy.location('pathname').should('eq', '/publication/add-single/import/confirm')

    assertUserMenuWorks()

    cy.updateFields(
      'Publication details',
      () => {
        cy.contains('label', 'Title').click().type('Test publication [CYPRESSTEST]')
        cy.contains('label', 'Publication year').click().type(new Date().getFullYear().toString())
      },
      true
    )

    cy.updateFields(
      'Authors',
      () => {
        cy.contains('label', 'First name').click().type('Dries')
        cy.contains('label', 'Last name').click().type('Moreels')

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
