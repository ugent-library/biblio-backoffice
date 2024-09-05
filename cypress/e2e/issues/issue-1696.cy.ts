// https://github.com/ugent-library/biblio-backoffice/issues/1696

describe("Issue #1696: <p> and </p> tags showing in project description", () => {
  it("should strip HTML from project descriptions in publications", () => {
    cy.loginAsResearcher("researcher1");

    cy.setUpPublication();
    cy.visitPublication();

    verifyProjectDescription();
  });

  it("should strip HTML from project descriptions in datasets", () => {
    cy.loginAsResearcher("researcher1");

    cy.setUpDataset();
    cy.visitDataset();

    verifyProjectDescription();
  });

  function verifyProjectDescription() {
    cy.contains(".card", "Project").contains(".btn", "Add project").click();

    cy.ensureModal("Select projects").within(() => {
      cy.setFieldByLabel("Search project", "001D07903");

      cy.get("#project-suggestions .list-group-item")
        .should("have.length", 1)
        .find(".list-group-item-main .text-muted")
        .invoke("text")
        .should(
          "start.with",
          "This interdisciplinary research project will try",
        )
        .and("not.contain", "<")
        .and("not.contain", ">");
    });
  }
});
