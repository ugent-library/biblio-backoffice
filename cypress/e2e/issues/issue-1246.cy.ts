// https://github.com/ugent-library/biblio-backoffice/issues/1246

describe("Issue #1246: Close button on toast does not work", () => {
  beforeEach(() => {
    cy.loginAsResearcher("researcher1");
  });

  it("should be possible to dismiss the delete publication toast", () => {
    cy.setUpPublication();
    cy.visitPublication();

    cy.get(".bc-toolbar")
      // The "..." dropdown toggle button
      .find(".dropdown .btn:has(i.if.if-more)")
      .click();
    cy.contains(".dropdown-item", "Delete").click();

    cy.ensureModal("Confirm deletion").closeModal("Delete");
    cy.ensureNoModal();

    assertToast("Publication was successfully deleted.");
  });

  it("should be possible to dismiss the publish publication toast", () => {
    cy.setUpPublication("Miscellaneous", { prepareForPublishing: true });
    cy.visitPublication();

    cy.contains(".btn", "Publish to Biblio").click();
    cy.ensureModal("Are you sure?").closeModal("Publish");
    cy.ensureNoModal();

    assertToast("Publication was successfully published.");
  });

  it("should be possible to dismiss the withdraw publication toast", () => {
    cy.setUpPublication("Miscellaneous", { prepareForPublishing: true });
    cy.visitPublication();

    cy.contains(".btn", "Publish to Biblio").click();
    cy.ensureModal("Are you sure?").closeModal("Publish");
    cy.closeToast();

    cy.contains(".btn", "Withdraw").click();
    cy.ensureModal("Are you sure?").closeModal("Withdraw");
    cy.ensureNoModal();

    assertToast("Publication was successfully withdrawn.");
  });

  it("should be possible to dismiss the republish publication toast", () => {
    cy.setUpPublication("Miscellaneous", { prepareForPublishing: true });
    cy.visitPublication();

    cy.contains(".btn", "Publish to Biblio").click();
    cy.ensureModal("Are you sure?").closeModal("Publish");
    cy.closeToast();

    cy.contains(".btn", "Withdraw").click();
    cy.ensureModal("Are you sure?").closeModal("Withdraw");
    cy.closeToast();

    // Make sure withdraw-toast is gone first
    cy.ensureNoToast();

    cy.contains(".btn", "Republish").click();
    cy.ensureModal("Are you sure?").closeModal("Republish");
    cy.ensureNoModal();

    assertToast("Publication was successfully republished.");
  });

  it("should be possible to dismiss the locked publication toast", () => {
    cy.setUpPublication("Miscellaneous");

    cy.loginAsLibrarian("librarian1");
    cy.visitPublication();

    cy.contains(".btn", "Lock record").click();

    assertToast("Publication was successfully locked.");
  });

  it("should be possible to dismiss the unlocked publication toast", () => {
    cy.setUpPublication("Miscellaneous");

    cy.loginAsLibrarian("librarian1");
    cy.visitPublication();

    cy.contains(".btn", "Lock record").click();
    cy.closeToast();

    // Make sure lock-toast is gone first
    cy.ensureNoToast();

    cy.contains(".btn", "Unlock record").click();

    assertToast("Publication was successfully unlocked.");
  });

  function assertToast(toastMessage: string) {
    cy.ensureToast(toastMessage).closeToast();

    // Reduced assertion timeout here so the test still works if someone decides to reduce the
    // toast dismissal timeout in the future.
    cy.ensureNoToast({ timeout: 500 });
  }
});
