import { getRandomText } from "support/util";

describe("Publication import", () => {
  beforeEach(() => {
    cy.loginAsResearcher();
  });

  it("should be a possible to import publications by DOI", () => {
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

    cy.contains("Import your publication via an identifier").click();
    cy.contains(".btn", "Add publication(s)").click();

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

    cy.get('select[name="source"]').should("have.value", "crossref"); // crossref = DOI
    cy.get('input[name="identifier"]').type("10.1016/j.ese.2024.100396");
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

  it("should detect duplicates by DOI", () => {
    const title = getRandomText();
    cy.setUpPublication("Miscellaneous", { title, prepareForPublishing: true });
    cy.visitPublication();

    cy.updateFields(
      "Publication details",
      () => {
        cy.setFieldByLabel("DOI", "DOI/test/123.98");
      },
      true,
    );

    //10.1016/j.ese.2024.100396

    cy.contains(".btn", "Publish to Biblio").click();
    cy.ensureModal("Are you sure?").closeModal("Publish");

    // Make the second publication
    cy.visit("/add-publication");

    cy.contains("Import your publication via an identifier").click();
    cy.contains(".btn", "Add publication(s)").click();

    cy.get('input[name="identifier"]').type("DOI/test/123.98");
    cy.contains(".btn", "Add publication(s)").click();

    // TODO: convert this to a regular dialog so we can use cy.ensureModal here
    cy.get(".modal-dialog").within(() => {
      cy.get(".modal-header").should(
        "contain",
        "Are you sure you want to import this publication?",
      );

      cy.get(".modal-body").should(
        "contain.text",
        "Biblio contains another publication with the same DOI:",
      );

      cy.get(".list-group-item").should("have.length", 1);
      cy.get(".list-group-item-title").should("contain.text", title);

      cy.contains(".modal-footer .btn", "Import anyway")
        .should("be.visible")
        .should("have.class", "btn-danger");
    });
  });
});
