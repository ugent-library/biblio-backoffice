// https://github.com/ugent-library/biblio-backoffice/issues/1724

describe("Issue #1724: Cannot remove projects from records", () => {
  beforeEach(() => {
    cy.loginAsResearcher("researcher1");
  });

  it("should be possible to remove a project from a publication that contains slashes in its ID", () => {
    cy.setUpPublication();
    cy.visitPublication();

    cy.get("#projects-body .list-group-item").should("have.length", 0);

    cy.contains(".btn", "Add project").click();
    cy.ensureModal("Select projects").within(() => {
      cy.intercept("/publication/*/projects/suggestions?q=*").as(
        "suggestProject",
      );

      cy.setFieldByLabel(
        "Search project",
        "Study of the exclusive electron scattering experiments.",
      );

      cy.wait("@suggestProject");

      cy.get(".list-group-item")
        .should("have.length", 1)
        .should(
          "contain",
          "Study of the exclusive electron scattering experiments.",
        )
        .should("contain", "IWETO ID BOF/COR/2022/009")
        .contains(".btn", "Add project")
        .click();
    });

    cy.ensureNoModal();

    cy.get("#projects-body .list-group-item")
      .should("have.length", 1)
      .should(
        "contain",
        "Study of the exclusive electron scattering experiments.",
      )
      .should("contain", "IWETO ID BOF/COR/2022/009")
      .within(() => {
        cy.get(".btn .if-more").click();
        cy.contains(".dropdown-item", "Remove from publication").click();
      });

    cy.ensureModal("Confirm deletion").closeModal("Delete");

    cy.ensureNoModal();
    cy.get("#projects-body .list-group-item").should("have.length", 0);
  });

  it("should be possible to remove a project from a dataset that contains slashes in its ID", () => {
    cy.setUpDataset();
    cy.visitDataset();

    cy.get("#projects-body .list-group-item").should("have.length", 0);

    cy.contains(".btn", "Add project").click();
    cy.ensureModal("Select projects").within(() => {
      cy.intercept("/dataset/*/projects/suggestions?q=*").as("suggestProject");

      cy.setFieldByLabel(
        "Search project",
        "Study of the exclusive electron scattering experiments.",
      );

      cy.wait("@suggestProject");

      cy.get(".list-group-item")
        .should("have.length", 1)
        .should(
          "contain",
          "Study of the exclusive electron scattering experiments.",
        )
        .should("contain", "IWETO ID BOF/COR/2022/009")
        .contains(".btn", "Add project")
        .click();
    });

    cy.ensureNoModal();

    cy.get("#projects-body .list-group-item")
      .should("have.length", 1)
      .should(
        "contain",
        "Study of the exclusive electron scattering experiments.",
      )
      .should("contain", "IWETO ID BOF/COR/2022/009")
      .within(() => {
        cy.get(".btn .if-more").click();
        cy.contains(".dropdown-item", "Remove from dataset").click();
      });

    cy.ensureModal("Confirm deletion").closeModal("Delete");

    cy.ensureNoModal();
    cy.get("#projects-body .list-group-item").should("have.length", 0);
  });
});
