import { getRandomText } from "support/util";

describe("The publication search (for researchers)", () => {
  it("should display filtered navigation tabs with all publications, my publications, supervised by me and registered by me", () => {
    const randomTitleSuffix = getRandomText();

    // Setup
    cy.loginAsResearcher("researcher1");
    cy.setUpPublication("Dissertation", {
      biblioIDAlias: "dissertation1",
      title: `Dissertation 1 ${randomTitleSuffix}`,
      shouldWaitForIndex: true,
    });
    cy.addAuthor("Biblio", "Researcher1", { biblioIDAlias: "@dissertation1" });
    cy.addSupervisor("Biblio", "Researcher2", {
      biblioIDAlias: "@dissertation1",
    });

    cy.loginAsResearcher("researcher2");
    cy.setUpPublication("Book", {
      biblioIDAlias: "book",
      title: `Book ${randomTitleSuffix}`,
      shouldWaitForIndex: true,
    });
    cy.addAuthor("John", "Doe", { biblioIDAlias: "@book" });
    cy.addEditor("Biblio", "Researcher2", { biblioIDAlias: "@book" });

    cy.setUpPublication("Dissertation", {
      biblioIDAlias: "dissertation2",
      title: `Dissertation 2 ${randomTitleSuffix}`,
      shouldWaitForIndex: true,
    });
    cy.addAuthor("Biblio", "Researcher2", { biblioIDAlias: "@dissertation2" });
    cy.addSupervisor("Biblio", "Researcher1", {
      biblioIDAlias: "@dissertation2",
    });

    cy.setUpPublication("Dissertation", {
      biblioIDAlias: "dissertation3",
      title: `Dissertation 3 ${randomTitleSuffix}`,
      shouldWaitForIndex: true,
    });
    cy.addAuthor("Biblio", "Researcher1", { biblioIDAlias: "@dissertation3" });
    cy.addSupervisor("Biblio", "Researcher2", {
      biblioIDAlias: "@dissertation3",
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
