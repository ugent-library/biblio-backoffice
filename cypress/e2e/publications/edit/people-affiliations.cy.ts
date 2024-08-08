import { testFormAccessibility } from "support/util";

describe("Editing publication people & affiliations", () => {
  beforeEach(() => {
    cy.loginAsResearcher("researcher1");

    cy.intercept("/publication/*/contributors/author/suggestions*").as(
      "suggest",
    );
    cy.intercept("/publication/*/contributors/author/confirm-create*").as(
      "confirmCreate",
    );
  });

  describe("Authors", () => {
    it("should be possible to add and delete authors", () => {
      cy.setUpPublication();
      cy.visitPublication();

      cy.contains(".nav-link", "People & Affiliations").click();
      cy.get("#authors .card-body").should(
        "contain",
        "Add at least one UGent author.",
      );

      cy.contains(".btn", "Add author").click();

      cy.ensureModal("Add author").within(() => {
        cy.setFieldByLabel("First name", "Jame");
        cy.wait("@suggest");
        cy.setFieldByLabel("Last name", "Dow");
        cy.wait("@suggest");

        cy.contains(".btn", "Add external author").click();
        cy.wait("@confirmCreate");

        // Made an error, let's go back
        cy.contains("Review author information").should("be.visible");
        cy.contains("Jame Dow").should("be.visible");

        cy.contains(".btn", "Back to search").click();

        cy.setFieldByLabel("First name", "Jane");
        cy.wait("@suggest");

        cy.contains(".btn", "Add external author").click();
        cy.wait("@confirmCreate");

        cy.contains("Review author information").should("be.visible");
        cy.contains("Jane Dow").should("be.visible");

        cy.contains(".btn", "Save and add next").click();
      });

      cy.ensureModal("Add author").within(() => {
        cy.setFieldByLabel("First name", "John");
        cy.wait("@suggest");
        cy.setFieldByLabel("Last name", "Doe");
        cy.wait("@suggest");

        cy.contains(".btn", "Add author").click();
        cy.wait("@confirmCreate");

        cy.contains(".btn", /^Save$/).click();
      });

      cy.get("#authors tbody tr")
        .should("have.length", 2)
        .contains("tr", "Jane Dow")
        .find(".btn .if-delete")
        .click();

      cy.ensureModal("Confirm deletion")
        .within(() => {
          cy.contains("Are you sure you want to remove this author?").should(
            "be.visible",
          );
        })
        .closeModal("Delete");
      cy.ensureNoModal();

      cy.get("#authors tbody tr")
        .should("have.length", 1)
        .contains("tr", "John Doe")
        .find(".btn .if-delete")
        .click();

      cy.ensureModal("Confirm deletion")
        .within(() => {
          cy.contains("Are you sure you want to remove this author?").should(
            "be.visible",
          );
        })
        .closeModal("Delete");
      cy.ensureNoModal();

      cy.contains("#authors", "Jane Dow").should("not.exist");
      cy.contains("#authors", "John Doe").should("not.exist");
      cy.get("#authors .card-body").should(
        "contain",
        "Add at least one UGent author.",
      );
    });

    it("should not be possible to delete the last UGent author of a published publication", () => {
      cy.setUpPublication(undefined, { prepareForPublishing: true });
      cy.visitPublication();

      cy.contains(".btn", "Publish to Biblio").click();
      cy.ensureModal("Are you sure?").closeModal("Publish");
      cy.ensureToast("Publication was successfully published.").closeToast();

      // Add other external author first
      cy.addAuthor("Jane", "Doe", { external: true });

      cy.contains(".nav-item", "People & Affiliations").click();
      cy.contains("#authors tr", "John Doe").find(".btn .if-delete").click();
      cy.ensureModal("Confirm deletion").closeModal("Delete");

      cy.ensureModal(
        "Can't delete this contributor due to the following errors",
      ).within(() => {
        cy.contains(
          ".alert-danger",
          "At least one UGent author is required",
        ).should("be.visible");
      });
    });

    it("should be possible to add author without first name", () => {
      cy.setUpPublication("Book");
      cy.visitPublication();

      cy.updateFields(
        "Authors",
        () => {
          cy.setFieldByLabel("Last name", "Dow");

          cy.contains("Dow External, non-UGent")
            .closest(".list-group-item")
            .contains(".btn", "Add external author")
            .click();
        },
        true,
      );

      cy.get(".card#authors")
        .find("#contributors-author-table tr")
        .last()
        .find("td")
        .first()
        .find(".bc-avatar-text")
        .should("contain", "[missing] Dow");
    });

    it("should be possible to add author without last name", () => {
      cy.setUpPublication("Book");
      cy.visitPublication();

      cy.updateFields(
        "Authors",
        () => {
          cy.setFieldByLabel("First name", "Jane");

          cy.contains("Jane External, non-UGent")
            .closest(".list-group-item")
            .contains(".btn", "Add external author")
            .click();
        },
        true,
      );

      cy.get(".card#authors")
        .find("#contributors-author-table tr")
        .last()
        .find("td")
        .first()
        .find(".bc-avatar-text")
        .should("contain", "Jane [missing]");
    });

    it("should have clickable labels in the add/edit author dialog", () => {
      cy.setUpPublication();
      cy.visitPublication();

      cy.updateFields("Authors", () => {
        testFormAccessibility(
          {
            "input[name=first_name]": "First name",
            "input[name=last_name]": "Last name",
          },
          "First name",
        );

        cy.setFieldByLabel("First name", "John");
        cy.setFieldByLabel("Last name", "Doe");
        cy.wait("@suggest");
        cy.contains(".btn", "Add author").click();
        cy.wait("@confirmCreate");

        testFormAccessibility(
          {
            "select[name=credit_role]": "Roles",
          },
          "Roles",
        );
      });
    });
  });

  describe("Editors", () => {
    it("should be possible to add and delete editors", () => {
      cy.setUpPublication("Book");
      cy.visitPublication();

      cy.contains(".nav-link", "People & Affiliations").click();
      cy.get("#editors .card-body").should("contain", "No editors.");

      cy.updateFields(
        "Editors",
        () => {
          cy.intercept("/publication/*/contributors/editor/suggestions*").as(
            "suggest",
          );

          cy.setFieldByLabel("First name", "Jane");
          cy.wait("@suggest");
          cy.setFieldByLabel("Last name", "Doe");
          cy.wait("@suggest");

          cy.contains(".btn", "Add external editor").click({
            scrollBehavior: "nearest",
          });
        },
        true,
      );

      cy.contains("#editors tr", "Jane Doe").find(".btn .if-delete").click();

      cy.ensureModal("Confirm deletion")
        .within(() => {
          cy.contains("Are you sure you want to remove this editor?").should(
            "be.visible",
          );
        })
        .closeModal("Delete");
      cy.ensureNoModal();

      cy.contains("#editors", "Jane Doe").should("not.exist");
      cy.get("#editors .card-body").should("contain", "No editors.");
    });

    it("should have clickable labels in the add/edit editor dialog", () => {
      cy.setUpPublication();
      cy.visitPublication();

      cy.updateFields(
        "Editors",
        () => {
          testFormAccessibility(
            {
              "input[name=first_name]": "First name",
              "input[name=last_name]": "Last name",
            },
            "First name",
          );

          cy.setFieldByLabel("First name", "Jane");
          cy.setFieldByLabel("Last name", "Dow");
          cy.contains(".btn", "Add external editor").click();
        },
        true,
      );

      cy.contains("#contributors-editor-body table tbody tr", "Jane Dow")
        .find(".if-edit")
        .click();

      cy.ensureModal("Edit or change editor").within(() => {
        testFormAccessibility(
          {
            "input[name=first_name]": "First name",
            "input[name=last_name]": "Last name",
          },
          "First name",
        );
      });
    });
  });

  describe("Supervisors", () => {
    it("should be possible to add and delete supervisors", () => {
      cy.setUpPublication("Dissertation");
      cy.visitPublication();

      cy.contains(".nav-link", "People & Affiliations").click();
      cy.get("#supervisors .card-body").should("contain", "No supervisors.");

      cy.updateFields(
        "Supervisors",
        () => {
          cy.intercept(
            "/publication/*/contributors/supervisor/suggestions?*",
          ).as("suggest");

          cy.setFieldByLabel("First name", "Jane");
          cy.wait("@suggest");
          cy.setFieldByLabel("Last name", "Doe");
          cy.wait("@suggest");

          cy.contains(".btn", "Add external supervisor").click();
        },
        true,
      );

      cy.contains("#supervisors tr", "Jane Doe")
        .find(".btn .if-delete")
        .click();

      cy.ensureModal("Confirm deletion")
        .within(() => {
          cy.contains(
            "Are you sure you want to remove this supervisor?",
          ).should("be.visible");
        })
        .closeModal("Delete");
      cy.ensureNoModal();

      cy.contains("#supervisors", "Jane Doe").should("not.exist");
      cy.get("#supervisors .card-body").should("contain", "No supervisors.");
    });

    it("should have clickable labels in the add/edit supervisor dialog", () => {
      cy.setUpPublication("Dissertation");
      cy.visitPublication();

      cy.updateFields(
        "Supervisors",
        () => {
          testFormAccessibility(
            {
              "input[name=first_name]": "First name",
              "input[name=last_name]": "Last name",
            },
            "First name",
          );

          cy.setFieldByLabel("First name", "Jane");
          cy.setFieldByLabel("Last name", "Dow");
          cy.contains(".btn", "Add external supervisor").click();
        },
        true,
      );

      cy.contains("#contributors-supervisor-body table tbody tr", "Jane Dow")
        .find(".if-edit")
        .click();

      cy.ensureModal("Edit or change supervisor").within(() => {
        testFormAccessibility(
          {
            "input[name=first_name]": "First name",
            "input[name=last_name]": "Last name",
          },
          "First name",
        );
      });
    });
  });

  describe("Departments", () => {
    it("should be possible to add and delete departments", () => {
      cy.setUpPublication();
      cy.visitPublication();

      cy.contains(".nav .nav-item", "People & Affiliations").click();

      cy.get("#departments .card-body").should(
        "contain",
        "Add at least one department.",
      );

      cy.get("#departments").contains(".btn", "Add department").click();

      cy.ensureModal("Select departments").within(() => {
        cy.intercept("/publication/*/departments/suggestions?*").as(
          "suggestDepartment",
        );

        cy.getLabel("Search").next("input").should("be.focused").type("LW17");
        cy.wait("@suggestDepartment");

        cy.contains(".list-group-item", "Department ID LW17")
          .contains(".btn", "Add department")
          .click();
      });
      cy.ensureNoModal();

      cy.get("#departments-body .list-group-item-text h4")
        .map("textContent")
        .should("have.ordered.members", [
          "Department of Art, music and theatre sciences",
        ]);

      cy.get("#departments").contains(".btn", "Add department").click();

      cy.ensureModal("Select departments").within(() => {
        cy.getLabel("Search").next("input").should("be.focused").type("DI62");
        cy.wait("@suggestDepartment");

        cy.contains(".list-group-item", "Department ID DI62")
          .contains(".btn", "Add department")
          .click();
      });
      cy.ensureNoModal();

      cy.get("#departments-body .list-group-item-text h4")
        .map("textContent")
        .should("have.ordered.members", [
          "Department of Art, music and theatre sciences",
          "Biocenter AgriVet",
        ]);

      cy.contains(
        "#departments-body tr",
        "Department of Art, music and theatre sciences",
      )
        .find(".if-delete")
        .click();

      cy.ensureModal("Confirm deletion")
        .within(() => {
          cy.get(".modal-body").should(
            "contain",
            "Are you sure you want to remove this department from the publication?",
          );
        })
        .closeModal("Delete");
      cy.ensureNoModal();

      cy.get("#departments-body .list-group-item-text h4")
        .map("textContent")
        .should("have.ordered.members", ["Biocenter AgriVet"]);
    });
  });
});
