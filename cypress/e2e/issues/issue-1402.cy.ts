// https://github.com/ugent-library/biblio-backoffice/issues/1402

describe("Issue #1402: Gohtml conversion to Templ", () => {
  beforeEach(() => {
    cy.loginAsLibrarian();
  });

  describe("for publications", () => {
    // TODO
    it("should be possible to change the publication type");

    it("should be possible to add, edit and delete abstracts", () => {
      cy.setUpPublication();
      cy.visitPublication();

      cy.get("#abstracts").find("table tbody tr").should("have.length", 0);

      cy.contains(".btn", "Add abstract").click();
      cy.ensureModal("Add abstract")
        .within(() => {
          cy.setFieldByLabel("Abstract", "The initial abstract");
          cy.setFieldByLabel("Language", "Danish");
        })
        .closeModal("Add abstract");
      cy.ensureNoModal();

      cy.get("#abstracts")
        .find("table tbody tr")
        .as("row")
        .should("have.length", 1);

      cy.get("@row").should("contain", "The initial abstract");
      cy.get("@row").should("contain", "Danish");

      cy.get("@row").within(() => {
        cy.get(".if-more").click();

        cy.contains("button", "Edit").click();
      });

      cy.ensureModal("Edit abstract")
        .within(() => {
          cy.setFieldByLabel("Abstract", "The updated abstract");
          cy.setFieldByLabel("Language", "Northern Sami");
        })
        .closeModal("Update abstract");
      cy.ensureNoModal();

      cy.get("@row").should("have.length", 1);

      cy.get("@row").should("not.contain", "The initial abstract");
      cy.get("@row").should("contain", "The updated abstract");
      cy.get("@row").should("not.contain", "Danish");
      cy.get("@row").should("contain", "Northern Sami");

      cy.get("@row").within(() => {
        cy.get(".if-more").click();

        cy.contains("button", "Delete").click();
      });

      cy.ensureModal("Are you sure").closeModal("Delete");
      cy.ensureNoModal();

      cy.get("#abstracts").find("table tbody tr").should("have.length", 0);
    });

    it("should be possible to add, edit and delete links", () => {
      cy.setUpPublication();
      cy.visitPublication();

      cy.get("#links").find("table tbody tr").should("have.length", 0);

      cy.contains(".btn", "Add link").click();
      cy.ensureModal("Add link")
        .within(() => {
          cy.setFieldByLabel("URL", "https://www.ugent.be");
          cy.setFieldByLabel("Relation", "Related information");
          cy.setFieldByLabel("Description", "The initial website");
        })
        .closeModal("Add link");
      cy.ensureNoModal();

      cy.get("#links")
        .find("table tbody tr")
        .as("row")
        .should("have.length", 1);

      cy.get("@row").should("contain", "https://www.ugent.be");
      cy.get("@row").should("contain", "Related information");
      cy.get("@row").should("contain", "The initial website");

      cy.get("@row").within(() => {
        cy.get(".if-more").click();

        cy.contains("button", "Edit").click();
      });

      cy.ensureModal("Edit link")
        .within(() => {
          cy.setFieldByLabel("URL", "https://lib.ugent.be");
          cy.setFieldByLabel("Relation", "Accompanying website");
          cy.setFieldByLabel("Description", "The updated website");
        })
        .closeModal("Update link");
      cy.ensureNoModal();

      cy.get("@row").should("have.length", 1);

      cy.get("@row").should("not.contain", "https://www.ugent.be");
      cy.get("@row").should("contain", "https://lib.ugent.be");
      cy.get("@row").should("not.contain", "Related information");
      cy.get("@row").should("contain", "Accompanying website");
      cy.get("@row").should("not.contain", "The initial website");
      cy.get("@row").should("contain", "The updated website");

      cy.get("@row").within(() => {
        cy.get(".if-more").click();

        cy.contains("button", "Delete").click();
      });

      cy.ensureModal("Are you sure").closeModal("Delete");
      cy.ensureNoModal();

      cy.get("#links").find("table tbody tr").should("have.length", 0);
    });

    it("should be possible to add, edit and delete lay summaries", () => {
      cy.setUpPublication("Dissertation");
      cy.visitPublication();

      cy.get("#lay-summaries").find("table tbody tr").should("have.length", 0);

      cy.contains(".btn", "Add lay summary").click();
      cy.ensureModal("Add lay summary")
        .within(() => {
          cy.setFieldByLabel("Lay summary", "The initial lay summary");
          cy.setFieldByLabel("Language", "Italian");
        })
        .closeModal("Add lay summary");
      cy.ensureNoModal();

      cy.get("#lay-summaries")
        .find("table tbody tr")
        .as("row")
        .should("have.length", 1);

      cy.get("@row").should("contain", "The initial lay summary");
      cy.get("@row").should("contain", "Italian");

      cy.get("@row").within(() => {
        cy.get(".if-more").click();

        cy.contains("button", "Edit").click();
      });

      cy.ensureModal("Edit lay summary")
        .within(() => {
          cy.setFieldByLabel("Lay summary", "The updated lay summary");
          cy.setFieldByLabel("Language", "Multiple languages");
        })
        .closeModal("Update lay summary");
      cy.ensureNoModal();

      cy.get("@row").should("have.length", 1);

      cy.get("@row").should("not.contain", "The initial lay summary");
      cy.get("@row").should("contain", "The updated lay summary");
      cy.get("@row").should("not.contain", "Italian");
      cy.get("@row").should("contain", "Multiple languages");

      cy.get("@row").within(() => {
        cy.get(".if-more").click();

        cy.contains("button", "Delete").click();
      });

      cy.ensureModal("Are you sure").closeModal("Delete");
      cy.ensureNoModal();

      cy.get("#lay-summaries").find("table tbody tr").should("have.length", 0);
    });

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

    it("should be possible to add and edit librarian tags", () => {
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

    it("should be possible to add and edit librarian notes", () => {
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

    it("should be possible to add and edit Biblio message", () => {
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

    it("should be possible to publish, withdraw, republish, lock and unlock a publication", () => {
      cy.setUpPublication(undefined, { prepareForPublishing: true });
      cy.visitPublication();

      cy.contains(".btn-success", "Publish to Biblio").click();
      cy.ensureModal("Are you sure?")
        .within(() => {
          cy.get(".modal-content").should(
            "contain",
            "Are you sure you want to publish this publication to Biblio?",
          );
        })
        .closeModal("Publish");
      cy.ensureNoModal();
      cy.ensureToast("Publication was successfully published.").closeToast();

      cy.contains(".btn-outline-danger", "Withdraw").click();
      cy.ensureModal("Are you sure?")
        .within(() => {
          cy.get(".modal-content").should(
            "contain",
            "Are you sure you want to withdraw this publication to Biblio?",
          ); // TODO: "from Biblio"
        })
        .closeModal("Withdraw");
      cy.ensureNoModal();
      cy.ensureToast("Publication was successfully withdrawn.").closeToast();

      cy.contains(".btn-success", "Republish to Biblio").click();
      cy.ensureModal("Are you sure?")
        .within(() => {
          cy.get(".modal-content").should(
            "contain",
            "Are you sure you want to republish this publication to Biblio?",
          );
        })
        .closeModal("Republish");
      cy.ensureNoModal();
      cy.ensureToast("Publication was successfully republished.").closeToast();

      cy.contains(".btn-outline-secondary", "Lock").click();
      cy.ensureNoModal();
      cy.ensureToast("Publication was successfully locked.").closeToast();

      cy.contains(".btn-outline-secondary", "Unlock").click();
      cy.ensureNoModal();
      cy.ensureToast("Publication was successfully unlocked.").closeToast();
    });
  });

  describe("for datasets", () => {
    beforeEach(() => {
      cy.setUpDataset();
      cy.visitDataset();
    });

    it("should be possible to add, edit and delete abstracts", () => {
      cy.get("#abstracts").find("table tbody tr").should("have.length", 0);

      cy.contains(".btn", "Add abstract").click();
      cy.ensureModal("Add abstract")
        .within(() => {
          cy.setFieldByLabel("Abstract", "The initial abstract");
          cy.setFieldByLabel("Language", "Danish");
        })
        .closeModal("Add abstract");
      cy.ensureNoModal();

      cy.get("#abstracts")
        .find("table tbody tr")
        .as("row")
        .should("have.length", 1);

      cy.get("@row").should("contain", "The initial abstract");
      cy.get("@row").should("contain", "Danish");

      cy.get("@row").within(() => {
        cy.get(".if-more").click();

        cy.contains("button", "Edit").click();
      });

      cy.ensureModal("Edit abstract")
        .within(() => {
          cy.setFieldByLabel("Abstract", "The updated abstract");
          cy.setFieldByLabel("Language", "Northern Sami");
        })
        .closeModal("Update abstract");
      cy.ensureNoModal();

      cy.get("@row").should("have.length", 1);

      cy.get("@row").should("not.contain", "The initial abstract");
      cy.get("@row").should("contain", "The updated abstract");
      cy.get("@row").should("not.contain", "Danish");
      cy.get("@row").should("contain", "Northern Sami");

      cy.get("@row").within(() => {
        cy.get(".if-more").click();

        cy.contains("button", "Delete").click();
      });

      cy.ensureModal("Are you sure").closeModal("Delete");
      cy.ensureNoModal();

      cy.get("#abstracts").find("table tbody tr").should("have.length", 0);
    });

    it("should be possible to add, edit and delete links", () => {
      cy.get("#links").find("table tbody tr").should("have.length", 0);

      cy.contains(".btn", "Add link").click();
      cy.ensureModal("Add link")
        .within(() => {
          cy.setFieldByLabel("URL", "https://www.ugent.be");
          cy.setFieldByLabel("Relation", "Related information");
          cy.setFieldByLabel("Description", "The initial website");
        })
        .closeModal("Add link");
      cy.ensureNoModal();

      cy.get("#links")
        .find("table tbody tr")
        .as("row")
        .should("have.length", 1);

      cy.get("@row").should("contain", "https://www.ugent.be");
      cy.get("@row").should("contain", "Related information");
      cy.get("@row").should("contain", "The initial website");

      cy.get("@row").within(() => {
        cy.get(".if-more").click();

        cy.contains("button", "Edit").click();
      });

      cy.ensureModal("Edit link")
        .within(() => {
          cy.setFieldByLabel("URL", "https://lib.ugent.be");
          cy.setFieldByLabel("Relation", "Accompanying website");
          cy.setFieldByLabel("Description", "The updated website");
        })
        .closeModal("Update link");
      cy.ensureNoModal();

      cy.get("@row").should("have.length", 1);

      cy.get("@row").should("not.contain", "https://www.ugent.be");
      cy.get("@row").should("contain", "https://lib.ugent.be");
      cy.get("@row").should("not.contain", "Related information");
      cy.get("@row").should("contain", "Accompanying website");
      cy.get("@row").should("not.contain", "The initial website");
      cy.get("@row").should("contain", "The updated website");

      cy.get("@row").within(() => {
        cy.get(".if-more").click();

        cy.contains("button", "Delete").click();
      });

      cy.ensureModal("Are you sure").closeModal("Delete");
      cy.ensureNoModal();

      cy.get("#links").find("table tbody tr").should("have.length", 0);
    });

    it("should be possible to add and delete departments", () => {
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
      cy.contains(".dropdown-item", "Remove from dataset").click();

      cy.ensureModal("Are you sure").closeModal("Delete");
      cy.ensureNoModal();

      cy.get("#departments-body .list-group-item-text h4")
        .map("textContent")
        .should("have.ordered.members", ["Faculty of Veterinary Medicine"]);
    });

    it("should be possible to add and edit librarian tags", () => {
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

    it("should be possible to add and edit librarian notes", () => {
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

    it("should be possible to add and edit Biblio message", () => {
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

    it("should be possible to publish, withdraw, republish, lock and unlock a dataset", () => {
      cy.setUpDataset({ prepareForPublishing: true });
      cy.visitDataset();

      cy.contains(".btn-success", "Publish to Biblio").click();
      cy.ensureModal("Are you sure?")
        .within(() => {
          cy.get(".modal-content").should(
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
          cy.get(".modal-content").should(
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
          cy.get(".modal-content").should(
            "contain",
            "Are you sure you want to republish this dataset to Biblio?",
          );
        })
        .closeModal("Republish");
      cy.ensureNoModal();
      cy.ensureToast("Dataset was successfully republished.").closeToast();

      cy.contains(".btn-outline-secondary", "Lock").click();
      cy.ensureNoModal();
      cy.ensureToast("Dataset was successfully locked.").closeToast();

      cy.contains(".btn-outline-secondary", "Unlock").click();
      cy.ensureNoModal();
      cy.ensureToast("Dataset was successfully unlocked.").closeToast();
    });
  });
});
