describe('Issue #1127: Cannot search any longer on book title, journal title, short journal title nor conference title', () => {
  beforeEach(() => {
    cy.loginAsResearcher()
  })

  it('should be possible to search by publisher', () => {
    const randomText = getRandomText()

    setUpPublication('Miscellaneous', () => {
      updateDescriptionFields('Publication details', () => {
        setField('Publisher', `Alternative title: ${randomText}`)
      })
    })

    cy.search(randomText).should('eq', 1)
  })

  it('should be possible to search by alternative title', () => {
    const randomText1 = getRandomText()
    const randomText2 = getRandomText()

    setUpPublication('Miscellaneous', () => {
      updateDescriptionFields('Publication details', () => {
        setField('Alternative title', `Alternative title: ${randomText1}`)
          .closest('.form-value')
          .contains('.btn', 'Add')
          .click()
          .closest('.form-value')
          .next('.form-value')
          .find('input')
          .type(`Alternative title: ${randomText2}`)
      })
    })

    cy.search(randomText1).should('eq', 1)
    cy.search(randomText2).should('eq', 1)
  })

  it('should be possible to search by conference title', () => {
    const randomText = getRandomText()

    setUpPublication('Conference contribution', () => {
      updateDescriptionFields('Conference Details', () => {
        setField('Conference', `The conference name: ${randomText}`)
      })
    })

    cy.search(randomText).should('eq', 1)
  })

  it('should be possible to search by journal title', () => {
    const randomText = getRandomText()

    setUpPublication('Journal Article', () => {
      updateDescriptionFields('Publication details', () => {
        setField('Journal title', `The journal name: ${randomText}`)
      })
    })

    cy.search(randomText).should('eq', 1)
  })

  it('should be possible to search by short journal title', () => {
    const randomText = getRandomText()

    setUpPublication('Journal Article', () => {
      updateDescriptionFields('Publication details', () => {
        setField('Short journal title', `The short journal name: ${randomText}`)
      })
    })

    cy.search(randomText).should('eq', 1)
  })

  it('should be possible to search by book title', () => {
    const randomText = getRandomText()

    setUpPublication('Book Chapter', () => {
      updateDescriptionFields('Publication details', () => {
        setField('Book', `The book title: ${randomText}`)
      })
    })

    cy.search(randomText).should('eq', 1)
  })

  type PublicationType = 'Journal Article' | 'Book Chapter' | 'Conference contribution' | 'Miscellaneous'

  function setUpPublication(publicationType: PublicationType, editPublicationCallback: () => void) {
    cy.visit('/publication/add')

    cy.contains('Enter a publication manually').click()
    cy.contains('.btn', 'Add publication(s)').click()

    cy.contains(publicationType).click()
    cy.contains('.btn', 'Add publication(s)').click()

    updateDescriptionFields('Publication details', () => {
      setField('Title', `The primary ${publicationType} title [CYPRESSTEST]`)
      setField('Publication year', new Date().getFullYear().toString())
    })

    // Custom callback here
    editPublicationCallback()

    cy.ensureNoModal()

    cy.contains('.btn', 'Complete Description').click()

    cy.contains('.btn', 'Save as draft').click()
  }

  function updateDescriptionFields(section: 'Publication details' | 'Conference Details', callback: () => void) {
    cy.contains('.card-header', section).contains('.btn', 'Edit').click({ scrollBehavior: 'nearest' })

    cy.ensureModal(`Edit ${section.toLowerCase()}`).within(() => {
      callback()

      cy.contains('.modal-footer .btn', 'Save').click()
    })

    cy.ensureNoModal()
  }

  function setField(fieldLabel: string, value: string) {
    return cy.contains('label', fieldLabel).next().find('input').type(value)
  }

  function getRandomText() {
    return crypto.randomUUID().replace(/-/g, '').toUpperCase()
  }
})
