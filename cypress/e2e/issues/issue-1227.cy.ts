describe('Issue #1127: Cannot search any longer on book title, journal title, short journal title nor conference title', () => {
  beforeEach(() => {
    cy.loginAsResearcher()
  })

  it('should be possible to search by publisher', () => {
    const randomText = getRandomText()

    cy.setUpPublication('Miscellaneous')

    cy.visitPublication()

    cy.updateFields(
      'Publication details',
      () => {
        cy.setFieldByLabel('Publisher', `Alternative title: ${randomText}`)
      },
      true
    )

    cy.visit('/publication')

    cy.search(randomText).should('eq', 1)
  })

  it('should be possible to search by alternative title', () => {
    const randomText1 = getRandomText()
    const randomText2 = getRandomText()

    cy.setUpPublication('Miscellaneous')

    cy.visitPublication()

    cy.updateFields(
      'Publication details',
      () => {
        cy.setFieldByLabel('Alternative title', `Alternative title: ${randomText1}`)
          .closest('.form-value')
          .contains('.btn', 'Add')
          .click()
          .closest('.form-value')
          .next('.form-value')
          .find('input')
          .type(`Alternative title: ${randomText2}`)
      },
      true
    )

    cy.visit('/publication')

    cy.search(randomText1).should('eq', 1)
    cy.search(randomText2).should('eq', 1)
  })

  it('should be possible to search by conference title', () => {
    const randomText = getRandomText()

    cy.setUpPublication('Conference contribution')

    cy.visitPublication()

    cy.updateFields(
      'Conference details',
      () => {
        cy.setFieldByLabel('Conference', `The conference name: ${randomText}`)
      },
      true
    )

    cy.visit('/publication')

    cy.search(randomText).should('eq', 1)
  })

  it('should be possible to search by journal title', () => {
    const randomText = getRandomText()

    cy.setUpPublication('Journal Article')

    cy.visitPublication()

    cy.updateFields(
      'Publication details',
      () => {
        cy.setFieldByLabel('Journal title', `The journal name: ${randomText}`)
      },
      true
    )

    cy.visit('/publication')

    cy.search(randomText).should('eq', 1)
  })

  it('should be possible to search by short journal title', () => {
    const randomText = getRandomText()

    cy.setUpPublication('Journal Article')

    cy.visitPublication()

    cy.updateFields(
      'Publication details',
      () => {
        cy.setFieldByLabel('Short journal title', `The short journal name: ${randomText}`)
      },
      true
    )

    cy.visit('/publication')

    cy.search(randomText).should('eq', 1)
  })

  it('should be possible to search by book title', () => {
    const randomText = getRandomText()

    cy.setUpPublication('Book Chapter')

    cy.visitPublication()

    cy.updateFields(
      'Publication details',
      () => {
        cy.setFieldByLabel('Book title', `The book title: ${randomText}`)
      },
      true
    )

    cy.visit('/publication')

    cy.search(randomText).should('eq', 1)
  })

  function getRandomText() {
    return crypto.randomUUID().replace(/-/g, '').toUpperCase()
  }
})
