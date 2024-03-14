describe("The home page", () => {
  it("should be able to load the home page anonymously", () => {
    cy.visit("/");

    cy.get("h1").should("have.text", "Biblio Back Office");

    cy.contains("a", "Log in").should("be.visible");

    cy.contains("h2", "Publications").should("be.visible");

    cy.contains("h2", "Datasets").should("be.visible");

    // On a regular page you don't have to do this to make the "be.visible" assertion work,
    // but in this case the elements are being clipped by an element with "overflow: scroll"
    cy.get(".u-scroll-wrapper__body").scrollTo("bottom", { duration: 250 });

    cy.contains("h2", "Biblio Academic Bibliography").should("be.visible");

    cy.contains("h2", "Help").should("be.visible");

    cy.get(".c-sidebar").should("not.have.class", "c-sidebar--dark-gray");
  });

  it("should redirect to the login page when clicking the Login buttons", () => {
    cy.visit("/");

    const assertLoginRedirection = (href) => {
      cy.request(href).then((response) => {
        expect(response).to.have.property("isOkStatusCode", true);
        expect(response)
          .to.have.property("redirects")
          .that.is.an("array")
          .that.has.length(1);

        const redirects = response.redirects
          .map((url) => url.replace(/^3\d\d\: /, "")) // Redirect entries are in form '3XX: {url}'
          .map((url) => new URL(url));

        expect(redirects[0]).to.have.property(
          "origin",
          Cypress.env("OIDC_ORIGIN"),
        );
      });
    };

    cy.get('header .btn:contains("Log in"), main .btn:contains("Log in")')
      .should("have.length", 2)
      .map("href")
      .unique() // No need to check the same URL more than once
      .each(assertLoginRedirection);
  });

  it("should be able to logon as researcher", () => {
    cy.loginAsResearcher();

    cy.visit("/");

    cy.get(".nav-main .dropdown-menu")
      .as("user-menu")
      .should("have.css", "display", "none");
    cy.get(".nav-main button.dropdown-toggle").click();
    cy.get("@user-menu").should("have.css", "display", "block");

    cy.get(".nav-main .dropdown-menu .dropdown-item").should("have.length", 1);
    cy.contains(".dropdown-menu .dropdown-item", "View as").should("not.exist");
    cy.contains(".dropdown-menu .dropdown-item", "Logout").should("exist");

    cy.get(".c-sidebar button.dropdown-toggle").should("not.exist");
    cy.get(".c-sidebar-menu .c-sidebar__item").should("have.length", 3);
    cy.contains(".c-sidebar__item", "Dashboard").should("be.visible");
    cy.contains(".c-sidebar__item", "Biblio Publications").should("be.visible");
    cy.contains(".c-sidebar__item", "Biblio Datasets").should("be.visible");
    cy.contains(".c-sidebar__item", "Batch").should("not.exist");

    cy.get(".c-sidebar").should("not.have.class", "c-sidebar--dark-gray");
  });

  it("should be able to logon as librarian and switch to librarian mode", () => {
    cy.loginAsLibrarian();

    cy.visit("/");

    cy.get(".nav-main .dropdown-menu .dropdown-item").should("have.length", 2);
    cy.contains(".dropdown-menu .dropdown-item", "View as").should("exist");
    cy.contains(".dropdown-menu .dropdown-item", "Logout").should("exist");

    cy.get(".c-sidebar button.dropdown-toggle")
      .find(".if-briefcase")
      .should("be.visible");
    cy.get(".c-sidebar button.dropdown-toggle")
      .find(".if-book")
      .should("not.exist");
    cy.get(".c-sidebar button.dropdown-toggle").should(
      "contain.text",
      "Researcher",
    );
    cy.get(".c-sidebar-menu .c-sidebar__item").should("have.length", 3);
    cy.contains(".c-sidebar__item", "Dashboard").should("be.visible");
    cy.contains(".c-sidebar__item", "Biblio Publications").should("be.visible");
    cy.contains(".c-sidebar__item", "Biblio Datasets").should("be.visible");
    cy.contains(".c-sidebar__item", "Batch").should("not.exist");

    cy.get(".c-sidebar").should("not.have.class", "c-sidebar--dark-gray");

    cy.switchMode("Librarian");

    cy.get(".c-sidebar button.dropdown-toggle")
      .find(".if-book")
      .should("be.visible");
    cy.get(".c-sidebar button.dropdown-toggle")
      .find(".if-briefcase")
      .should("not.exist");
    cy.get(".c-sidebar button.dropdown-toggle").should(
      "contain.text",
      "Librarian",
    );
    cy.get(".c-sidebar-menu .c-sidebar__item").should("have.length", 4);
    cy.contains(".c-sidebar__item", "Dashboard").should("be.visible");
    cy.contains(".c-sidebar__item", "Biblio Publications").should("be.visible");
    cy.contains(".c-sidebar__item", "Biblio Datasets").should("be.visible");
    cy.contains(".c-sidebar__item", "Batch").should("be.visible");

    cy.get(".c-sidebar").should("have.class", "c-sidebar--dark-gray");

    cy.switchMode("Researcher");

    cy.get(".c-sidebar button.dropdown-toggle")
      .find(".if-briefcase")
      .should("be.visible");
    cy.get(".c-sidebar button.dropdown-toggle")
      .find(".if-book")
      .should("not.exist");
    cy.get(".c-sidebar button.dropdown-toggle").should(
      "contain.text",
      "Researcher",
    );
    cy.get(".c-sidebar-menu .c-sidebar__item").should("have.length", 3);
    cy.contains(".c-sidebar__item", "Dashboard").should("be.visible");
    cy.contains(".c-sidebar__item", "Biblio Publications").should("be.visible");
    cy.contains(".c-sidebar__item", "Biblio Datasets").should("be.visible");
    cy.contains(".c-sidebar__item", "Batch").should("not.exist");

    cy.get(".c-sidebar").should("not.have.class", "c-sidebar--dark-gray");
  });

  it("should not set the biblio-backoffice cookie twice when switching roles", () => {
    cy.loginAsLibrarian();

    cy.visit("/");

    cy.intercept({ method: "PUT", pathname: "/role/curator" }).as(
      "role-curator",
    );
    cy.intercept({ method: "PUT", pathname: "/role/user" }).as("role-user");

    cy.switchMode("Librarian");

    cy.wait("@role-curator")
      .its("response.headers[set-cookie]")
      .then((cookies) => {
        expect(
          cookies.filter((c) => c.startsWith("biblio-backoffice=")),
        ).to.have.length(1);
      });

    cy.switchMode("Researcher");

    cy.wait("@role-user")
      .its("response.headers[set-cookie]")
      .then((cookies) => {
        expect(
          cookies.filter((c) => c.startsWith("biblio-backoffice=")),
        ).to.have.length(1);
      });
  });
});
