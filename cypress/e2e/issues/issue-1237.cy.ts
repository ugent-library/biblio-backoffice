describe('Issue #1237: Accessibility and mark-up: make sure labels are clickable in form elements', () => {
  beforeEach(() => {
    cy.loginAsResearcher()
  })

  it('should have clickable labels in the Journal Article form', () => {
    setUpPublication('Journal Article', () => {
      updateFields('Publication details', () => {
        testFocusForLabel('Publication type', 'select[name="type"]')
        testFocusForLabel('Article type', 'select[name="journal_article_type"]')
        testFocusForLabel('DOI', 'input[type=text][name="doi"]')

        testFocusForLabel('Title', 'input[type=text][name="title"]')
        testFocusForLabel('Alternative title', 'input[type=text][name="alternative_title[0]"]')
        testFocusForLabel('Journal title', 'input[type=text][name="publication"]')
        testFocusForLabel('Short journal title', 'input[type=text][name="publication_abbreviation"]')

        testFocusForLabel('Languages', 'select[name="language[0]"]')
        testFocusForLabel('Publishing status', 'select[name="publication_status"]')
        testFocusForLabel(
          'Published while none of the authors and editors were employed at UGent',
          ':checkbox[name="extern"]'
        )
        testFocusForLabel('Publication year', 'input[type=text][name="year"]')
        testFocusForLabel('Place of publication', 'input[type=text][name="place_of_publication"]')
        testFocusForLabel('Publisher', 'input[type=text][name="publisher"]')

        testFocusForLabel('Volume', 'input[type=text][name="volume"]')
        testFocusForLabel('Issue', 'input[type=text][name="issue"]')
        testFocusForLabel('Special issue title', 'input[type=text][name="issue_title"]')
        testFocusForLabel('First page', 'input[type=text][name="page_first"]')
        testFocusForLabel('Last page', 'input[type=text][name="page_last"]')
        testFocusForLabel('Number of pages', 'input[type=text][name="page_count"]')
        testFocusForLabel('Article number', 'input[type=text][name="article_number"]')

        testFocusForLabel('Web of Science ID', 'input[type=text][name="wos_id"]')
        testFocusForLabel('ISSN', 'input[type=text][name="issn[0]"]')
        testFocusForLabel('E-ISSN', 'input[type=text][name="eissn[0]"]')
        testFocusForLabel('ISBN', 'input[type=text][name="isbn[0]"]')
        testFocusForLabel('E-ISBN', 'input[type=text][name="eisbn[0]"]')
        testFocusForLabel('PubMed ID', 'input[type=text][name="pubmed_id"]')
        testFocusForLabel('Arxiv ID', 'input[type=text][name="arxiv_id"]')
        testFocusForLabel('ESCI ID', 'input[type=text][name="esci_id"]')
      })

      testConferenceDetailsSection()

      testAbstractSection()

      testLinkSection()

      testAdditionalInformationSection()

      cy.contains('.nav-tabs .nav-item', 'People & Affiliations').click()

      testAuthorSection()

      testEditorSection()

      cy.contains('.nav-tabs .nav-item', 'Biblio Messages').click()

      testMessagesSection()
    })
  })

  it('should have clickable labels in the Book Chapter form', () => {
    setUpPublication('Book Chapter', () => {
      updateFields('Publication details', () => {
        testFocusForLabel('Publication type', 'select[name="type"]')
        testFocusForLabel('DOI', 'input[type=text][name="doi"]')

        testFocusForLabel('Title', 'input[type=text][name="title"]')
        testFocusForLabel('Alternative title', 'input[type=text][name="alternative_title[0]"]')
        testFocusForLabel('Book title', 'input[type=text][name="publication"]')

        testFocusForLabel('Languages', 'select[name="language[0]"]')
        testFocusForLabel('Publishing status', 'select[name="publication_status"]')
        testFocusForLabel(
          'Published while none of the authors and editors were employed at UGent',
          ':checkbox[name="extern"]'
        )
        testFocusForLabel('Publication year', 'input[type=text][name="year"]')
        testFocusForLabel('Place of publication', 'input[type=text][name="place_of_publication"]')
        testFocusForLabel('Publisher', 'input[type=text][name="publisher"]')

        testFocusForLabel('Series title', 'input[type=text][name="series_title"]')
        testFocusForLabel('Volume', 'input[type=text][name="volume"]')
        testFocusForLabel('Edition', 'input[type=text][name="edition"]')
        testFocusForLabel('First page', 'input[type=text][name="page_first"]')
        testFocusForLabel('Last page', 'input[type=text][name="page_last"]')
        testFocusForLabel('Number of pages', 'input[type=text][name="page_count"]')

        testFocusForLabel('Web of Science ID', 'input[type=text][name="wos_id"]')
        testFocusForLabel('ISSN', 'input[type=text][name="issn[0]"]')
        testFocusForLabel('E-ISSN', 'input[type=text][name="eissn[0]"]')
        testFocusForLabel('ISBN', 'input[type=text][name="isbn[0]"]')
        testFocusForLabel('E-ISBN', 'input[type=text][name="eisbn[0]"]')
      })

      testConferenceDetailsSection()

      testAbstractSection()

      testLinkSection()

      testAdditionalInformationSection()

      cy.contains('.nav-tabs .nav-item', 'People & Affiliations').click()

      testAuthorSection()

      testEditorSection()

      cy.contains('.nav-tabs .nav-item', 'Biblio Messages').click()

      testMessagesSection()
    })
  })

  it('should have clickable labels in the Book form', () => {
    setUpPublication('Book', () => {
      updateFields('Publication details', () => {
        testFocusForLabel('Publication type', 'select[name="type"]')
        testFocusForLabel('DOI', 'input[type=text][name="doi"]')

        testFocusForLabel('Title', 'input[type=text][name="title"]')
        testFocusForLabel('Alternative title', 'input[type=text][name="alternative_title[0]"]')

        testFocusForLabel('Languages', 'select[name="language[0]"]')
        testFocusForLabel('Publishing status', 'select[name="publication_status"]')
        testFocusForLabel(
          'Published while none of the authors and editors were employed at UGent',
          ':checkbox[name="extern"]'
        )
        testFocusForLabel('Publication year', 'input[type=text][name="year"]')
        testFocusForLabel('Place of publication', 'input[type=text][name="place_of_publication"]')
        testFocusForLabel('Publisher', 'input[type=text][name="publisher"]')

        testFocusForLabel('Series title', 'input[type=text][name="series_title"]')
        testFocusForLabel('Volume', 'input[type=text][name="volume"]')
        testFocusForLabel('Edition', 'input[type=text][name="edition"]')
        testFocusForLabel('Number of pages', 'input[type=text][name="page_count"]')

        testFocusForLabel('Web of Science ID', 'input[type=text][name="wos_id"]')
        testFocusForLabel('ISSN', 'input[type=text][name="issn[0]"]')
        testFocusForLabel('E-ISSN', 'input[type=text][name="eissn[0]"]')
        testFocusForLabel('ISBN', 'input[type=text][name="isbn[0]"]')
        testFocusForLabel('E-ISBN', 'input[type=text][name="eisbn[0]"]')
      })

      testAbstractSection()

      testLinkSection()

      testAdditionalInformationSection()

      cy.contains('.nav-tabs .nav-item', 'People & Affiliations').click()

      testAuthorSection()

      testEditorSection()

      cy.contains('.nav-tabs .nav-item', 'Biblio Messages').click()

      testMessagesSection()
    })
  })

  function testFocusForLabel(labelText: string, fieldSelector: string, autoFocus = false) {
    getLabel(labelText)
      .as('theLabel')
      .should('have.length', 1)
      .parent({ log: false })
      .find(fieldSelector)
      .as('theField')
      .should('have.length', 1)
      .should(autoFocus ? 'have.focus' : 'not.have.focus')

    // TODO: remove
    // cy.get('@theLabel').then(label => {
    //   cy.get('@theField').then(field => {
    //     if (!field.attr('id')) {
    //       field.attr('id', label.attr('for'))
    //     }
    //   })
    // })

    cy.get('@theLabel').click()

    return cy.get('@theField').should('have.focus')
  }

  function testConferenceDetailsSection() {
    updateFields('Conference Details', () => {
      testFocusForLabel('Conference', 'input[type=text][name="name"]')
      testFocusForLabel('Conference location', 'input[type=text][name="location"]')
      testFocusForLabel('Conference organiser', 'input[type=text][name="organizer"]')
      testFocusForLabel('Conference start date', 'input[type=text][name="start_date"]')
      testFocusForLabel('Conference end date', 'input[type=text][name="end_date"]')
    })
  }

  function testAbstractSection() {
    updateFields('Abstract', () => {
      testFocusForLabel('Abstract', 'textarea[name="text"]')
      testFocusForLabel('Language', 'select[name="lang"]')
    })
  }

  function testLinkSection() {
    updateFields('Link', () => {
      testFocusForLabel('URL', 'input[type=text][name="url"]')
      testFocusForLabel('Relation', 'select[name="relation"]')
      testFocusForLabel('Description', 'input[type=text][name="description"]')
    })
  }

  function testAdditionalInformationSection() {
    updateFields('Additional information', () => {
      testFocusForLabel('Research field', 'select[name="research_field[0]"]')
      // TODO: tagify component doesn't support focussing by label
      // testFocusForLabel('Keywords', 'tags > .tagify__input')
      testFocusForLabel('Additional information', 'textarea[name="additional_info"]')
    })
  }

  function testAuthorSection() {
    updateFields('Author', () => {
      testFocusForLabel('First name', 'input[name="first_name"]', true).type('Dries')
      testFocusForLabel('Last name', 'input[name="last_name"]').type('Moreels')

      cy.contains('.btn', 'Add author').click()

      testFocusForLabel('Roles', 'select[name="credit_role[0]"]')
    })
  }

  function testEditorSection() {
    updateFields('Editor', () => {
      testFocusForLabel('First name', 'input[name="first_name"]', true)
      testFocusForLabel('Last name', 'input[name="last_name"]')
    })
  }

  function testMessagesSection() {
    updateFields('Messages from and for Biblio team', () => {
      testFocusForLabel('Message', 'textarea[name="message"]')
    })
  }

  type PublicationType = 'Journal Article' | 'Book Chapter' | 'Book' | 'Conference contribution' | 'Miscellaneous'

  function setUpPublication(publicationType: PublicationType, editPublicationCallback: () => void) {
    cy.visit('/publication/add')

    cy.contains('Enter a publication manually').click()
    cy.contains('.btn', 'Add publication(s)').click()

    cy.contains(new RegExp(`^${publicationType}$`)).click()
    cy.contains('.btn', 'Add publication(s)').click()

    updateFields(
      'Publication details',
      () => {
        setField('Title', `The primary ${publicationType} title [CYPRESSTEST]`)
        setField('Publication year', new Date().getFullYear().toString())
      },
      true
    )

    // Custom callback here
    editPublicationCallback()

    cy.ensureNoModal()
  }

  type FieldsSection =
    | 'Publication details'
    | 'Conference Details'
    | 'Abstract'
    | 'Link'
    | 'Additional information'
    | 'Author'
    | 'Editor'
    | 'Messages from and for Biblio team'

  function updateFields(section: FieldsSection, callback: () => void, persist = false) {
    cy.contains('.card-header', section).find('.btn').click({ scrollBehavior: 'nearest' })

    const modalTitle = new RegExp(`(Edit|Add) ${section}`, 'i')

    cy.ensureModal(modalTitle).within(() => {
      callback()

      cy.contains('.modal-footer .btn', persist ? 'Save' : 'Cancel').click()
    })

    cy.ensureNoModal()
  }

  function setField(fieldLabel: string, value: string) {
    return getLabel(fieldLabel).next().find('input').type(value)
  }

  function getLabel(labelText: string) {
    return cy.get(`label:contains("${labelText}")`).filter(
      (_, el) => {
        // Make sure to match the exact label text (excluding badges)
        const $label = Cypress.$(el).clone()

        $label.find('.badge').remove()

        return $label.text().trim() === labelText
      },
      { log: false }
    )
  }
})
