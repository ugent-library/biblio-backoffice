describe('The home page', () => {
  it('should be able to load the home page anonymously', () => {
    cy.visit('/')

    cy.get('h2').should('have.text', 'Home')

    cy.contains('a', 'Log in').should('be.visible')

    cy.contains('h3', 'Publications').should('be.visible')
    cy.contains('a', 'Go to publications').should('be.visible')

    cy.contains('h3', 'Datasets').should('be.visible')
    cy.contains('a', 'Go to datasets').should('be.visible')

    // On a regular page you don't have to do this to make the "be.visible" assertion work,
    // but in this case the elements are being clipped by an element with "overflow: scroll"
    cy.get('.u-scroll-wrapper__body').scrollTo('bottom', { duration: 250 })

    cy.contains('h3', 'Biblio Academic Bibliography').should('be.visible')
    cy.contains('a', 'Go to Biblio Academic Bibliography').should('be.visible')

    cy.contains('h3', 'Help').should('be.visible')
    cy.contains('a', 'How to register and deposit').should('be.visible')
  })

  it('should redirect to the login page when browsing publications anonymously', () => {
    cy.visit('/')

    cy.contains('Biblio Publications').click()

    cy.url().should('contain', 'liblogin.ugent.be')
  })

  it('should redirect to the login page when browsing datasets anonymously', () => {
    cy.visit('/')

    cy.contains('Biblio Datasets').click()

    cy.url().should('contain', 'liblogin.ugent.be')
  })

  it('should be able to logon as researcher', () => {
    cy.loginAsResearcher()

    cy.visit('/')

    cy.get('.nav-main .dropdown-menu').as('user-menu').should('have.css', 'display', 'none')
    cy.get('.nav-main button.dropdown-toggle').click()
    cy.get('@user-menu').should('have.css', 'display', 'block')

    cy.get('.nav-main .dropdown-menu .dropdown-item').should('have.length', 1)
    cy.contains('.dropdown-menu .dropdown-item', 'View as').should('not.exist')
    cy.contains('.dropdown-menu .dropdown-item', 'Logout').should('exist')

    cy.get('.c-sidebar button.dropdown-toggle').should('not.exist')
    cy.get('.c-sidebar').should('not.have.class', 'c-sidebar--dark-gray')
    cy.get('.c-sidebar-menu .c-sidebar__item').should('have.length', 2)
    cy.contains('.c-sidebar__item', 'Biblio Publications').should('be.visible')
    cy.contains('.c-sidebar__item', 'Biblio Datasets').should('be.visible')
    cy.contains('.c-sidebar__item', 'Dashboard').should('not.exist')
  })

  it('should be able to logon as librarian and switch to librarian mode', () => {
    cy.loginAsLibrarian()

    cy.visit('/')

    cy.get('.nav-main .dropdown-menu .dropdown-item').should('have.length', 2)
    cy.contains('.dropdown-menu .dropdown-item', 'View as').should('exist')
    cy.contains('.dropdown-menu .dropdown-item', 'Logout').should('exist')

    cy.get('.c-sidebar button.dropdown-toggle').should('contain.text', 'Researcher')
    cy.get('.c-sidebar').should('not.have.class', 'c-sidebar--dark-gray')
    cy.get('.c-sidebar-menu .c-sidebar__item').should('have.length', 2)
    cy.contains('.c-sidebar__item', 'Biblio Publications').should('be.visible')
    cy.contains('.c-sidebar__item', 'Biblio Datasets').should('be.visible')
    cy.contains('.c-sidebar__item', 'Dashboard').should('not.exist')

    cy.switchMode('Librarian')

    cy.get('.c-sidebar button.dropdown-toggle').should('contain.text', 'Librarian')
    cy.get('.c-sidebar').should('have.class', 'c-sidebar--dark-gray')
    cy.get('.c-sidebar-menu .c-sidebar__item').should('have.length', 4)
    cy.contains('.c-sidebar__item', 'Biblio Publications').should('be.visible')
    cy.contains('.c-sidebar__item', 'Biblio Datasets').should('be.visible')
    cy.contains('.c-sidebar__item', 'Dashboard').should('be.visible')
    cy.contains('.c-sidebar__item', 'Batch').should('be.visible')

    cy.switchMode('Researcher')

    cy.get('.c-sidebar button.dropdown-toggle').should('contain.text', 'Researcher')
    cy.get('.c-sidebar').should('not.have.class', 'c-sidebar--dark-gray')
    cy.get('.c-sidebar-menu .c-sidebar__item').should('have.length', 2)
    cy.contains('.c-sidebar__item', 'Biblio Publications').should('be.visible')
    cy.contains('.c-sidebar__item', 'Biblio Datasets').should('be.visible')
    cy.contains('.c-sidebar__item', 'Dashboard').should('not.exist')
  })

  it('should not set the biblio-backoffice cookie twice when switching roles', () => {
    cy.loginAsLibrarian()

    cy.visit('/')

    cy.intercept({ method: 'PUT', pathname: '/role/curator' }).as('role-curator')
    cy.intercept({ method: 'PUT', pathname: '/role/user' }).as('role-user')

    cy.switchMode('Librarian')

    cy.wait('@role-curator')
      .its('response.headers[set-cookie]')
      .then(cookies => {
        expect(cookies.filter(c => c.startsWith('biblio-backoffice='))).to.have.length(1)
      })

    cy.switchMode('Researcher')

    cy.wait('@role-user')
      .its('response.headers[set-cookie]')
      .then(cookies => {
        expect(cookies.filter(c => c.startsWith('biblio-backoffice='))).to.have.length(1)
      })
  })
})
