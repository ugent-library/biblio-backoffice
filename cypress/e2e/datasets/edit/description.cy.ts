import { testFocusForForm } from "support/util";

describe("Editing dataset description", () => {
  beforeEach(() => {
    cy.loginAsResearcher();

    cy.setUpDataset();
    cy.visitDataset();
  });

  describe("Dataset details", () => {
    it("should have clickable labels in the dataset form", () => {
      cy.updateFields("Dataset details", () => {
        testFocusForForm(
          {
            "input[type=text][name=title]": "Title",
            "select[name=identifier_type]": "Persistent identifier type",
            "input[type=text][name=identifier]": "Identifier",

            "select[name=language]": "Languages",
            "input[type=text][name=year]": "Publication year",
            "input[type=text][name=publisher]": "Publisher",

            "input[type=text][name=format]": "Data format",
            ".tags:has(textarea#keyword) tags span.tagify__input[contenteditable]":
              "Keywords",

            "select[name=license]": "License",

            "select[name=access_level]": "Access level",
          },
          "Title",
          ["textarea[data-input-name=keyword]"],
        );

        // Also display the "Other license" field
        cy.intercept("PUT", "/dataset/*/details/edit/refresh*").as(
          "refreshForm",
        );
        cy.setFieldByLabel("License", "The license is not listed here");
        cy.wait("@refreshForm");

        testFocusForForm(
          {
            "input[type=text][name=title]": "Title",
            "select[name=identifier_type]": "Persistent identifier type",
            "input[type=text][name=identifier]": "Identifier",

            "select[name=language]": "Languages",
            "input[type=text][name=year]": "Publication year",
            "input[type=text][name=publisher]": "Publisher",

            "input[type=text][name=format]": "Data format",
            ".tags:has(textarea#keyword) tags span.tagify__input[contenteditable]":
              "Keywords",

            "select[name=license]": "License",
            "input[type=text][name=other_license]": "Other license",

            "select[name=access_level]": "Access level",
          },
          undefined,
          ["textarea[data-input-name=keyword]"],
        );
      });
    });

    it("should not set autofocus when popup is refreshed", () => {
      cy.updateFields("Dataset details", () => {
        cy.focused().should("have.attr", "id", "title");

        cy.intercept("/dataset/*/details/edit/refresh*", (req) => {
          req.on("response", (res) => {
            // Pre-check of assertion so command log doesn't get bloated with massive HTML blocks
            if (
              typeof res.body === "string" &&
              res.body.includes("autofocus")
            ) {
              expect(res.body).to.not.contain("autofocus");
            }
          });
        }).as("refreshForm");

        cy.setFieldByLabel("License", "The license is not listed here");
        cy.wait("@refreshForm");
        cy.focused().should("have.attr", "id", "license");

        cy.setFieldByLabel("Access level", "Restricted access");
        cy.wait("@refreshForm");
        cy.focused().should("have.attr", "id", "access_level");

        cy.get("@refreshForm.all").should("have.length", 2);

        cy.setFieldByLabel("Publication year", "ABCD");
        cy.contains(".btn", "Save").click();
        cy.contains(
          ".alert-danger",
          "Publication year is an invalid value",
        ).should("be.visible");
        cy.focused().should("have.length", 0);
      });
    });
  });

  describe("Projects", () => {
    it("should be possible to add and delete projects", () => {
      cy.get("#projects-body").should("contain", "No projects");

      cy.contains(".card", "Project").contains(".btn", "Add project").click();

      cy.ensureModal("Select projects").within(() => {
        cy.intercept("/dataset/*/projects/suggestions?*").as("suggestProject");

        cy.setFieldByLabel("Search project", "001D07903");
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

    it("should have clickable labels in the project dialog", () => {
      cy.contains(".card", "Project").contains(".btn", "Add project").click();

      cy.ensureModal("Select projects").within(() => {
        cy.get("#project-q").should("be.focused");

        testFocusForForm(
          {
            "#project-q": "Search project",
          },
          "Search project",
        );
      });
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
        testFocusForForm(
          {
            "textarea[name=text]": "Abstract",
            "select[name=lang]": "Language",
          },
          "Abstract",
        );
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
        testFocusForForm(
          {
            "input[type=text][name=url]": "URL",
            "select[name=relation]": "Relation",
            "input[type=text][name=description]": "Description",
          },
          "URL",
        );
      });
    });
  });
});
