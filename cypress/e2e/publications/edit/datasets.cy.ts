import { getRandomText } from "support/util";

describe("Editing publication datasets", () => {
  beforeEach(() => {
    cy.loginAsResearcher("researcher1");
  });

  it("should be possible to add and delete related datasets", () => {
    const title = getRandomText();
    cy.setUpDataset({ title, prepareForPublishing: true });
    cy.visitDataset();
    cy.contains(".btn", "Publish to Biblio").click();
    cy.ensureModal("Are you sure?").closeModal("Publish");

    cy.setUpPublication();
    cy.visitPublication();

    cy.contains(".nav .nav-item", "Datasets").click();

    cy.get("#datasets-body").should("contain", "No datasets");

    cy.contains(".card", "Related datasets")
      .contains(".btn", "Add dataset")
      .click();

    cy.ensureModal("Select datasets").within(() => {
      cy.intercept("/publication/*/datasets/suggestions?*").as(
        "suggestDataset",
      );

      cy.getLabel("Search datasets")
        .next("input")
        .should("be.focused")
        .type(title);
      cy.wait("@suggestDataset");

      cy.contains(".list-group-item", title)
        .contains(".btn", "Add dataset")
        .click();
    });
    cy.ensureNoModal();

    cy.get("#datasets-body")
      .contains(".list-group-item", title)
      .find(".if-more")
      .click();
    cy.contains(".dropdown-item", "Remove from publication").click();

    cy.ensureModal("Confirm deletion")
      .within(() => {
        cy.get(".modal-body").should(
          "contain",
          "Are you sure you want to remove this dataset from the publication?",
        );
      })
      .closeModal("Delete");
    cy.ensureNoModal();

    cy.get("#datasets-body").should("contain", "No datasets");
  });
});
