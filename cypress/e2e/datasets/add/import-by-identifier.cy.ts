import { getRandomText, testFormAccessibility } from "support/util";

describe("Dataset import", () => {
  it("should be possible to import datasets by DOI", () => {
    cy.login("researcher1");
    cy.visit("/");

    cy.contains(".btn", "Add research").click();
    cy.contains(".dropdown-item", "Add dataset").click();

    // Step 1a
    cy.contains("Step 1")
      .should("be.visible")
      .next()
      .should("have.text", "Add dataset");
    cy.get(".c-stepper__step").as("steps").should("have.length", 4);
    cy.get("@steps").eq(0).should("have.class", "c-stepper__step--active");
    cy.get("@steps").eq(1).should("not.have.class", "c-stepper__step--active");
    cy.get("@steps").eq(2).should("not.have.class", "c-stepper__step--active");
    cy.get("@steps").eq(3).should("not.have.class", "c-stepper__step--active");

    cy.contains("Register your dataset via a DOI").click();
    cy.contains(".btn", "Add dataset").click();

    // Step 1b
    cy.contains("Step 1")
      .should("be.visible")
      .next()
      .should("have.text", "Add dataset");
    cy.get(".c-stepper__step").as("steps").should("have.length", 4);
    cy.get("@steps").eq(0).should("have.class", "c-stepper__step--active");
    cy.get("@steps").eq(1).should("not.have.class", "c-stepper__step--active");
    cy.get("@steps").eq(2).should("not.have.class", "c-stepper__step--active");
    cy.get("@steps").eq(3).should("not.have.class", "c-stepper__step--active");

    testFormAccessibility({ "input[name=identifier]": "DOI" });

    cy.get("input[name=identifier]").type("10.6084/M9.FIGSHARE.22067864.V1");
    cy.contains(".btn", "Add dataset").click();

    // Step 2
    cy.contains("Step 2")
      .should("be.visible")
      .next()
      .should("have.text", "Complete Description");
    cy.get(".c-stepper__step").as("steps").should("have.length", 4);
    cy.get("@steps").eq(0).should("have.class", "c-stepper__step--complete");
    cy.get("@steps").eq(1).should("have.class", "c-stepper__step--active");
    cy.get("@steps").eq(2).should("not.have.class", "c-stepper__step--active");
    cy.get("@steps").eq(3).should("not.have.class", "c-stepper__step--active");

    cy.get("#summary").should(
      "contain",
      "SI: Janssen et al. (2023): modelling copper toxicity on brook trout populations",
    );
    cy.contains(".btn", "Complete Description").click();

    // Step 3
    cy.contains("Step 3")
      .should("be.visible")
      .next()
      .should("have.text", "Review and publish");
    cy.get(".c-stepper__step").as("steps").should("have.length", 4);
    cy.get("@steps").eq(0).should("have.class", "c-stepper__step--complete");
    cy.get("@steps").eq(1).should("have.class", "c-stepper__step--complete");
    cy.get("@steps").eq(2).should("have.class", "c-stepper__step--active");
    cy.get("@steps").eq(3).should("not.have.class", "c-stepper__step--active");

    cy.extractBiblioId();

    cy.contains(".btn", "Save as draft").click();

    cy.location("pathname").should("eq", "/dataset");

    cy.visitDataset();
  });

  it("should show an error toast if the DOI is invalid", () => {
    cy.login("researcher1");

    cy.visit("/add-dataset");

    cy.contains("Register your dataset via a DOI").click();
    cy.contains(".btn", "Add dataset").click();

    cy.get("input[name=identifier]").type("SOME/random/DOI.123");
    cy.contains(".btn", "Add dataset").click();

    cy.ensureToast("Sorry, something went wrong. Could not import the dataset");
  });

  it("should detect duplicates by DOI", () => {
    const DOI = "10.48804/A76XM9";

    // First clean up existing datasets with the same DOI
    cy.login("librarian1");
    cy.switchMode("Librarian");
    const selector =
      ".card .card-body .list-group .list-group-item .c-button-toolbar .dropdown .dropdown-item:contains('Delete')";

    deleteNextDataset();

    function deleteNextDataset() {
      cy.visit("/dataset", { qs: { q: DOI, "page-size": 1000 } }).then(() => {
        const deleteButton = Cypress.$(selector).first();

        if (deleteButton.length > 0) {
          cy.wrap(deleteButton).click({ force: true });

          cy.intercept("DELETE", "/dataset/*").as("deleteDataset");
          cy.ensureModal("Confirm deletion").closeModal("Delete");
          cy.wait("@deleteDataset").then(deleteNextDataset);
        }
      });
    }

    // Actual test starts here
    cy.login("researcher1");

    // First make and publish the first dataset manually
    const title = getRandomText();
    cy.setUpDataset({
      title,
      otherFields: {
        identifier_type: "DOI",
        identifier: DOI,
      },
      publish: true,
    });

    // Some extra time for the dataset to be indexed
    cy.wait(1000);

    // Make the second dataset (via DOI import)
    cy.visit("/add-dataset");

    cy.contains("Register your dataset via a DOI").click();
    cy.contains(".btn", "Add dataset").click();

    cy.get("input[name=identifier]").type(DOI);
    cy.contains(".btn", "Add dataset").click();

    cy.ensureModal("Are you sure you want to import this dataset?").within(
      () => {
        cy.get(".modal-body").should(
          "contain.text",
          "Biblio contains another dataset with the same DOI:",
        );

        cy.get(".list-group-item").should("have.length", 1);
        cy.get(".list-group-item-title").should("contain.text", title);

        cy.contains(".modal-footer .btn", "Import Anyway")
          .should("be.visible")
          .should("have.class", "btn-danger")
          .click();
      },
    );

    cy.ensureNoModal();
  });
});
