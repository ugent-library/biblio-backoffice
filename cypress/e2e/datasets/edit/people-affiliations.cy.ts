import { testFormAccessibility } from "support/util";

describe("Editing dataset people & affiliations", () => {
  beforeEach(() => {
    cy.loginAsResearcher("researcher1");
  });

  describe("Creators", () => {
    it("should be possible to add and delete creators", () => {
      cy.setUpDataset();
      cy.visitDataset();

      cy.contains(".nav-link", "People & Affiliations").click();
      cy.get("#authors .card-body").should(
        "contain",
        "Add at least one UGent creator.",
      );

      cy.intercept("/dataset/*/contributors/author/suggestions*").as("suggest");
      cy.intercept("/dataset/*/contributors/author/confirm-create*").as(
        "confirmCreate",
      );

      cy.contains(".btn", "Add creator").click();

      cy.ensureModal("Add creator").within(() => {
        cy.setFieldByLabel("First name", "Jame");
        cy.wait("@suggest");
        cy.setFieldByLabel("Last name", "Dow");
        cy.wait("@suggest");

        cy.contains(".btn", "Add external creator").click();
        cy.wait("@confirmCreate");

        // Made an error, let's go back
        cy.contains("Review creator information").should("be.visible");
        cy.contains("Jame Dow").should("be.visible");

        cy.contains(".btn", "Back to search").click();

        cy.setFieldByLabel("First name", "Jane");
        cy.wait("@suggest");

        cy.contains(".btn", "Add external creator").click();
        cy.wait("@confirmCreate");

        cy.contains("Review creator information").should("be.visible");
        cy.contains("Jane Dow").should("be.visible");

        cy.contains(".btn", "Save and add next").click();
      });

      cy.ensureModal("Add creator").within(() => {
        cy.setFieldByLabel("First name", "John");
        cy.wait("@suggest");
        cy.setFieldByLabel("Last name", "Doe");
        cy.wait("@suggest");

        cy.contains(".btn", "Add creator").click();
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
          cy.contains("Are you sure you want to remove this creator?").should(
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
          cy.contains("Are you sure you want to remove this creator?").should(
            "be.visible",
          );
        })
        .closeModal("Delete");
      cy.ensureNoModal();

      cy.contains("#authors", "Jane Dow").should("not.exist");
      cy.contains("#authors", "John Doe").should("not.exist");
      cy.get("#authors .card-body").should(
        "contain",
        "Add at least one UGent creator.",
      );
    });

    it("should not be possible to delete the last UGent creator of a published dataset", () => {
      cy.setUpDataset({ prepareForPublishing: true });
      cy.visitDataset();

      cy.contains(".btn", "Publish to Biblio").click();
      cy.ensureModal("Are you sure?").closeModal("Publish");
      cy.ensureToast("Dataset was successfully published.").closeToast();

      // Add other external creator first
      cy.addCreator("Jane", "Doe", { external: true });

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

    it("should be possible to add creator without first name", () => {
      cy.setUpDataset();
      cy.visitDataset();

      cy.updateFields(
        "Creators",
        () => {
          cy.setFieldByLabel("Last name", "Dow");

          cy.contains("Dow External, non-UGent")
            .closest(".list-group-item")
            .contains(".btn", "Add external creator")
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

    it("should be possible to add creator without last name", () => {
      cy.setUpDataset();
      cy.visitDataset();

      cy.updateFields(
        "Creators",
        () => {
          cy.setFieldByLabel("First name", "Jane");

          cy.contains("Jane External, non-UGent")
            .closest(".list-group-item")
            .contains(".btn", "Add external creator")
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

    it("should have clickable labels in the add/edit creator dialog", () => {
      cy.setUpDataset();
      cy.visitDataset();

      cy.updateFields(
        "Creators",
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
          cy.contains(".btn", "Add external creator").click();
        },
        true,
      );

      cy.contains("#contributors-author-body table tbody tr", "Jane Dow")
        .find(".if-edit")
        .click();

      cy.ensureModal("Edit or change creator").within(() => {
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
      cy.setUpDataset();
      cy.visitDataset();

      cy.contains(".nav .nav-item", "People & Affiliations").click();

      cy.get("#departments .card-body").should(
        "contain",
        "Add at least one department.",
      );

      cy.get("#departments").contains(".btn", "Add department").click();

      cy.ensureModal("Select departments").within(() => {
        cy.intercept("/dataset/*/departments/suggestions?*").as(
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
            "Are you sure you want to remove this department from the dataset?",
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
