describe("Media type suggestions", () => {
  it("should provide format type suggestions", () => {
    cy.loginAsResearcher("researcher1");

    cy.visit("media_type/suggestions", {
      qs: { input: "format", format: "earth" },
    });

    cy.get(".card .list-group .list-group-item")
      .as("items")
      .should("have.length", 3);

    cy.get("@items")
      .eq(0)
      .should("have.attr", "data-value", "earth")
      .should("contain.text", 'Use custom data format "earth"')
      .find(".badge")
      .should("not.exist");

    cy.get("@items")
      .eq(1)
      .should("have.attr", "data-value", "application/vnd.google-earth.kmz")
      .find(".badge")
      .should("contains.text", "application/vnd.google-earth.kmz")
      .parent()
      .prop("innerText")
      .should("contain", "application/vnd.google-earth.kmz")
      .should("not.contain", "(")
      .should("not.contain", ")");

    cy.get("@items")
      .eq(2)
      .should("have.attr", "data-value", "application/vnd.google-earth.kml+xml")
      .find(".badge")
      .should("contains.text", "application/vnd.google-earth.kml+xml")
      .parent()
      .prop("innerText")
      .should("contain", "application/vnd.google-earth.kml+xml (.kml)");
  });
});
