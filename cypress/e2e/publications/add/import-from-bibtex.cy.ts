describe("Publication import", () => {
  it("should be possible to import publications from BibTeX and save as draft", () => {
    cy.loginAsResearcher("researcher1");

    cy.visit("/");

    cy.contains("Biblio Publications").click();

    cy.contains("Add Publication").click();

    // Add publication(s)
    cy.contains("Step 1").should("be.visible");
    cy.contains(".bc-toolbar-title", "Start: add publication(s)").should(
      "be.visible",
    );
    cy.get(".c-stepper__step").as("steps").should("have.length", 3);
    cy.get("@steps").eq(0).should("have.class", "c-stepper__step--active");
    cy.get("@steps").eq(1).should("not.have.class", "c-stepper__step--active");
    cy.get("@steps").eq(2).should("not.have.class", "c-stepper__step--active");

    cy.contains(".card", "Import via BibTeX file")
      .contains(".btn", "Add")
      .click();

    // Upload BibTeX file
    cy.get(".c-file-upload").should(
      "contain.text",
      "Drag and drop your .bib file or",
    );
    cy.contains(".btn", "Upload .bib file")
      .get(".spinner-border")
      .should("not.be.visible");
    cy.get("input[name=file]").selectFile(
      "cypress/fixtures/import-from-bibtex.bib",
    );
    cy.contains(".btn", "Upload .bib file")
      .get(".spinner-border")
      .should("be.visible");

    // Review and publish
    cy.contains("Step 2").should("be.visible");
    cy.contains(".bc-toolbar-title", "Review and publish").should("be.visible");
    cy.get("@steps").eq(0).should("have.class", "c-stepper__step--complete");
    cy.get("@steps").eq(1).should("have.class", "c-stepper__step--active");
    cy.get("@steps").eq(2).should("not.have.class", "c-stepper__step--active");

    cy.contains("Review and publish").should("be.visible");
    cy.wait(1000); // Give elastic some extra time to index imports
    cy.reload();
    cy.contains(".card-header", "Imported publications").should(
      "contain",
      "Showing 1",
    );

    cy.extractBiblioId();

    // Try publishing remaining publication and verify validation error
    cy.ensureNoModal();

    cy.contains(".btn", "Publish all to Biblio").click();

    cy.ensureModal(
      "Unable to publish a publication due to the following errors",
    ).within(() => {
      cy.get(".alert.alert-danger")
        .should("be.visible")
        .should("contain.text", "At least one UGent author is required");

      cy.contains(".btn", "Close").click();
    });

    // Add internal UGent author
    cy.ensureNoModal();

    cy.contains("People & Affiliations").click();

    cy.ensureNoModal();

    cy.contains(".btn", "Add author").click();

    cy.ensureModal("Add author").within(function () {
      cy.intercept({
        pathname: `/publication/${this.biblioId}/contributors/author/suggestions`,
        query: {
          first_name: "John",
          last_name: /^(|Doe)$/, // This forces an exact string match. Just '' matches any string.
        },
      }).as("user-search");

      cy.contains("Search author").should("be.visible");

      cy.setFieldByLabel("First name", "John");
      cy.wait("@user-search");

      cy.setFieldByLabel("Last name", "Doe");
      cy.wait("@user-search");

      cy.contains(".list-group-item", "800000000001")
        // Make sure the right author is selected
        .should("contain.text", "John Doe")
        .should("contain.text", "Active UGent member")
        .should("contain.text", "0000-0003-4217-153X") // ORCID
        .contains(".btn", "Add author")
        .click();
    });

    cy.ensureModal("Add author")
      .within(() => {
        cy.contains("h3", "Review author information").should("be.visible");

        cy.get(".list-group-item")
          .should("have.length", 1)
          .should("contain.text", "John Doe");
      })
      .closeModal(true);

    cy.ensureNoModal();

    cy.contains(".btn", 'Back to "Review and publish" overview').click();

    // Verify publication is still draft
    cy.get(".list-group-item .badge")
      .should("have.class", "badge-warning-light")
      .find(".badge-text")
      .should("have.text", "Biblio draft");

    // Publish
    cy.intercept("POST", "/add-publication/import/multiple/*/publish").as(
      "publish",
    );
    cy.contains(".btn", "Publish all to Biblio").click();
    cy.wait("@publish");

    // Finished
    cy.contains(".bc-toolbar-title", "Congratulations!").should("be.visible");
    cy.get("@steps").eq(0).should("have.class", "c-stepper__step--complete");
    cy.get("@steps").eq(1).should("have.class", "c-stepper__step--complete");
    cy.get("@steps").eq(2).should("have.class", "c-stepper__step--active");
    cy.contains(".card-header", "Imported publications").should(
      "contain",
      "Showing 1",
    );

    cy.contains(".btn", "Continue to overview").click();
    cy.location("pathname").should("eq", "/publication");

    // Verify publication is published
    cy.visitPublication();

    cy.get("#summary .badge")
      .should("have.class", "badge-success-light")
      .find(".badge-text")
      .should("have.text", "Biblio public");
  });

  it("should show an error toast if the import file contains an error", () => {
    cy.loginAsResearcher("researcher1");

    cy.visit("/add-publication");

    cy.contains(".card", "Import via BibTeX file")
      .contains(".btn", "Add")
      .click();

    cy.get("input[name=file]").selectFile(
      "cypress/fixtures/import-from-bibtex-error.bib",
    );

    cy.ensureToast(
      "Sorry, something went wrong. Could not import the publication(s).",
    );
  });

  // TODO: Not yet implemented
  it("should report errors after import");
});
