describe("Publishing publications", () => {
  describe("as researcher", () => {
    beforeEach(() => {
      cy.loginAsResearcher();
    });

    it("should be possible to publish, withdraw and republish a publication", () => {
      cy.setUpPublication(undefined, { prepareForPublishing: true });
      cy.visitPublication();

      cy.contains(".btn-success", "Publish to Biblio").click();
      cy.ensureModal("Are you sure?")
        .within(() => {
          cy.get(".modal-body").should(
            "contain",
            "Are you sure you want to publish this publication to Biblio?",
          );
        })
        .closeModal("Publish");
      cy.ensureNoModal();
      cy.ensureToast("Publication was successfully published.").closeToast();

      cy.contains(".btn-outline-danger", "Withdraw").click();
      cy.ensureModal("Are you sure?")
        .within(() => {
          cy.get(".modal-body").should(
            "contain",
            "Are you sure you want to withdraw this publication from Biblio?",
          );
        })
        .closeModal("Withdraw");
      cy.ensureNoModal();
      cy.ensureToast("Publication was successfully withdrawn.").closeToast();

      cy.contains(".btn-success", "Republish to Biblio").click();
      cy.ensureModal("Are you sure?")
        .within(() => {
          cy.get(".modal-body").should(
            "contain",
            "Are you sure you want to republish this publication to Biblio?",
          );
        })
        .closeModal("Republish");
      cy.ensureNoModal();
      cy.ensureToast("Publication was successfully republished.").closeToast();
    });

    it("should error when publication is not ready for publication", () => {
      cy.setUpPublication("Miscellaneous", { prepareForPublishing: false });
      cy.visitPublication();

      cy.contains(".btn-success", "Publish to Biblio").click();
      cy.ensureModal("Are you sure?")
        .within(() => {
          cy.get(".modal-body").should(
            "contain",
            "Are you sure you want to publish this publication to Biblio?",
          );
        })
        .closeModal("Publish");

      cy.ensureModal(
        "Unable to publish this publication due to the following errors",
      )
        .within(() => {
          cy.get("ul > li")
            .map("textContent")
            .should("have.members", [
              "Publication year is required",
              "One or more authors are required",
              "At least one UGent author is required",
            ]);
        })
        .closeModal("Close");
      cy.ensureNoModal();
      cy.ensureNoToast();

      cy.reload();

      cy.contains(".btn-success", "Publish to Biblio").should("be.visible");
    });

    it("should error when publication is not ready for republication", () => {
      cy.setUpPublication("Miscellaneous", { prepareForPublishing: true });
      cy.visitPublication();

      cy.contains(".btn-success", "Publish to Biblio").click();
      cy.ensureModal("Are you sure?").closeModal("Publish");
      cy.ensureNoModal();
      cy.ensureToast("Publication was successfully published.").closeToast();

      cy.contains(".btn-outline-danger", "Withdraw").click();
      cy.ensureModal("Are you sure?").closeModal("Withdraw");
      cy.ensureNoModal();
      cy.ensureToast("Publication was successfully withdrawn.").closeToast();

      cy.updateFields(
        "Publication details",
        () => {
          cy.setFieldByLabel("Publication year", "");
        },
        true,
      );

      cy.contains(".nav-link", "People & Affiliations").click();

      // Delete interal author
      cy.get("#authors button:has(.if-delete)").click();
      cy.ensureModal("Confirm deletion").closeModal("Delete");

      // Add external author
      cy.addAuthor("John", "Doe", true);

      cy.contains(".btn-success", "Republish to Biblio").click();
      cy.ensureModal("Are you sure?")
        .within(() => {
          cy.get(".modal-body").should(
            "contain",
            "Are you sure you want to republish this publication to Biblio?",
          );
        })
        .closeModal("Republish");

      cy.ensureModal(
        "Unable to republish this publication due to the following errors",
      )
        .within(() => {
          cy.get("ul > li")
            .map("textContent")
            .should("have.members", [
              "Publication year is required",
              "At least one UGent author is required",
            ]);
        })
        .closeModal("Close");
      cy.ensureNoModal();
      cy.ensureNoToast();

      cy.reload();

      cy.contains(".btn-success", "Republish to Biblio").should("be.visible");
    });
  });

  describe("as librarian", () => {
    beforeEach(() => {
      cy.loginAsLibrarian();
    });

    it("should be possible to lock and unlock a publication", () => {
      cy.setUpPublication();
      cy.visitPublication();

      cy.contains(".btn-outline-secondary", "Lock").click();
      cy.ensureNoModal();
      cy.ensureToast("Publication was successfully locked.").closeToast();

      cy.contains(".btn-outline-secondary", "Unlock").click();
      cy.ensureNoModal();
      cy.ensureToast("Publication was successfully unlocked.").closeToast();
    });
  });
});
