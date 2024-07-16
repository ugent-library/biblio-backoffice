import { getRandomText } from "support/util";

describe("Publication import", () => {
  it("should be possible to import publications by DOI", () => {
    cy.loginAsResearcher();
    cy.visit("/");

    cy.contains(".btn", "Add research").click();
    cy.contains(".dropdown-item", "Add publication").click();

    // Step 1a
    cy.contains(".bc-toolbar-title", "Start: add publication(s)")
      .should("be.visible")
      .prev()
      .should("have.text", "Step 1");
    cy.get(".c-stepper__step").as("steps").should("have.length", 3);
    cy.get("@steps").eq(0).should("have.class", "c-stepper__step--active");
    cy.get("@steps").eq(1).should("not.have.class", "c-stepper__step--active");
    cy.get("@steps").eq(2).should("not.have.class", "c-stepper__step--active");

    cy.contains(".card", "Import your publication via an identifier")
      .contains(".btn", "Add")
      .click();

    // Step 1b
    cy.contains(".bc-toolbar-title", "Add publication(s)")
      .should("be.visible")
      .prev()
      .should("have.text", "Step 1");
    cy.get(".c-stepper__step").as("steps").should("have.length", 4);
    cy.get("@steps").eq(0).should("have.class", "c-stepper__step--active");
    cy.get("@steps").eq(1).should("not.have.class", "c-stepper__step--active");
    cy.get("@steps").eq(2).should("not.have.class", "c-stepper__step--active");
    cy.get("@steps").eq(3).should("not.have.class", "c-stepper__step--active");

    cy.get("select[name=source]").should("have.value", "crossref"); // crossref = DOI
    cy.get("input[name=identifier]").type("10.1016/j.ese.2024.100396");
    cy.contains(".btn", "Add publication(s)").click();

    // Step 2
    cy.contains(".bc-toolbar-title", "Complete Description")
      .should("be.visible")
      .prev()
      .should("have.text", "Step 2");
    cy.get(".c-stepper__step").as("steps").should("have.length", 4);
    cy.get("@steps").eq(0).should("have.class", "c-stepper__step--complete");
    cy.get("@steps").eq(1).should("have.class", "c-stepper__step--active");
    cy.get("@steps").eq(2).should("not.have.class", "c-stepper__step--active");
    cy.get("@steps").eq(3).should("not.have.class", "c-stepper__step--active");

    cy.get("#summary").should(
      "contain",
      "Synergies from off-gas analysis and mass balances for wastewater treatment — Some personal reflections on our experiences",
    );
    cy.contains(".btn", "Complete Description").click();

    // Step 3
    cy.contains(".bc-toolbar-title", "Publish to Biblio")
      .should("be.visible")
      .prev()
      .should("have.text", "Step 3");
    cy.get(".c-stepper__step").as("steps").should("have.length", 4);
    cy.get("@steps").eq(0).should("have.class", "c-stepper__step--complete");
    cy.get("@steps").eq(1).should("have.class", "c-stepper__step--complete");
    cy.get("@steps").eq(2).should("have.class", "c-stepper__step--active");
    cy.get("@steps").eq(3).should("not.have.class", "c-stepper__step--active");

    cy.extractBiblioId();

    cy.contains(".btn", "Save as draft").click();

    cy.location("pathname").should("eq", "/publication");

    cy.visitPublication();
  });

  it("should show an error toast if the DOI is invalid", () => {
    cy.loginAsResearcher();

    cy.visit("/add-publication");

    cy.contains(".card", "Import your publication via an identifier")
      .contains(".btn", "Add")
      .click();

    cy.get("input[name=identifier]").type("SOME/random/DOI.123");
    cy.contains(".btn", "Add publication(s)").click();

    cy.ensureToast(
      "Sorry, something went wrong. Could not import the publication",
    );
  });

  it("should detect duplicates by DOI", () => {
    const DOI = "10.2307/2323707";

    // First clean up existing publications with the same DOI
    cy.loginAsLibrarian();
    cy.switchMode("Librarian");
    const selector =
      ".card .card-body .list-group .list-group-item .c-button-toolbar .dropdown .dropdown-item:contains('Delete')";

    deleteNextPublication();

    function deleteNextPublication() {
      cy.visit("/publication", { qs: { q: DOI, "page-size": 1000 } }).then(
        () => {
          const deleteButton = Cypress.$(selector).first();

          if (deleteButton.length > 0) {
            cy.wrap(deleteButton).click({ force: true });

            cy.intercept("DELETE", "/publication/*").as("deletePublication");
            cy.ensureModal("Confirm deletion").closeModal("Delete");
            cy.wait("@deletePublication").then(deleteNextPublication);
          }
        },
      );
    }

    // Actual test starts here
    cy.loginAsResearcher();

    // First make and publish the first publication manually
    const title = getRandomText();
    cy.setUpPublication("Miscellaneous", {
      title,
      otherFields: { doi: DOI },
      publish: true,
    });

    // Some extra time for the dataset to be indexed
    cy.wait(1000);

    // Make the second publication
    cy.visit("/add-publication");

    cy.contains(".card", "Import your publication via an identifier")
      .contains(".btn", "Add")
      .click();

    cy.get("input[name=identifier]").type(DOI);
    cy.contains(".btn", "Add publication(s)").click();

    cy.ensureModal("Are you sure you want to import this publication?").within(
      () => {
        cy.get(".modal-body").should(
          "contain.text",
          "Biblio contains another publication with the same DOI:",
        );

        cy.get(".list-group-item").should("have.length", 1);
        cy.get(".list-group-item-title").should("contain.text", title);

        cy.contains(".modal-footer .btn", "Import anyway")
          .should("be.visible")
          .should("have.class", "btn-danger");
      },
    );
  });
});
