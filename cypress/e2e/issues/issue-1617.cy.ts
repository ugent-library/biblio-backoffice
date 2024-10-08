// https://github.com/ugent-library/biblio-backoffice/issues/1617

describe("Issue #1617: Librarian tags field does not work when you hit save before blurring the field", () => {
  beforeEach(() => {
    cy.loginAsLibrarian("librarian1");
  });

  describe("for publications", () => {
    beforeEach(() => {
      cy.setUpPublication();
      cy.visitPublication();
    });

    it("should save keywords when you directly hit save", () => {
      cy.updateFields(
        "Additional information",
        () => {
          cy.setFieldByLabel("Keywords", "Keyword 1{enter}Keyword 2");
        },
        true,
      );

      cy.contains("#additional-info-body .list-group-item", "Keywords")
        .find(".list-inline .badge")
        .should("have.length", 2)
        .should("contain", "Keyword 1")
        .should("contain", "Keyword 2");
    });

    it("should save librarian tags when you directly hit save", () => {
      cy.updateFields(
        "Librarian tags",
        () => {
          cy.setFieldByLabel("Librarian tags", "Tag 1{enter}Tag 2;Tag 3");
        },
        true,
      );

      cy.get("#reviewer-tags-body .badge-list .badge")
        .should("have.length", 3)
        .should("contain", "Tag 1")
        .should("contain", "Tag 2")
        .should("contain", "Tag 3");
    });
  });

  describe("for datasets", () => {
    beforeEach(() => {
      cy.setUpDataset();
      cy.visitDataset();
    });

    it("should save keywords when you directly hit save", () => {
      cy.updateFields(
        "Dataset details",
        () => {
          cy.setFieldByLabel("Keywords", "Keyword 1{enter}Keyword 2");
        },
        true,
      );

      cy.contains("#details-body .list-group-item .row", "Keywords")
        .find(".list-inline .badge")
        .should("have.length", 2)
        .should("contain", "Keyword 1")
        .should("contain", "Keyword 2");
    });

    it("should save librarian tags when you directly hit save", () => {
      cy.updateFields(
        "Librarian tags",
        () => {
          cy.setFieldByLabel("Librarian tags", "Tag 1{enter}Tag 2;Tag 3");
        },
        true,
      );

      cy.get("#reviewer-tags-body .badge-list .badge")
        .should("have.length", 3)
        .should("contain", "Tag 1")
        .should("contain", "Tag 2")
        .should("contain", "Tag 3");
    });
  });
});
