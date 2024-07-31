describe("Deleting a publication", () => {
  it("should be possible to delete a publication", () => {
    cy.login("researcher1");

    cy.setUpPublication();
    cy.visitPublication();

    cy.get(".btn .if-more").click();
    cy.contains(".dropdown-item", "Delete").click();

    cy.ensureModal("Confirm deletion")
      .within(() => {
        cy.get(".modal-body").should(
          "contain",
          "Are you sure you want to delete this publication?",
        );
      })
      .closeModal("Delete");
    cy.ensureNoModal();

    cy.location("pathname").should("eq", "/publication");

    cy.get<string>("@biblioId").then((biblioId) => {
      cy.request({
        url: `/publication/${biblioId}`,
        failOnStatusCode: false,
      }).should("have.property", "isOkStatusCode", false);
    });
  });
});
