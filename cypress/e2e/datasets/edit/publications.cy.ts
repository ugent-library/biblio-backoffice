import { getRandomText } from "support/util";

describe("Editing publication datasets", () => {
  beforeEach(() => {
    cy.loginAsResearcher();
  });

  it("should be possible to add and delete related publications", () => {
    const title = getRandomText();
    cy.setUpPublication("Miscellaneous", {
      title,
      prepareForPublishing: true,
    });
    cy.visitPublication();
    cy.contains(".btn", "Publish to Biblio").click();
    cy.ensureModal("Are you sure?").closeModal("Publish");

    cy.setUpDataset();
    cy.visitDataset();

    cy.contains(".nav .nav-item", "Publications").click();

    cy.get("#publications-body").should("contain", "No publications");

    cy.contains(".card", "Related publications")
      .contains(".btn", "Add publication")
      .click();

    cy.ensureModal("Select publications").within(() => {
      cy.intercept("/dataset/*/publications/suggestions?*").as(
        "suggestPublication",
      );

      cy.getLabel("Search").next("input").type(title);
      cy.wait("@suggestPublication");

      cy.contains(".list-group-item", title)
        .contains(".btn", "Add publication")
        .click();
    });
    cy.ensureNoModal();

    cy.get("#publications-body")
      .contains(".list-group-item", title)
      .find(".if-more")
      .click();
    cy.contains(".dropdown-item", "Remove from dataset").click();

    cy.ensureModal("Confirm deletion")
      .within(() => {
        cy.get(".modal-body").should(
          "contain",
          "Are you sure you want to remove this publication from the dataset?",
        );
      })
      .closeModal("Delete");
    cy.ensureNoModal();

    cy.get("#publications-body").should("contain", "No publications");
  });
});
