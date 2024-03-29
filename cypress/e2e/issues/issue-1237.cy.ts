describe('Issue #1237: Accessibility and mark-up: make sure labels are clickable in form elements', () => {
  describe('In researcher mode', () => {
    beforeEach(() => {
      cy.loginAsResearcher()
    })

    describe('Add/edit publication forms', () => {
      it('should have clickable labels in the Journal Article form', () => {
        cy.setUpPublication('Journal Article')

        cy.visitPublication()

        cy.updateFields('Publication details', () => {
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
        cy.setUpPublication('Book Chapter')

        cy.visitPublication()

        cy.updateFields('Publication details', () => {
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
        cy.setUpPublication('Book')

        cy.visitPublication()

        cy.updateFields('Publication details', () => {
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
        cy.setUpPublication('Conference contribution')

        cy.visitPublication()

        cy.updateFields('Publication details', () => {
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
        cy.setUpPublication('Dissertation')

        cy.visitPublication()

        cy.updateFields('Publication details', () => {
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
        cy.setUpPublication('Miscellaneous')

        cy.visitPublication()

        cy.updateFields('Publication details', () => {
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
        cy.setUpPublication('Issue')

        cy.visitPublication()

        cy.updateFields('Publication details', () => {
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
        cy.setUpPublication('Miscellaneous')

        cy.visitPublication()

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
        cy.setUpDataset()

        cy.visitDataset()

        cy.updateFields('Dataset details', () => {
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
  })

  describe('In librarian mode', () => {
    beforeEach(() => {
      cy.loginAsLibrarian()

      cy.switchMode('Librarian')
    })

    describe('Add/edit publication forms', () => {
      it('should have clickable labels in librarian only sections', () => {
        cy.setUpPublication('Book')

        cy.visitPublication()

        testLibrarianTagsSection()

        testLibrarianNoteSection()
      })
    })

    describe('Add/edit dataset form', () => {
      it('should have clickable labels in librarian only sections', () => {
        cy.setUpDataset()

        cy.visitDataset()

        testLibrarianTagsSection()

        testLibrarianNoteSection()
      })
    })
  })

  function testFocusForLabel(labelText: string, fieldSelector: string, autoFocus = false) {
    cy.getLabel(labelText)
      .as('theLabel', { type: 'static' })
      .should('have.length', 1)
      .parent({ log: false })
      .find(fieldSelector)
      .should('have.length', 1)
      .first({ log: false })
      .as('theField', { type: 'static' })
      .should(autoFocus ? 'be.focused' : 'not.be.focused')

    if (autoFocus) {
      cy.focused().blur()

      cy.get('@theField').should('not.be.focused')
    }

    cy.get('@theLabel').click()

    cy.get('@theField').should('be.focused')
  }

  function testConferenceDetailsSection() {
    cy.updateFields('Conference details', () => {
      testFocusForLabel('Conference', 'input[type=text][name="name"]')
      testFocusForLabel('Conference location', 'input[type=text][name="location"]')
      testFocusForLabel('Conference organiser', 'input[type=text][name="organizer"]')
      testFocusForLabel('Conference start date', 'input[type=text][name="start_date"]')
      testFocusForLabel('Conference end date', 'input[type=text][name="end_date"]')
    })
  }

  function testAbstractSection() {
    cy.updateFields('Abstract', () => {
      testFocusForLabel('Abstract', 'textarea[name="text"]')
      testFocusForLabel('Language', 'select[name="lang"]')
    })
  }

  function testLinkSection() {
    cy.updateFields('Link', () => {
      testFocusForLabel('URL', 'input[type=text][name="url"]')
      testFocusForLabel('Relation', 'select[name="relation"]')
      testFocusForLabel('Description', 'input[type=text][name="description"]')
    })
  }
  function testLaySummarySection() {
    cy.updateFields('Lay summary', () => {
      testFocusForLabel('Lay summary', 'textarea[name="text"]')
      testFocusForLabel('Language', 'select[name="lang"]')
    })
  }

  function testAdditionalInformationSection() {
    cy.updateFields('Additional information', () => {
      testFocusForLabel('Research field', 'select[name="research_field"]')
      // Keywords field: tagify component doesn't support focussing by label
      testFocusForLabel('Additional information', 'textarea[name="additional_info"]')
    })
  }

  function testAuthorSection() {
    cy.updateFields('Authors', () => {
      testFocusForLabel('First name', 'input[name="first_name"]', true)
      testFocusForLabel('Last name', 'input[name="last_name"]')

      cy.setFieldByLabel('First name', 'Griet')
      cy.setFieldByLabel('Last name', 'Alleman')

      cy.contains('.btn', 'Add author').click()

      testFocusForLabel('Roles', 'select[name="credit_role"]')
    })
  }

  function testEditorSection() {
    cy.updateFields('Editors', () => {
      testFocusForLabel('First name', 'input[name="first_name"]', true)
      testFocusForLabel('Last name', 'input[name="last_name"]')
    })
  }

  function testCreatorSection() {
    cy.updateFields('Creators', () => {
      testFocusForLabel('First name', 'input[name="first_name"]', true)
      testFocusForLabel('Last name', 'input[name="last_name"]')
    })
  }

  function testSupervisorSection() {
    cy.updateFields('Supervisors', () => {
      testFocusForLabel('First name', 'input[name="first_name"]', true)
      testFocusForLabel('Last name', 'input[name="last_name"]')
    })
  }

  function testLibrarianTagsSection() {
    cy.updateFields('Librarian tags', () => {
      testFocusForLabel('Librarian tags', 'input[name="reviewer_tags"]')
    })
  }

  function testLibrarianNoteSection() {
    cy.updateFields('Librarian note', () => {
      testFocusForLabel('Librarian note', 'textarea[name="reviewer_note"]')
    })
  }

  function testMessagesSection() {
    cy.updateFields('Messages from and for Biblio team', () => {
      testFocusForLabel('Message', 'textarea[name="message"]')
    })
  }
})
