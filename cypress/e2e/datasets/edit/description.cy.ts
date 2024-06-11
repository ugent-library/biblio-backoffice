import { testFocusForLabel } from "support/util";

describe("Editing dataset description", () => {
  beforeEach(() => {
    cy.loginAsResearcher();

    cy.setUpDataset();
    cy.visitDataset();
  });

  describe("Dataset details", () => {
    it("should have clickable labels in the dataset form", () => {
      cy.updateFields("Dataset details", () => {
        cy.intercept("PUT", "/dataset/*/details/edit/refresh*").as(
          "refreshForm",
        );
        cy.setFieldByLabel("License", "The license is not listed here");

        cy.wait("@refreshForm");

        testFocusForLabel("Title", 'input[type=text][name="title"]');
        testFocusForLabel(
          "Persistent identifier type",
          'select[name="identifier_type"]',
        );
        testFocusForLabel("Identifier", 'input[type=text][name="identifier"]');

        testFocusForLabel("Languages", 'select[name="language"]');
        testFocusForLabel("Publication year", 'input[type=text][name="year"]');
        testFocusForLabel("Publisher", 'input[type=text][name="publisher"]');

        testFocusForLabel("Data format", 'input[type=text][name="format"]');
        // Keywords field: tagify component doesn't support focussing by label

        testFocusForLabel("License", 'select[name="license"]');
        testFocusForLabel(
          "Other license",
          'input[type=text][name="other_license"]',
        );

        testFocusForLabel("Access level", 'select[name="access_level"]');
      });
    });
  });

  describe("Projects", () => {
    it("should be possible to add and delete projects", () => {
      cy.get("#projects-body").should("contain", "No projects");

      cy.contains(".card", "Project").contains(".btn", "Add project").click();

      cy.ensureModal("Select projects").within(() => {
        cy.intercept("/dataset/*/projects/suggestions?*").as("suggestProject");

        cy.getLabel("Search project").next("input").type("001D07903");
        cy.wait("@suggestProject");

        cy.contains(".list-group-item", "001D07903")
          .contains(".btn", "Add project")
          .click();
      });
      cy.ensureNoModal();

      cy.get("#projects-body")
        .contains(".list-group-item", "001D07903")
        .find(".if-more")
        .click();
      cy.contains(".dropdown-item", "Remove from dataset").click();

      cy.ensureModal("Confirm deletion")
        .within(() => {
          cy.get(".modal-body").should(
            "contain",
            "Are you sure you want to remove this project from the dataset?",
          );
        })
        .closeModal("Delete");
      cy.ensureNoModal();

      cy.get("#projects-body").should("contain", "No projects");
    });
  });

  describe("Abstracts", () => {
    it("should be possible to add, edit and delete abstracts", () => {
      cy.get("#abstracts-body").should("contain", "No abstracts");

      cy.contains(".btn", "Add abstract").click();
      cy.ensureModal("Add abstract")
        .within(() => {
          cy.setFieldByLabel("Abstract", "");
          cy.setFieldByLabel("Language", "Danish");
        })
        .closeModal("Add abstract");

      cy.ensureModal("Add abstract")
        .within(() => {
          cy.contains(".alert-danger", "Abstract text can't be empty").should(
            "be.visible",
          );
          cy.get("textarea[name=text]")
            .should("have.class", "is-invalid")
            .next(".invalid-feedback")
            .should("have.text", "Abstract text can't be empty");

          cy.setFieldByLabel("Abstract", "The initial abstract");
        })
        .closeModal("Add abstract");
      cy.ensureNoModal();

      cy.get("#abstracts-body")
        .find("table tbody tr")
        .as("row")
        .should("have.length", 1);

      cy.get("@row").should("contain", "The initial abstract");
      cy.get("@row").should("contain", "Danish");

      cy.get("@row").within(() => {
        cy.get(".if-more").click();

        cy.contains("button", "Edit").click();
      });

      cy.ensureModal("Edit abstract")
        .within(() => {
          cy.setFieldByLabel("Abstract", "");
          cy.setFieldByLabel("Language", "Northern Sami");
        })
        .closeModal("Update abstract");

      cy.ensureModal("Edit abstract")
        .within(() => {
          cy.contains(".alert-danger", "Abstract text can't be empty").should(
            "be.visible",
          );
          cy.get("textarea[name=text]")
            .should("have.class", "is-invalid")
            .next(".invalid-feedback")
            .should("have.text", "Abstract text can't be empty");

          cy.setFieldByLabel("Abstract", "The updated abstract");
        })
        .closeModal("Update abstract");
      cy.ensureNoModal();

      cy.get("@row").should("have.length", 1);

      cy.get("@row").should("not.contain", "The initial abstract");
      cy.get("@row").should("contain", "The updated abstract");
      cy.get("@row").should("not.contain", "Danish");
      cy.get("@row").should("contain", "Northern Sami");

      cy.get("@row").within(() => {
        cy.get(".if-more").click();

        cy.contains("button", "Delete").click();
      });

      cy.ensureModal("Confirm deletion")
        .within(() => {
          cy.get(".modal-body").should(
            "contain",
            "Are you sure you want to remove this abstract?",
          );
        })
        .closeModal("Delete");
      cy.ensureNoModal();

      cy.get("#abstracts-body").should("contain", "No abstracts");
    });

    it("should error when trying to edit/delete an abstract that was already deleted", () => {
      cy.updateFields(
        "Abstract",
        () => {
          cy.setFieldByLabel("Abstract", "The abstract text");
        },
        "Add abstract",
      );

      cy.get("#abstracts-body .if-more").click();
      cy.contains("#abstracts-body .dropdown-item", "Delete").click();
      cy.ensureModal("Confirm deletion")
        .within(() => {
          cy.contains(".btn", "Delete").triggerHtmx("hx-delete");
        })
        .closeModal("Delete");
      cy.ensureModal(null)
        .within(() => {
          cy.get(".modal-body").should(
            "contain",
            "Dataset has been modified by another user. Please reload the page.",
          );
        })
        .closeModal("Close");
      cy.ensureNoModal();

      cy.get("#abstracts-body .if-more").click();
      cy.contains("#abstracts-body .dropdown-item", "Edit").click();
      cy.ensureModal(null)
        .within(() => {
          cy.get(".modal-body").should(
            "contain",
            "Dataset has been modified by another user. Please reload the page.",
          );
        })
        .closeModal("Close");

      cy.get("#abstracts-body .if-more").click();
      cy.contains("#abstracts-body .dropdown-item", "Delete").click();
      cy.ensureModal(null).within(() => {
        cy.get(".modal-body").should(
          "contain",
          "Dataset has been modified by another user. Please reload the page.",
        );
      });
    });

    it("should have clickable labels in the Abstract dialog", () => {
      cy.updateFields("Abstract", () => {
        testFocusForLabel("Abstract", 'textarea[name="text"]');
        testFocusForLabel("Language", 'select[name="lang"]');
      });
    });
  });

  describe("Links", () => {
    it("should be possible to add, edit and delete links", () => {
      cy.get("#links-body").should("contain", "No links");

      cy.contains(".btn", "Add link").click();
      cy.ensureModal("Add link")
        .within(() => {
          cy.setFieldByLabel("URL", "https://www.ugent.be");
          cy.setFieldByLabel("Relation", "Related information");
          cy.setFieldByLabel("Description", "The initial website");
        })
        .closeModal("Add link");
      cy.ensureNoModal();

      cy.get("#links-body")
        .find("table tbody tr")
        .as("row")
        .should("have.length", 1);

      cy.get("@row").should("contain", "https://www.ugent.be");
      cy.get("@row").should("contain", "Related information");
      cy.get("@row").should("contain", "The initial website");

      cy.get("@row").within(() => {
        cy.get(".if-more").click();

        cy.contains("button", "Edit").click();
      });

      cy.ensureModal("Edit link")
        .within(() => {
          cy.setFieldByLabel("URL", "https://lib.ugent.be");
          cy.setFieldByLabel("Relation", "Accompanying website");
          cy.setFieldByLabel("Description", "The updated website");
        })
        .closeModal("Update link");
      cy.ensureNoModal();

      cy.get("@row").should("have.length", 1);

      cy.get("@row").should("not.contain", "https://www.ugent.be");
      cy.get("@row").should("contain", "https://lib.ugent.be");
      cy.get("@row").should("not.contain", "Related information");
      cy.get("@row").should("contain", "Accompanying website");
      cy.get("@row").should("not.contain", "The initial website");
      cy.get("@row").should("contain", "The updated website");

      cy.get("@row").within(() => {
        cy.get(".if-more").click();

        cy.contains("button", "Delete").click();
      });

      cy.ensureModal("Confirm deletion")
        .within(() => {
          cy.get(".modal-body").should(
            "contain",
            "Are you sure you want to remove this link?",
          );
        })
        .closeModal("Delete");
      cy.ensureNoModal();

      cy.get("#links-body").should("contain", "No links");
    });

    it("should have clickable labels in the Link dialog", () => {
      cy.updateFields("Link", () => {
        testFocusForLabel("URL", 'input[type=text][name="url"]');
        testFocusForLabel("Relation", 'select[name="relation"]');
        testFocusForLabel(
          "Description",
          'input[type=text][name="description"]',
        );
      });
    });
  });
});
