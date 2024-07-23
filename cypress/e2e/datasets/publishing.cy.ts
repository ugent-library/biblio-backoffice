describe("Publishing datasets", () => {
  describe("as researcher", () => {
    beforeEach(() => {
      cy.login("researcher1");
    });

    it("should be possible to publish, withdraw and republish a dataset", () => {
      cy.setUpDataset({ prepareForPublishing: true });
      cy.visitDataset();

      cy.contains(".btn-success", "Publish to Biblio").click();
      cy.ensureModal("Are you sure?")
        .within(() => {
          cy.get(".modal-body").should(
            "contain",
            "Are you sure you want to publish this dataset to Biblio?",
          );
        })
        .closeModal("Publish");
      cy.ensureNoModal();
      cy.ensureToast("Dataset was successfully published.").closeToast();

      cy.contains(".btn-outline-danger", "Withdraw").click();
      cy.ensureModal("Are you sure?")
        .within(() => {
          cy.get(".modal-body").should(
            "contain",
            "Are you sure you want to withdraw this dataset from Biblio?",
          );
        })
        .closeModal("Withdraw");
      cy.ensureNoModal();
      cy.ensureToast("Dataset was successfully withdrawn.").closeToast();

      cy.contains(".btn-success", "Republish to Biblio").click();
      cy.ensureModal("Are you sure?")
        .within(() => {
          cy.get(".modal-body").should(
            "contain",
            "Are you sure you want to republish this dataset to Biblio?",
          );
        })
        .closeModal("Republish");
      cy.ensureNoModal();
      cy.ensureToast("Dataset was successfully republished.").closeToast();
    });

    it("should error when dataset is not ready for publication", () => {
      cy.setUpDataset({ prepareForPublishing: false });
      cy.visitDataset();

      cy.contains(".btn-success", "Publish to Biblio").click();
      cy.ensureModal("Are you sure?")
        .within(() => {
          cy.get(".modal-body").should(
            "contain",
            "Are you sure you want to publish this dataset to Biblio?",
          );
        })
        .closeModal("Publish");

      cy.ensureModal(
        "Unable to publish this dataset due to the following errors",
      )
        .within(() => {
          cy.get("ul > li")
            .map("textContent")
            .should("have.members", [
              "Access level is required",
              "Format is required",
              "Publisher is required",
              "Publication year is required",
              "One or more authors are required",
              "At least one UGent author is required",
              "License is required",
            ]);
        })
        .closeModal("Close");
      cy.ensureNoModal();
      cy.ensureNoToast();

      cy.reload();

      cy.contains(".btn-success", "Publish to Biblio").should("be.visible");
    });

    it("should error when dataset is not ready for republication", () => {
      cy.setUpDataset({ prepareForPublishing: true });
      cy.visitDataset();

      cy.contains(".btn-success", "Publish to Biblio").click();
      cy.ensureModal("Are you sure?").closeModal("Publish");
      cy.ensureNoModal();
      cy.ensureToast("Dataset was successfully published.").closeToast();

      cy.contains(".btn-outline-danger", "Withdraw").click();
      cy.ensureModal("Are you sure?").closeModal("Withdraw");
      cy.ensureNoModal();
      cy.ensureToast("Dataset was successfully withdrawn.").closeToast();

      cy.updateFields(
        "Dataset details",
        () => {
          cy.setFieldByLabel("Publisher", "");
          cy.setFieldByLabel("Publication year", "");
        },
        true,
      );

      cy.contains(".nav-link", "People & Affiliations").click();

      // Delete interal author
      cy.get("#authors button:has(.if-delete)").click();
      cy.ensureModal("Confirm deletion").closeModal("Delete");

      // Add external author
      cy.addCreator("John", "Doe", { external: true });

      cy.contains(".btn-success", "Republish to Biblio").click();
      cy.ensureModal("Are you sure?")
        .within(() => {
          cy.get(".modal-body").should(
            "contain",
            "Are you sure you want to republish this dataset to Biblio?",
          );
        })
        .closeModal("Republish");

      cy.ensureModal(
        "Unable to republish this dataset due to the following errors",
      )
        .within(() => {
          cy.get("ul > li")
            .map("textContent")
            .should("have.members", [
              "Publisher is required",
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
      cy.login("librarian1");
    });

    it("should be possible to lock and unlock a dataset", () => {
      cy.setUpDataset();
      cy.visitDataset();

      cy.contains(".btn-outline-secondary", "Lock").click();
      cy.ensureNoModal();
      cy.ensureToast("Dataset was successfully locked.").closeToast();

      cy.contains(".btn-outline-secondary", "Unlock").click();
      cy.ensureNoModal();
      cy.ensureToast("Dataset was successfully unlocked.").closeToast();
    });
  });
});
