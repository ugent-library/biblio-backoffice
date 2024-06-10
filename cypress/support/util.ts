export function getRandomText() {
  return crypto.randomUUID().replace(/-/g, "").toUpperCase();
}

export function testFocusForLabel(
  labelText: string,
  fieldSelector: string,
  autoFocus = false,
) {
  cy.getLabel(labelText)
    .as("theLabel")
    .should("have.length", 1)
    .parent({ log: false })
    .find(fieldSelector)
    .should("have.length", 1)
    .first({ log: false })
    .as("theField")
    .should(autoFocus ? "be.focused" : "not.be.focused");

  if (autoFocus) {
    cy.focused().blur();

    cy.get("@theField").should("not.be.focused");
  }

  cy.get("@theLabel").click();

  cy.get("@theField").should("be.focused");
}
