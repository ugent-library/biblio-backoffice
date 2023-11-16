describe('Issue #1237: Accessibility and mark-up: make sure labels are clickable in form elements', () => {
  beforeEach(() => {
    cy.loginAsResearcher()
  })

  describe('Add/edit publication forms', () => {
    it('should have clickable labels in the Journal Article form', () => {
      setUpPublication('Journal Article')

      updateFields('Publication details', () => {
        testFocusForLabel('Publication type', 'select[name="type"]')
        testFocusForLabel('Article type', 'select[name="journal_article_type"]')
        testFocusForLabel('DOI', 'input[type=text][name="doi"]')

        testFocusForLabel('Title', 'input[type=text][name="title"]')
        testFocusForLabel('Alternative title', 'input[type=text][name="alternative_title"]')
        testFocusForLabel('Journal title', 'input[type=text][name="publication"]')
        testFocusForLabel('Short journal title', 'input[type=text][name="publication_abbreviation"]')

        testFocusForLabel('Languages', 'select[name="language"]')
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
        testFocusForLabel('ISSN', 'input[type=text][name="issn"]')
        testFocusForLabel('E-ISSN', 'input[type=text][name="eissn"]')
        testFocusForLabel('ISBN', 'input[type=text][name="isbn"]')
        testFocusForLabel('E-ISBN', 'input[type=text][name="eisbn"]')
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

    it('should have clickable labels in the Book Chapter form', () => {
      setUpPublication('Book Chapter')

      updateFields('Publication details', () => {
        testFocusForLabel('Publication type', 'select[name="type"]')
        testFocusForLabel('DOI', 'input[type=text][name="doi"]')

        testFocusForLabel('Title', 'input[type=text][name="title"]')
        testFocusForLabel('Alternative title', 'input[type=text][name="alternative_title"]')
        testFocusForLabel('Book title', 'input[type=text][name="publication"]')

        testFocusForLabel('Languages', 'select[name="language"]')
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
        testFocusForLabel('ISSN', 'input[type=text][name="issn"]')
        testFocusForLabel('E-ISSN', 'input[type=text][name="eissn"]')
        testFocusForLabel('ISBN', 'input[type=text][name="isbn"]')
        testFocusForLabel('E-ISBN', 'input[type=text][name="eisbn"]')
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

    it('should have clickable labels in the Book form', () => {
      setUpPublication('Book')
      updateFields('Publication details', () => {
        testFocusForLabel('Publication type', 'select[name="type"]')
        testFocusForLabel('DOI', 'input[type=text][name="doi"]')

        testFocusForLabel('Title', 'input[type=text][name="title"]')
        testFocusForLabel('Alternative title', 'input[type=text][name="alternative_title"]')

        testFocusForLabel('Languages', 'select[name="language"]')
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
        testFocusForLabel('ISSN', 'input[type=text][name="issn"]')
        testFocusForLabel('E-ISSN', 'input[type=text][name="eissn"]')
        testFocusForLabel('ISBN', 'input[type=text][name="isbn"]')
        testFocusForLabel('E-ISBN', 'input[type=text][name="eisbn"]')
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

    it('should have clickable labels in the Conference contribution form', () => {
      setUpPublication('Conference contribution')

      updateFields('Publication details', () => {
        testFocusForLabel('Publication type', 'select[name="type"]')
        testFocusForLabel('Conference type', 'select[name="conference_type"]')
        testFocusForLabel('DOI', 'input[type=text][name="doi"]')

        testFocusForLabel('Title', 'input[type=text][name="title"]')
        testFocusForLabel('Alternative title', 'input[type=text][name="alternative_title"]')
        testFocusForLabel('Proceedings title', 'input[type=text][name="publication"]')
        testFocusForLabel('Publication short title', 'input[type=text][name="publication_abbreviation"]')

        testFocusForLabel('Languages', 'select[name="language"]')
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
        testFocusForLabel('Issue', 'input[type=text][name="issue"]')
        testFocusForLabel('Special issue title', 'input[type=text][name="issue_title"]')
        testFocusForLabel('First page', 'input[type=text][name="page_first"]')
        testFocusForLabel('Last page', 'input[type=text][name="page_last"]')
        testFocusForLabel('Number of pages', 'input[type=text][name="page_count"]')
        testFocusForLabel('Article number', 'input[type=text][name="article_number"]')

        testFocusForLabel('Web of Science ID', 'input[type=text][name="wos_id"]')
        testFocusForLabel('ISSN', 'input[type=text][name="issn"]')
        testFocusForLabel('E-ISSN', 'input[type=text][name="eissn"]')
        testFocusForLabel('ISBN', 'input[type=text][name="isbn"]')
        testFocusForLabel('E-ISBN', 'input[type=text][name="eisbn"]')
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

    it('should have clickable labels in the Dissertation form', () => {
      setUpPublication('Dissertation')

      updateFields('Publication details', () => {
        testFocusForLabel('Publication type', 'select[name="type"]')
        testFocusForLabel('DOI', 'input[type=text][name="doi"]')

        testFocusForLabel('Title', 'input[type=text][name="title"]')
        testFocusForLabel('Alternative title', 'input[type=text][name="alternative_title"]')

        testFocusForLabel('Languages', 'select[name="language"]')
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
        testFocusForLabel('Number of pages', 'input[type=text][name="page_count"]')

        testFocusForLabel('Date of defense', 'input[type=text][name="defense_date"]')
        testFocusForLabel('Place of defense', 'input[type=text][name="defense_place"]')

        testFocusForLabel('Web of Science ID', 'input[type=text][name="wos_id"]')
        testFocusForLabel('ISSN', 'input[type=text][name="issn"]')
        testFocusForLabel('E-ISSN', 'input[type=text][name="eissn"]')
        testFocusForLabel('ISBN', 'input[type=text][name="isbn"]')
        testFocusForLabel('E-ISBN', 'input[type=text][name="eisbn"]')
      })

      testAbstractSection()

      testLinkSection()

      testLaySummarySection()

      testAdditionalInformationSection()

      cy.contains('.nav-tabs .nav-item', 'People & Affiliations').click()

      testAuthorSection()

      testSupervisorSection()

      cy.contains('.nav-tabs .nav-item', 'Biblio Messages').click()

      testMessagesSection()
    })

    it('should have clickable labels in the Miscellaneous form', () => {
      setUpPublication('Miscellaneous')

      updateFields('Publication details', () => {
        testFocusForLabel('Publication type', 'select[name="type"]')
        testFocusForLabel('Miscellaneous type', 'select[name="miscellaneous_type"]')
        testFocusForLabel('DOI', 'input[type=text][name="doi"]')

        testFocusForLabel('Title', 'input[type=text][name="title"]')
        testFocusForLabel('Alternative title', 'input[type=text][name="alternative_title"]')
        testFocusForLabel('Publication title', 'input[type=text][name="publication"]')
        testFocusForLabel('Publication short title', 'input[type=text][name="publication_abbreviation"]')

        testFocusForLabel('Languages', 'select[name="language"]')
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
        testFocusForLabel('Issue', 'input[type=text][name="issue"]')
        testFocusForLabel('Special issue title', 'input[type=text][name="issue_title"]')
        testFocusForLabel('Edition', 'input[type=text][name="edition"]')
        testFocusForLabel('First page', 'input[type=text][name="page_first"]')
        testFocusForLabel('Last page', 'input[type=text][name="page_last"]')
        testFocusForLabel('Number of pages', 'input[type=text][name="page_count"]')
        testFocusForLabel('Article number', 'input[type=text][name="article_number"]')
        testFocusForLabel('Report number', 'input[type=text][name="report_number"]')

        testFocusForLabel('Web of Science ID', 'input[type=text][name="wos_id"]')
        testFocusForLabel('ISSN', 'input[type=text][name="issn"]')
        testFocusForLabel('E-ISSN', 'input[type=text][name="eissn"]')
        testFocusForLabel('ISBN', 'input[type=text][name="isbn"]')
        testFocusForLabel('PubMed ID', 'input[type=text][name="pubmed_id"]')
        testFocusForLabel('Arxiv ID', 'input[type=text][name="arxiv_id"]')
        testFocusForLabel('ESCI ID', 'input[type=text][name="esci_id"]')
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

    it('should have clickable labels in the Issue editor form', () => {
      setUpPublication('Issue')

      updateFields('Publication details', () => {
        testFocusForLabel('Publication type', 'select[name="type"]')
        testFocusForLabel('DOI', 'input[type=text][name="doi"]')

        testFocusForLabel('Title', 'input[type=text][name="title"]')
        testFocusForLabel('Alternative title', 'input[type=text][name="alternative_title"]')
        testFocusForLabel('Journal title', 'input[type=text][name="publication"]')
        testFocusForLabel('Short journal title', 'input[type=text][name="publication_abbreviation"]')

        testFocusForLabel('Languages', 'select[name="language"]')
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
        testFocusForLabel('Edition', 'input[type=text][name="edition"]')
        testFocusForLabel('Number of pages', 'input[type=text][name="page_count"]')

        testFocusForLabel('Web of Science ID', 'input[type=text][name="wos_id"]')
        testFocusForLabel('ISSN', 'input[type=text][name="issn"]')
        testFocusForLabel('E-ISSN', 'input[type=text][name="eissn"]')
        testFocusForLabel('ISBN', 'input[type=text][name="isbn"]')
      })

      testConferenceDetailsSection()

      testAbstractSection()

      testLinkSection()

      testAdditionalInformationSection()

      cy.contains('.nav-tabs .nav-item', 'People & Affiliations').click()

      testEditorSection()

      cy.contains('.nav-tabs .nav-item', 'Biblio Messages').click()

      testMessagesSection()
    })

    it('should have clickable labels in the file upload form', () => {
      setUpPublication('Miscellaneous')

      cy.contains('.nav-tabs .nav-item', 'Full text & Files').click()

      cy.get('input[type=file][name=file]').selectFile('cypress/fixtures/empty-pdf.pdf')

      cy.ensureModal('Document details for file empty-pdf.pdf').within(() => {
        cy.intercept('/publication/*/files/*/refresh-form*').as('refreshForm')

        cy.getLabel(/Embargoed access/).click()

        cy.wait('@refreshForm')

        testFocusForLabel('Document type', 'select[name="relation"]')
        testFocusForLabel('Publication version', 'select[name="publication_version"]')
        testFocusForLabel('Access level during embargo', 'select[name="access_level_during_embargo"]')
        testFocusForLabel('Access level after embargo', 'select[name="access_level_after_embargo"]')
        testFocusForLabel('Embargo end', 'input[type=date][name="embargo_date"]')
        testFocusForLabel('License granted by the rights holder', 'select[name="license"]')
      })
    })
  })

  describe('Add/edit dataset form', () => {
    it('should have clickable labels in the dataset form', () => {
      setUpDataset()

      updateFields('Dataset details', () => {
        cy.intercept('PUT', '/dataset/*/details/edit/refresh-form*').as('refreshForm')
        cy.setFieldByLabel('License', 'The license is not listed here')

        cy.wait('@refreshForm')

        testFocusForLabel('Title', 'input[type=text][name="title"]')
        testFocusForLabel('Persistent identifier type', 'select[name="identifier_type"]')
        testFocusForLabel('Identifier', 'input[type=text][name="identifier"]')

        testFocusForLabel('Languages', 'select[name="language"]')
        testFocusForLabel('Publication year', 'input[type=text][name="year"]')
        testFocusForLabel('Publisher', 'input[type=text][name="publisher"]')

        testFocusForLabel('Data format', 'input[type=text][name="format"]')
        // Keywords field: tagify component doesn't support focussing by label

        testFocusForLabel('License', 'select[name="license"]')
        testFocusForLabel('Other license', 'input[type=text][name="other_license"]')

        testFocusForLabel('Access level', 'select[name="access_level"]')
      })

      testAbstractSection()

      testLinkSection()

      cy.contains('.nav-tabs .nav-item', 'People & Affiliations').click()

      testCreatorSection()

      cy.contains('.nav-tabs .nav-item', 'Biblio Messages').click()

      testMessagesSection()
    })
  })

  function testFocusForLabel(labelText: string, fieldSelector: string, autoFocus = false) {
    cy.getLabel(labelText)
      .as('theLabel')
      .should('have.length', 1)
      .parent({ log: false })
      .find(fieldSelector)
      .first()
      .as('theField')
      .should('have.length', 1)
      .should(autoFocus ? 'have.focus' : 'not.have.focus')

    if (autoFocus) {
      cy.focused().blur()
    }

    cy.get('@theLabel').click()

    cy.get('@theField').should('have.focus')
  }

  function testConferenceDetailsSection() {
    updateFields('Conference details', () => {
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
  function testLaySummarySection() {
    updateFields('Lay summary', () => {
      testFocusForLabel('Lay summary', 'textarea[name="text"]')
      testFocusForLabel('Language', 'select[name="lang"]')
    })
  }

  function testAdditionalInformationSection() {
    updateFields('Additional information', () => {
      testFocusForLabel('Research field', 'select[name="research_field"]')
      // Keywords field: tagify component doesn't support focussing by label
      testFocusForLabel('Additional information', 'textarea[name="additional_info"]')
    })
  }

  function testAuthorSection() {
    updateFields('Author', () => {
      testFocusForLabel('First name', 'input[name="first_name"]', true)
      testFocusForLabel('Last name', 'input[name="last_name"]')

      cy.setFieldByLabel('First name', 'Griet')
      cy.setFieldByLabel('Last name', 'Alleman')

      cy.contains('.btn', 'Add author').click()

      testFocusForLabel('Roles', 'select[name="credit_role"]')
    })
  }

  function testEditorSection() {
    updateFields('Editor', () => {
      testFocusForLabel('First name', 'input[name="first_name"]', true)
      testFocusForLabel('Last name', 'input[name="last_name"]')
    })
  }

  function testCreatorSection() {
    updateFields('Creator', () => {
      testFocusForLabel('First name', 'input[name="first_name"]', true)
      testFocusForLabel('Last name', 'input[name="last_name"]')
    })
  }

  function testSupervisorSection() {
    updateFields('Supervisor', () => {
      testFocusForLabel('First name', 'input[name="first_name"]', true)
      testFocusForLabel('Last name', 'input[name="last_name"]')
    })
  }

  function testMessagesSection() {
    updateFields('Messages from and for Biblio team', () => {
      testFocusForLabel('Message', 'textarea[name="message"]')
    })
  }

  type PublicationType =
    | 'Journal Article'
    | 'Book Chapter'
    | 'Book'
    | 'Conference contribution'
    | 'Dissertation'
    | 'Miscellaneous'
    | 'Issue'

  function setUpPublication(publicationType: PublicationType) {
    cy.visit('/publication/add')

    cy.contains('Enter a publication manually').find(':radio').click()
    cy.contains('.btn', 'Add publication(s)').click()

    cy.contains(new RegExp(`^${publicationType}$`)).click()
    cy.contains('.btn', 'Add publication(s)').click()

    updateFields(
      'Publication details',
      () => {
        cy.setFieldByLabel('Title', `The ${publicationType} title [CYPRESSTEST]`)
      },
      true
    )
  }

  function setUpDataset() {
    cy.visit('/dataset/add')

    cy.contains('Register a dataset manually').find(':radio').click()
    cy.contains('.btn', 'Add dataset').click()

    updateFields(
      'Dataset details',
      () => {
        cy.setFieldByLabel('Title', `The dataset title [CYPRESSTEST]`)

        cy.setFieldByLabel('Persistent identifier type', 'DOI')

        cy.setFieldByLabel('Identifier', '10.5072/test/t')
      },
      true
    )
  }

  type FieldsSection =
    | 'Publication details'
    | 'Conference details'
    | 'Abstract'
    | 'Link'
    | 'Lay summary'
    | 'Additional information'
    | 'Author'
    | 'Editor'
    | 'Supervisor'
    | 'Messages from and for Biblio team'
    | 'Dataset details'
    | 'Creator'

  function updateFields(section: FieldsSection, callback: () => void, persist = false) {
    cy.contains('.card-header', section).find('.btn').click()

    const modalTitle = new RegExp(`(Edit|Add) ${section}`, 'i')

    cy.ensureModal(modalTitle).within(callback).closeModal(persist)

    cy.ensureNoModal()
  }
})
