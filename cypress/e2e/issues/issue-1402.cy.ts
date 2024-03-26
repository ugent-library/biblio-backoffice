// https://github.com/ugent-library/biblio-backoffice/issues/1402

import type { PublicationType } from "support/commands/set-up-publication";

describe("Issue #1402: Gohtml conversion to Templ", () => {
  beforeEach(() => {
    cy.loginAsLibrarian();
  });

  describe("for publications", () => {
    type Test = {
      name: string;
      publicationType?: PublicationType;
      selector: string;
      fields: { label: string; initial: string; updated: string }[];
    };

    const TESTS: Record<string, Test> = {
      // TODO: project

      abstracts: {
        name: "abstract",
        selector: "#abstracts",
        fields: [
          {
            label: "Abstract",
            initial: "The initial abstract",
            updated: "The updated abstract",
          },
          {
            label: "Language",
            initial: "Danish",
            updated: "Northern Sami",
          },
        ],
      },

      links: {
        name: "link",
        selector: "#links",
        fields: [
          {
            label: "URL",
            initial: "https://www.ugent.be",
            updated: "https://lib.ugent.be",
          },
          {
            label: "Relation",
            initial: "Related information",
            updated: "Accompanying website",
          },
          {
            label: "Description",
            initial: "The initial website",
            updated: "The updated website",
          },
        ],
      },

      "lay summaries": {
        name: "lay summary",
        publicationType: "Dissertation",
        selector: "#lay-summaries",
        fields: [
          {
            label: "Lay summary",
            initial: "The initial lay summary",
            updated: "The updated lay summary",
          },
          {
            label: "Language",
            initial: "Italian",
            updated: "Multiple languages",
          },
        ],
      },
    };

    Object.entries(TESTS).forEach(
      ([type, { name, publicationType, selector, fields }]) => {
        it(`should be possible to add, edit and delete ${type}`, () => {
          cy.setUpPublication(publicationType);
          cy.visitPublication();

          cy.get(selector).find("table tbody tr").should("have.length", 0);

          cy.contains(".btn", `Add ${name}`).click();
          cy.ensureModal(`Add ${name}`)
            .within(() => {
              for (const field of fields) {
                cy.setFieldByLabel(field.label, field.initial);
              }
            })
            .closeModal(`Add ${name}`);
          cy.ensureNoModal();

          cy.get(selector)
            .find("table tbody tr")
            .as("row")
            .should("have.length", 1);

          for (const field of fields) {
            cy.get("@row").should("contain", field.initial);
          }

          cy.get("@row").within(() => {
            cy.get(".if-more").click();

            cy.contains("button", "Edit").click();
          });

          cy.ensureModal(`Edit ${name}`)
            .within(() => {
              for (const field of fields) {
                cy.setFieldByLabel(field.label, field.updated);
              }
            })
            .closeModal(`Update ${name}`);
          cy.ensureNoModal();

          cy.get("@row").should("have.length", 1);

          for (const field of fields) {
            cy.get("@row").should("not.contain", field.initial);
            cy.get("@row").should("contain", field.updated);
          }

          cy.get("@row").within(() => {
            cy.get(".if-more").click();

            cy.contains("button", "Delete").click();
          });

          cy.ensureModal("Are you sure").closeModal("Delete");
          cy.ensureNoModal();

          cy.get(selector).find("table tbody tr").should("have.length", 0);
        });
      },
    );

    it("should be possible to add and edit conference details", () => {
      cy.setUpPublication("Conference contribution");
      cy.visitPublication();

      cy.get("#conference-details").contains(".btn", "Edit").click();
      cy.ensureModal("Edit conference details")
        .within(() => {
          cy.setFieldByLabel("Conference", "The conference name");
          cy.setFieldByLabel("Conference location", "The conference location");
          cy.setFieldByLabel(
            "Conference organiser",
            "The conference organiser",
          );
          cy.setFieldByLabel("Conference start date", "2021-01-01");
          cy.setFieldByLabel("Conference end date", "2022-02-02");
        })
        .closeModal(true);
      cy.ensureNoModal();

      cy.get("#conference-details")
        .should("contain", "The conference name")
        .should("contain", "The conference location")
        .should("contain", "The conference organiser")
        .should("contain", "2021-01-01")
        .should("contain", "2022-02-02");

      cy.get("#conference-details").contains(".btn", "Edit").click();
      cy.ensureModal("Edit conference details")
        .within(() => {
          cy.setFieldByLabel("Conference", "The updated conference name");
          cy.setFieldByLabel(
            "Conference location",
            "The updated conference location",
          );
          cy.setFieldByLabel(
            "Conference organiser",
            "The updated conference organiser",
          );
          cy.setFieldByLabel("Conference start date", "2023-03-03");
          cy.setFieldByLabel("Conference end date", "2024-04-04");
        })
        .closeModal(true);
      cy.ensureNoModal();

      cy.get("#conference-details")
        .should("contain", "The updated conference name")
        .should("contain", "The updated conference location")
        .should("contain", "The updated conference organiser")
        .should("contain", "2023-03-03")
        .should("contain", "2024-04-04");
    });

    it("should be possible to add and edit additional info", () => {
      cy.setUpPublication();
      cy.visitPublication();

      cy.get("#additional-information").contains(".btn", "Edit").click();
      cy.ensureModal("Edit additional information")
        .within(() => {
          cy.setFieldByLabel("Research field", "Performing Arts")
            .next("button.form-value-add")
            .click()
            .closest(".form-value")
            .next(".form-value")
            .find("select")
            .select("Social Sciences");

          cy.getLabel("Keywords")
            .next(".tags")
            .find("span[contenteditable]")
            .type("these{enter}are{enter}the{enter}keywords", { delay: 10 });

          cy.setFieldByLabel(
            "Additional information",
            "The additional information",
          );
        })
        .closeModal(true);
      cy.ensureNoModal();

      cy.get("#additional-information")
        .should("contain", "Performing Arts")
        .should("contain", "Social Sciences")
        .should("contain", "The additional information")
        .find(".badge")
        .should("have.length", 4)
        .map("textContent")
        .should("have.ordered.members", ["these", "are", "the", "keywords"]);

      cy.get("#additional-information").contains(".btn", "Edit").click();
      cy.ensureModal("Edit additional information")
        .within(() => {
          cy.getLabel("Research field")
            .next(".form-values")
            .as("formValues")
            .contains("select", "Performing Arts")
            .next("button:contains(Delete)")
            .click();

          cy.get("@formValues")
            .find(".form-value")
            .last()
            .find("select")
            .select("Chemistry");

          cy.get("tags").contains("tag", "these").find("x").click();
          cy.get("tags").contains("tag", "are").find("x").click();
          cy.get("tags").find("span[contenteditable]").type("updated");

          cy.setFieldByLabel(
            "Additional information",
            "The updated information",
          );
        })
        .closeModal(true);
      cy.ensureNoModal();

      cy.get("#additional-information")
        .should("contain", "Social Sciences")
        .should("contain", "Chemistry")
        .should("not.contain", "Performing Arts")
        .should("contain", "The updated information")
        .should("not.contain", "The additional information")
        .find(".badge")
        .should("have.length", 3)
        .map("textContent")
        .should("have.ordered.members", ["the", "keywords", "updated"]);
    });

    it("should be possible to add and delete departments", () => {
      cy.setUpPublication();
      cy.visitPublication();

      cy.contains(".nav .nav-item", "People & Affiliations").click();

      cy.get("#departments").contains(".btn", "Add department").click();

      cy.ensureModal("Select departments").within(() => {
        cy.getLabel("Search").next("input").type("LW");

        cy.contains(".btn", "Add department").click();
      });
      cy.ensureNoModal();

      cy.get("#departments-body .list-group-item-text h4")
        .map("textContent")
        .should("have.ordered.members", ["Faculty of Arts and Philosophy"]);

      cy.get("#departments").contains(".btn", "Add department").click();

      cy.ensureModal("Select departments").within(() => {
        cy.getLabel("Search").next("input").type("DI");

        cy.contains(".btn", "Add department").click();
      });
      cy.ensureNoModal();

      cy.get("#departments-body .list-group-item-text h4")
        .map("textContent")
        .should("have.ordered.members", [
          "Faculty of Arts and Philosophy",
          "Faculty of Veterinary Medicine",
        ]);

      cy.contains("#departments-body tr", "Faculty of Arts and Philosophy")
        .find(".if-more")
        .click();
      cy.contains(".dropdown-item", "Remove from publication").click();

      cy.ensureModal("Are you sure").closeModal("Delete");
      cy.ensureNoModal();

      cy.get("#departments-body .list-group-item-text h4")
        .map("textContent")
        .should("have.ordered.members", ["Faculty of Veterinary Medicine"]);
    });

    it("should be possible to edit librarian tags", () => {
      cy.setUpPublication();
      cy.visitPublication();

      cy.contains(".nav .nav-item", "Biblio Messages").click();

      cy.contains(".card", "Librarian tags").contains(".btn", "Edit").click();

      cy.ensureModal(/^Edit Librarian tags/)
        .within(() => {
          cy.setFieldByLabel("Librarian tags", "initial tag 1");
          cy.contains("button", "Add").click();
          cy.get("input").last().type("initial tag 2");
          cy.contains("button", "Add").click();
          cy.get("input").last().type("initial tag 3");
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
          cy.get("input").last().type("updated tag 4");

          cy.get("input")
            .eq(1)
            .should("have.value", "initial tag 2")
            .next("button:contains(Delete)")
            .click();

          cy.contains("button", "Add").click();
          cy.get("input").last().type("updated tag 5");
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

    it("should be possible to edit librarian notes", () => {
      cy.setUpPublication();
      cy.visitPublication();

      cy.contains(".nav .nav-item", "Biblio Messages").click();

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

    it("should be possible to edit Biblio message", () => {
      cy.setUpPublication();
      cy.visitPublication();

      cy.contains(".nav .nav-item", "Biblio Messages").click();

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
  });
});
