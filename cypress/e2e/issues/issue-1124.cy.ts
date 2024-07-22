// https://github.com/ugent-library/biblio-backoffice/issues/1124

describe("Issue #1124:  Add friendlier consistent confirmation toaster when locking or unlocking a record", () => {
  it("should display a toast and badge when publications are locked/unlocked", () => {
    cy.login("researcher1");

    cy.setUpPublication("Dissertation");

    cy.get<string>("@biblioId").then((biblioId) => {
      // Publication does not have lock icon in the list view
      cy.visit("/publication");
      cy.get('input[name=q][placeholder="Search..."]').type(biblioId);
      cy.contains(".btn", "Search").click();
      cy.contains("Showing 1 publications").should("be.visible");
      cy.get(".list-group-item .if.if-lock").should("not.exist");
      cy.contains(".list-group-item", "Locked").should("not.exist");

      // Publication does not have lock icon in the detail view
      cy.visitPublication();
      cy.get(".c-subline .if.if-lock").should("not.exist");
      cy.contains(".c-subline", "Locked").should("not.exist");

      // Edit button is still available for researchers for unlocked publications
      cy.contains(".btn", "Edit").should("be.visible").click();
      cy.ensureModal("Edit publication details");

      // Publication cannot be locked by researcher
      cy.contains(".btn", "Lock record").should("not.exist");

      // Publication can be locked by librarian
      cy.login("librarian1");
      cy.switchMode("Librarian");

      cy.visit("/publication");
      cy.get('input[name=q][placeholder="Search..."]').type(biblioId);
      cy.contains(".btn", "Search").click();
      cy.contains("Showing 1 publications").should("be.visible");
      cy.get(".list-group-item .if.if-lock").should("not.exist");
      cy.contains(".list-group-item", "Locked").should("not.exist");

      cy.visitPublication();
      cy.clocked(() => {
        cy.contains(".btn", "Lock record").should("be.visible").click();

        // Confirmation toast is displayed upon locking (and hidden after 5s)
        cy.contains(".toast", "Publication was successfully locked.")
          .as("lockedToast")
          .should("be.visible");

        cy.tick(5000);

        cy.get("@lockedToast").should("not.exist");
      });

      // Publication now has lock icon in the detail view
      cy.get(".c-subline .if.if-lock").should("be.visible");
      cy.contains(".c-subline", "Locked").should("be.visible");

      // Edit button is still available for librarians for locked publications
      cy.contains(".btn", "Edit").should("be.visible").click();
      cy.ensureModal("Edit publication details");

      // Publication does have lock icon in the list view for librarians
      cy.visit("/publication", { qs: { q: biblioId } });
      cy.contains("Showing 1 publications").should("be.visible");
      cy.get(".list-group-item .if.if-lock").should("be.visible");
      cy.contains(".list-group-item", "Locked").should("be.visible");

      // Publication does have lock icon in the list view for researcher
      cy.login("researcher1");
      cy.visit("/publication", { qs: { q: biblioId, "f[scope]": "all" } });
      cy.contains("Showing 1 publications").should("be.visible");
      cy.get(".list-group-item .if.if-lock").should("be.visible");
      cy.contains(".list-group-item", "Locked").should("be.visible");

      // In the detail view, as a researcher...
      cy.visitPublication();

      // Publication now has a lock icon
      cy.get(".c-subline .if.if-lock").should("be.visible");
      cy.contains(".c-subline", "Locked").should("be.visible");

      // Edit button is no longer available
      cy.contains(".btn", "Edit").should("not.exist");

      // Now go back as librarian to unlock the publication
      cy.login("librarian1");
      cy.switchMode("Librarian");

      cy.visitPublication();

      cy.clocked(() => {
        cy.contains(".btn", "Unlock record").should("be.visible").click();

        // Confirmation toast is displayed upon locking (and hidden after 5s)
        cy.contains(".toast", "Publication was successfully unlocked.")
          .as("unlockedToast")
          .should("be.visible");

        cy.tick(5000);

        cy.get("@unlockedToast").should("not.exist");
      });

      // Lock icon is removed again
      cy.get(".c-subline .if.if-lock").should("not.exist");
      cy.contains(".c-subline", "Locked").should("not.exist");
    });
  });
});
