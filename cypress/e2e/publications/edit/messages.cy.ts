describe("Editing publication messages", () => {
  describe("as researcher", () => {
    beforeEach(() => {
      cy.loginAsResearcher("researcher1");

      cy.setUpPublication();
      cy.visitPublication();

      cy.get('.nav-link[title="Biblio messages"]').as("nav").click();
    });

    it("should be possible to add and edit Biblio message", () => {
      cy.get("textarea[name=message]")
        .as("message")
        .setField("initial message")
        .closest(".form-group")
        .contains(".btn", "Save")
        .click();

      cy.reload();
      cy.get("@nav").click();
      cy.get("@message")
        .should("have.value", "initial message")
        .setField("updated message")
        .closest(".form-group")
        .contains(".btn", "Save")
        .click();

      cy.reload();
      cy.get("@nav").click();
      cy.get("@message").should("have.value", "updated message");
    });
  });

  describe("as librarian", () => {
    beforeEach(() => {
      cy.loginAsLibrarian("librarian1");

      cy.setUpPublication();
      cy.visitPublication();

      cy.get('.nav-link[title="Biblio messages"]').as("nav").click();
    });

    it("should be possible to add and edit librarian tags", () => {
      // Edit
      cy.get(
        ".tags:has(textarea#reviewer_tags) tags span.tagify__input[contenteditable]",
      )
        .as("reviewer_tags")
        .setField(
          "initial tag 1{enter}initial tag 2{enter}initial tag 3{enter}",
        );

      // Save
      cy.get("@reviewer_tags")
        .closest(".form-group")
        .contains(".btn", "Save")
        .click();

      // Reload
      cy.reload();
      cy.get("@nav").click();

      // Verify
      cy.get("@reviewer_tags")
        .closest("tags")
        .find("tag")
        .map("innerText")
        .should("have.ordered.members", [
          "initial tag 1",
          "initial tag 2",
          "initial tag 3",
        ]);

      // Edit
      cy.get("@reviewer_tags")
        .setField("updated tag 4{enter}")
        .closest("tags")
        .contains("tag", "initial tag 2")
        .find("x")
        .click();
      cy.get("@reviewer_tags").setField("updated tag 5{enter}");

      // Save
      cy.get("@reviewer_tags")
        .closest(".form-group")
        .contains(".btn", "Save")
        .click();

      // Reload
      cy.reload();
      cy.get("@nav").click();

      // Verify
      cy.get("@reviewer_tags")
        .closest("tags")
        .find("tag")
        .map("innerText")
        .should("have.ordered.members", [
          "initial tag 1",
          "initial tag 3",
          "updated tag 4",
          "updated tag 5",
        ]);
    });

    it("should be possible to add and edit librarian notes", () => {
      cy.get('textarea[name="reviewer_note"]')
        .as("reviewer_note")
        .setField("initial note");

      cy.get("@reviewer_note")
        .closest(".form-group")
        .contains(".btn", "Save")
        .click();

      cy.reload();
      cy.get("@nav").click();

      cy.get("#reviewer_note")
        .should("have.value", "initial note")
        .setField("updated note");

      cy.get("@reviewer_note")
        .closest(".form-group")
        .contains(".btn", "Save")
        .click();

      cy.reload();
      cy.get("@nav").click();

      cy.get("#reviewer_note").should("have.value", "updated note");
    });
  });
});
