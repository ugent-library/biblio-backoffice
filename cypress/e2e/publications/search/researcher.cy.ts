import { getRandomText } from "support/util";

describe("The publication search (for researchers)", () => {
  it("should display filtered navigation tabs with all publications, my publications, supervised by me and registered by me", () => {
    const randomTitleSuffix = getRandomText();

    // Setup
    cy.loginAsLibrarian();
    cy.setUpPublication("Dissertation", {
      biblioIDAlias: "dissertation1",
      title: `Dissertation 1 ${randomTitleSuffix}`,
      shouldWaitForIndex: true,
    });
    cy.addAuthor("Biblio", "Librarian", { biblioIdAlias: "@dissertation1" });
    cy.addSupervisor("Biblio", "Researcher", {
      biblioIdAlias: "@dissertation1",
    });

    cy.loginAsResearcher();
    cy.setUpPublication("Book", {
      biblioIDAlias: "book",
      title: `Book ${randomTitleSuffix}`,
      shouldWaitForIndex: true,
    });
    cy.addAuthor("John", "Doe", { biblioIdAlias: "@book" });
    cy.addEditor("Biblio", "Researcher", { biblioIdAlias: "@book" });

    cy.setUpPublication("Dissertation", {
      biblioIDAlias: "dissertation2",
      title: `Dissertation 2 ${randomTitleSuffix}`,
      shouldWaitForIndex: true,
    });
    cy.addAuthor("Biblio", "Researcher", { biblioIdAlias: "@dissertation2" });
    cy.addSupervisor("Biblio", "Librarian", {
      biblioIdAlias: "@dissertation2",
    });

    cy.setUpPublication("Dissertation", {
      biblioIDAlias: "dissertation3",
      title: `Dissertation 3 ${randomTitleSuffix}`,
      shouldWaitForIndex: true,
    });
    cy.addAuthor("Biblio", "Librarian", { biblioIdAlias: "@dissertation3" });
    cy.addSupervisor("Biblio", "Researcher", {
      biblioIdAlias: "@dissertation3",
    });

    // Extra time for ES to index this
    cy.wait(1000);

    // Actual test
    cy.then(function () {
      cy.visit("/publication", { qs: { q: randomTitleSuffix } });

      // Test "All"
      cy.get(".card")
        .should("contain", "Showing 4 publications")
        .find(".list-group")
        .should("contain", this.book)
        .should("contain", this.dissertation1)
        .should("contain", this.dissertation2)
        .should("contain", this.dissertation3);

      // Test "My publications"
      cy.contains(".nav-tabs .nav-item a", "My publications").click();
      cy.get(".card")
        .should("contain", "Showing 1 publications")
        .find(".list-group")
        .should("not.contain", this.dissertation1)
        .should("not.contain", this.book)
        .should("contain", this.dissertation2)
        .should("not.contain", this.dissertation3);

      // Test "Supervised by me"
      cy.contains(".nav-tabs .nav-item a", "Supervised by me").click();
      cy.get(".card")
        .should("contain", "Showing 2 publications")
        .find(".list-group")
        .should("contain", this.dissertation1)
        .should("not.contain", this.book)
        .should("not.contain", this.dissertation2)
        .should("contain", this.dissertation3);

      // Test "Registered by me"
      cy.contains(".nav-tabs .nav-item a", "Registered by me").click();
      cy.get(".card")
        .should("contain", "Showing 3 publications")
        .find(".list-group")
        .should("not.contain", this.dissertation1)
        .should("contain", this.book)
        .should("contain", this.dissertation2)
        .should("contain", this.dissertation3);
    });
  });
});
