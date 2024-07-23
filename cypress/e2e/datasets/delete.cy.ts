describe("Deleting a dataset", () => {
  it("should be possible to delete a dataset", () => {
    cy.login("researcher1");

    cy.setUpDataset();
    cy.visitDataset();

    cy.get(".btn .if-more").click();
    cy.contains(".dropdown-item", "Delete").click();

    cy.ensureModal("Confirm deletion")
      .within(() => {
        cy.get(".modal-body").should(
          "contain",
          "Are you sure you want to delete this dataset?",
        );
      })
      .closeModal("Delete");
    cy.ensureNoModal();

    cy.location("pathname").should("eq", "/dataset");

    cy.get<string>("@biblioId").then((biblioId) => {
      cy.request({
        url: `/dataset/${biblioId}`,
        failOnStatusCode: false,
      }).should("have.property", "isOkStatusCode", false);
    });
  });
});
