import { testFocusForLabel } from "support/util";

describe("Editing dataset Biblio messages", () => {
  describe("as researcher", () => {
    beforeEach(() => {
      cy.loginAsResearcher();

      cy.setUpDataset();
      cy.visitDataset();

      cy.contains(".nav .nav-item", "Biblio Messages").click();
    });

    it("should be possible to add and edit Biblio message", () => {
      cy.contains(".card", "Messages from and for Biblio team")
        .contains(".btn", "Edit")
        .click();

      cy.ensureModal("Edit messages from and for Biblio team")
        .within(() => {
          cy.setFieldByLabel("Message", "initial message");
        })
        .closeModal(true);
      cy.ensureNoModal();

      cy.get("#message-body").should("contain", "initial message");

      cy.contains(".card", "Messages from and for Biblio team")
        .contains(".btn", "Edit")
        .click();

      cy.ensureModal("Edit messages from and for Biblio team")
        .within(() => {
          cy.setFieldByLabel("Message", "updated message");
        })
        .closeModal(true);
      cy.ensureNoModal();

      cy.get("#message-body").should("contain", "updated message");
    });

    it("should have clickable labels in the edit Biblio message dialog", () => {
      cy.updateFields("Messages from and for Biblio team", () => {
        testFocusForLabel("Message", 'textarea[name="message"]');
      });
    });
  });

  describe("as librarian", () => {
    beforeEach(() => {
      cy.loginAsLibrarian();

      cy.setUpDataset();
      cy.visitDataset();

      cy.contains(".nav .nav-item", "Biblio Messages").click();
    });

    it("should be possible to add and edit librarian tags", () => {
      cy.contains(".card", "Librarian tags").contains(".btn", "Edit").click();

      cy.ensureModal(/^Edit Librarian tags/)
        .within(() => {
          cy.getLabel("Librarian tags")
            .next()
            .find("tags span[contenteditable]")
            .type(
              "initial tag 1{enter}initial tag 2{enter}initial tag 3{enter}",
              { delay: 10 },
            );

          // Give Tagify a bit of time to process this
          cy.wait(50);
        })
        .closeModal(true);
      cy.ensureNoModal();

      cy.get("#reviewer-tags-body .badge")
        .map("textContent")
        .should("have.ordered.members", [
          "initial tag 1",
          "initial tag 2",
          "initial tag 3",
        ]);

      cy.contains(".card", "Librarian tags").contains(".btn", "Edit").click();

      cy.ensureModal(/^Edit Librarian tags/)
        .within(() => {
          cy.getLabel("Librarian tags")
            .next()
            .find("tags span[contenteditable]")
            .as("tags")
            .type("updated tag 4{enter}");

          cy.contains("tags tag", "initial tag 2").find("x").click();

          cy.get("@tags").type("updated tag 5{enter}");

          // Give Tagify a bit of time to process this
          cy.wait(50);
        })
        .closeModal(true);
      cy.ensureNoModal();

      cy.get("#reviewer-tags-body .badge")
        .map("textContent")
        .should("have.ordered.members", [
          "initial tag 1",
          "initial tag 3",
          "updated tag 4",
          "updated tag 5",
        ]);
    });

    it("should be possible to add and edit librarian notes", () => {
      cy.contains(".card", "Librarian note").contains(".btn", "Edit").click();

      cy.ensureModal(/^Edit Librarian note/)
        .within(() => {
          cy.setFieldByLabel("Librarian note", "initial note");
        })
        .closeModal(true);
      cy.ensureNoModal();

      cy.get("#reviewer-note-body").should("contain", "initial note");

      cy.contains(".card", "Librarian note").contains(".btn", "Edit").click();

      cy.ensureModal(/^Edit Librarian note/)
        .within(() => {
          cy.setFieldByLabel("Librarian note", "updated note");
        })
        .closeModal(true);
      cy.ensureNoModal();

      cy.get("#reviewer-note-body").should("contain", "updated note");
    });

    it("should have clickable labels in the edit librarian note dialog", () => {
      cy.updateFields("Librarian note", () => {
        testFocusForLabel("Librarian note", 'textarea[name="reviewer_note"]');
      });
    });
  });
});
