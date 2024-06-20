import { testFocusForLabel } from "support/util";

describe("Editing dataset people & affiliations", () => {
  beforeEach(() => {
    cy.loginAsResearcher();
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

      cy.updateFields(
        "Creators",
        () => {
          cy.intercept("/dataset/*/contributors/author/suggestions?*").as(
            "suggestContributor",
          );

          cy.setFieldByLabel("First name", "Jane");
          cy.wait("@suggestContributor");
          cy.setFieldByLabel("Last name", "Doe");
          cy.wait("@suggestContributor");

          cy.contains(".btn", "Add external creator").click({
            scrollBehavior: false,
          });
        },
        true,
      );

      cy.contains("#authors tr", "Jane Doe").find(".btn .if-delete").click();

      cy.ensureModal("Confirm deletion")
        .within(() => {
          cy.contains("Are you sure you want to remove this creator?").should(
            "be.visible",
          );
        })
        .closeModal("Delete");
      cy.ensureNoModal();

      cy.contains("#authors", "Jane Doe").should("not.exist");
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
      cy.addCreator("Jane", "Doe", true);

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

      cy.updateFields("Creators", () => {
        testFocusForLabel("First name", 'input[name="first_name"]', true);
        testFocusForLabel("Last name", 'input[name="last_name"]');
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

        cy.contains(".list-group-item", "Department ID: LW17")
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

        cy.contains(".list-group-item", "Department ID: DI62")
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
