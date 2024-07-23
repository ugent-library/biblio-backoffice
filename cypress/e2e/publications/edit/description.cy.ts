import { testFormAccessibility } from "support/util";

describe("Editing publication description", () => {
  beforeEach(() => {
    cy.login("researcher1");
  });

  describe("Publication details", () => {
    it("should be possible to change the publication type", () => {
      cy.setUpPublication("Dissertation");
      cy.visitPublication();

      cy.contains(".card", "Publication details")
        .contains(".btn", "Edit")
        .click();

      cy.intercept("/publication/*/type/confirm?type=journal_article").as(
        "changeType",
      );

      cy.ensureModal("Edit publication details").within(() => {
        cy.getLabel("Publication type")
          .next()
          .find("select > option:selected")
          .should("have.value", "dissertation")
          .should("have.text", "Dissertation");

        cy.setFieldByLabel("Publication type", "Journal article");
      });

      cy.wait("@changeType");

      cy.ensureModal("Changing the publication type might result in data loss")
        .within(() => {
          cy.get(".modal-body").should(
            "contain",
            "Are you sure you want to change the type to Journal article?",
          );
        })
        .closeModal("Proceed");
      cy.ensureNoModal();

      cy.contains(".card", "Publication details")
        .find(".card-body")
        .within(() => {
          cy.getLabel("Publication type")
            .next()
            .should("contain", "Journal article");
        });
    });

    it("should error when changing publication type and publication was updated concurrently", () => {
      cy.setUpPublication();
      cy.visitPublication();

      // First perform an update but also capture the snapshot ID
      cy.updateFields(
        "Abstract",
        () => {
          cy.setFieldByLabel("Abstract", "This is an abstract");
          cy.setFieldByLabel("Language", "Danish");

          cy.contains(".modal-footer .btn", "Add abstract")
            .attr("hx-headers")
            .as("initialSnapshot");
        },
        "Add abstract",
      );

      cy.contains(".card", "Publication details")
        .contains(".btn", "Edit")
        .click();

      cy.intercept("/publication/*/type/confirm?type=journal_article").as(
        "changeType",
      );

      cy.ensureModal("Edit publication details").within(() => {
        cy.setFieldByLabel("Publication type", "Journal article");
      });
      cy.wait("@changeType");

      cy.ensureModal("Changing the publication type might result in data loss")
        .within(() => {
          cy.get(".modal-body").should(
            "contain",
            "Are you sure you want to change the type to Journal article?",
          );

          // "Fix" the proceed button with the old (outdated) snapshot ID
          cy.contains(".modal-footer .btn", "Proceed").then((button) => {
            cy.get<string>("@initialSnapshot").then((initialSnapshot) => {
              button.attr("hx-headers", initialSnapshot);
            });
          });
        })
        .closeModal("Proceed");

      cy.ensureModal(null).within(() => {
        cy.contains(
          "Publication has been modified by another user. Please reload the page.",
        ).should("be.visible");
      });
    });

    it("should have clickable labels in the Journal Article form", () => {
      cy.login("librarian1");
      cy.setUpPublication("Journal Article");
      cy.visitPublication();

      cy.updateFields("Publication details", () => {
        const form = {
          "select[name=type]": "Publication type",
          "select[name=journal_article_type]": "Article type",
          "input[type=text][name=doi]": "DOI",
          "select[name=classification]": "Classification",
          "input[type=checkbox][name=legacy]": "Legacy publication",

          "input[type=text][name=title]": "Title",
          "input[type=text][name=alternative_title]": "Alternative title",
          "input[type=text][name=publication]": "Journal title",
          "input[type=text][name=publication_abbreviation]":
            "Short journal title",

          "select[name=language]": "Languages",
          "select[name=publication_status]": "Publishing status",
          ":checkbox[name=extern]":
            "Published while none of the authors and editors were employed at UGent",
          "input[type=text][name=year]": "Publication year",
          "input[type=text][name=place_of_publication]": "Place of publication",
          "input[type=text][name=publisher]": "Publisher",

          "input[type=text][name=volume]": "Volume",
          "input[type=text][name=issue]": "Issue",
          "input[type=text][name=issue_title]": "Special issue title",
          "input[type=text][name=page_first]": "First page",
          "input[type=text][name=page_last]": "Last page",
          "input[type=text][name=page_count]": "Number of pages",
          "input[type=text][name=article_number]": "Article number",

          "input[type=text][name=wos_type]": "Web of Science type",
          "input[type=text][name=wos_id]": "Web of Science ID",
          "input[type=text][name=issn]": "ISSN",
          "input[type=text][name=eissn]": "E-ISSN",
          "input[type=text][name=isbn]": "ISBN",
          "input[type=text][name=eisbn]": "E-ISBN",
          "input[type=text][name=pubmed_id]": "PubMed ID",
          "input[type=text][name=arxiv_id]": "Arxiv ID",
          "input[type=text][name=esci_id]": "ESCI ID",
        };

        testFormAccessibility(form, "Publication type");
      });
    });

    it("should have clickable labels in the Book Chapter form", () => {
      cy.login("librarian1");
      cy.setUpPublication("Book Chapter");

      cy.visitPublication();

      cy.updateFields("Publication details", () => {
        const form = {
          "select[name=type]": "Publication type",
          "input[type=text][name=doi]": "DOI",
          "select[name=classification]": "Classification",
          "input[type=checkbox][name=legacy]": "Legacy publication",

          "input[type=text][name=title]": "Title",
          "input[type=text][name=alternative_title]": "Alternative title",
          "input[type=text][name=publication]": "Book title",

          "select[name=language]": "Languages",
          "select[name=publication_status]": "Publishing status",
          ":checkbox[name=extern]":
            "Published while none of the authors and editors were employed at UGent",
          "input[type=text][name=year]": "Publication year",
          "input[type=text][name=place_of_publication]": "Place of publication",
          "input[type=text][name=publisher]": "Publisher",

          "input[type=text][name=series_title]": "Series title",
          "input[type=text][name=volume]": "Volume",
          "input[type=text][name=edition]": "Edition",
          "input[type=text][name=page_first]": "First page",
          "input[type=text][name=page_last]": "Last page",
          "input[type=text][name=page_count]": "Number of pages",

          "input[type=text][name=wos_type]": "Web of Science type",
          "input[type=text][name=wos_id]": "Web of Science ID",
          "input[type=text][name=issn]": "ISSN",
          "input[type=text][name=eissn]": "E-ISSN",
          "input[type=text][name=isbn]": "ISBN",
          "input[type=text][name=eisbn]": "E-ISBN",
        };

        testFormAccessibility(form, "Publication type");
      });
    });

    it("should have clickable labels in the Book form", () => {
      cy.login("librarian1");
      cy.setUpPublication("Book");

      cy.visitPublication();

      cy.updateFields("Publication details", () => {
        const form = {
          "select[name=type]": "Publication type",
          "input[type=text][name=doi]": "DOI",
          "select[name=classification]": "Classification",
          "input[type=checkbox][name=legacy]": "Legacy publication",

          "input[type=text][name=title]": "Title",
          "input[type=text][name=alternative_title]": "Alternative title",

          "select[name=language]": "Languages",
          "select[name=publication_status]": "Publishing status",
          ":checkbox[name=extern]":
            "Published while none of the authors and editors were employed at UGent",
          "input[type=text][name=year]": "Publication year",
          "input[type=text][name=place_of_publication]": "Place of publication",
          "input[type=text][name=publisher]": "Publisher",

          "input[type=text][name=series_title]": "Series title",
          "input[type=text][name=volume]": "Volume",
          "input[type=text][name=edition]": "Edition",
          "input[type=text][name=page_count]": "Number of pages",

          "input[type=text][name=wos_type]": "Web of Science type",
          "input[type=text][name=wos_id]": "Web of Science ID",
          "input[type=text][name=issn]": "ISSN",
          "input[type=text][name=eissn]": "E-ISSN",
          "input[type=text][name=isbn]": "ISBN",
          "input[type=text][name=eisbn]": "E-ISBN",
        };

        testFormAccessibility(form, "Publication type");
      });
    });

    it("should have clickable labels in the Conference contribution form", () => {
      cy.login("librarian1");
      cy.setUpPublication("Conference contribution");

      cy.visitPublication();

      cy.updateFields("Publication details", () => {
        const form = {
          "select[name=type]": "Publication type",
          "select[name=conference_type]": "Conference type",
          "input[type=text][name=doi]": "DOI",
          "select[name=classification]": "Classification",
          "input[type=checkbox][name=legacy]": "Legacy publication",

          "input[type=text][name=title]": "Title",
          "input[type=text][name=alternative_title]": "Alternative title",
          "input[type=text][name=publication]": "Proceedings title",
          "input[type=text][name=publication_abbreviation]":
            "Publication short title",

          "select[name=language]": "Languages",
          "select[name=publication_status]": "Publishing status",
          ":checkbox[name=extern]":
            "Published while none of the authors and editors were employed at UGent",
          "input[type=text][name=year]": "Publication year",
          "input[type=text][name=place_of_publication]": "Place of publication",
          "input[type=text][name=publisher]": "Publisher",

          "input[type=text][name=series_title]": "Series title",
          "input[type=text][name=volume]": "Volume",
          "input[type=text][name=issue]": "Issue",
          "input[type=text][name=issue_title]": "Special issue title",
          "input[type=text][name=page_first]": "First page",
          "input[type=text][name=page_last]": "Last page",
          "input[type=text][name=page_count]": "Number of pages",
          "input[type=text][name=article_number]": "Article number",

          "input[type=text][name=wos_type]": "Web of Science type",
          "input[type=text][name=wos_id]": "Web of Science ID",
          "input[type=text][name=issn]": "ISSN",
          "input[type=text][name=eissn]": "E-ISSN",
          "input[type=text][name=isbn]": "ISBN",
          "input[type=text][name=eisbn]": "E-ISBN",
        };

        testFormAccessibility(form, "Publication type");
      });
    });

    it("should have clickable labels in the Dissertation form", () => {
      cy.login("librarian1");
      cy.setUpPublication("Dissertation");

      cy.visitPublication();

      cy.updateFields("Publication details", () => {
        const form = {
          "select[name=type]": "Publication type",
          "input[type=text][name=doi]": "DOI",
          "select[name=classification]": "Classification",
          "input[type=checkbox][name=legacy]": "Legacy publication",

          "input[type=text][name=title]": "Title",
          "input[type=text][name=alternative_title]": "Alternative title",

          "select[name=language]": "Languages",
          "select[name=publication_status]": "Publishing status",
          ":checkbox[name=extern]":
            "Published while none of the authors and editors were employed at UGent",
          "input[type=text][name=year]": "Publication year",
          "input[type=text][name=place_of_publication]": "Place of publication",
          "input[type=text][name=publisher]": "Publisher",

          "input[type=text][name=series_title]": "Series title",
          "input[type=text][name=volume]": "Volume",
          "input[type=text][name=page_count]": "Number of pages",

          "input[type=text][name=defense_date]": "Date of defense",
          "input[type=text][name=defense_place]": "Place of defense",

          "input[type=text][name=wos_type]": "Web of Science type",
          "input[type=text][name=wos_id]": "Web of Science ID",
          "input[type=text][name=issn]": "ISSN",
          "input[type=text][name=eissn]": "E-ISSN",
          "input[type=text][name=isbn]": "ISBN",
          "input[type=text][name=eisbn]": "E-ISBN",
        };

        // TODO: fix for radio button fields

        testFormAccessibility(form, "Publication type", [
          "input[type=radio][name=has_confidential_data]",
          "input[type=radio][name=has_patent_application]",
          "input[type=radio][name=has_publications_planned]",
          "input[type=radio][name=has_published_material]",
        ]);
      });
    });

    it("should have clickable labels in the Miscellaneous form", () => {
      cy.login("librarian1");
      cy.setUpPublication("Miscellaneous");

      cy.visitPublication();

      cy.updateFields("Publication details", () => {
        const form = {
          "select[name=type]": "Publication type",
          "select[name=miscellaneous_type]": "Miscellaneous type",
          "input[type=text][name=doi]": "DOI",
          "select[name=classification]": "Classification",
          "input[type=checkbox][name=legacy]": "Legacy publication",

          "input[type=text][name=title]": "Title",
          "input[type=text][name=alternative_title]": "Alternative title",
          "input[type=text][name=publication]": "Publication title",
          "input[type=text][name=publication_abbreviation]":
            "Publication short title",

          "select[name=language]": "Languages",
          "select[name=publication_status]": "Publishing status",
          ":checkbox[name=extern]":
            "Published while none of the authors and editors were employed at UGent",
          "input[type=text][name=year]": "Publication year",
          "input[type=text][name=place_of_publication]": "Place of publication",
          "input[type=text][name=publisher]": "Publisher",

          "input[type=text][name=series_title]": "Series title",
          "input[type=text][name=volume]": "Volume",
          "input[type=text][name=issue]": "Issue",
          "input[type=text][name=issue_title]": "Special issue title",
          "input[type=text][name=edition]": "Edition",
          "input[type=text][name=page_first]": "First page",
          "input[type=text][name=page_last]": "Last page",
          "input[type=text][name=page_count]": "Number of pages",
          "input[type=text][name=article_number]": "Article number",
          "input[type=text][name=report_number]": "Report number",

          "input[type=text][name=wos_type]": "Web of Science type",
          "input[type=text][name=wos_id]": "Web of Science ID",
          "input[type=text][name=issn]": "ISSN",
          "input[type=text][name=eissn]": "E-ISSN",
          "input[type=text][name=isbn]": "ISBN",
          "input[type=text][name=eisbn]": "E-ISBN",
          "input[type=text][name=pubmed_id]": "PubMed ID",
          "input[type=text][name=arxiv_id]": "Arxiv ID",
          "input[type=text][name=esci_id]": "ESCI ID",
        };

        testFormAccessibility(form, "Publication type");
      });
    });

    it("should have clickable labels in the Book editor form", () => {
      cy.login("librarian1");
      cy.setUpPublication("Book editor");

      cy.visitPublication();

      cy.updateFields("Publication details", () => {
        const form = {
          "select[name=type]": "Publication type",
          "input[type=text][name=doi]": "DOI",
          "select[name=classification]": "Classification",
          "input[type=checkbox][name=legacy]": "Legacy publication",

          "input[type=text][name=title]": "Title",
          "input[type=text][name=alternative_title]": "Alternative title",

          "select[name=language]": "Languages",
          "select[name=publication_status]": "Publishing status",
          ":checkbox[name=extern]":
            "Published while none of the authors and editors were employed at UGent",
          "input[type=text][name=year]": "Publication year",
          "input[type=text][name=place_of_publication]": "Place of publication",
          "input[type=text][name=publisher]": "Publisher",

          "input[type=text][name=series_title]": "Series title",
          "input[type=text][name=volume]": "Volume",
          "input[type=text][name=edition]": "Edition",
          "input[type=text][name=page_count]": "Number of pages",

          "input[type=text][name=wos_type]": "Web of Science type",
          "input[type=text][name=wos_id]": "Web of Science ID",
          "input[type=text][name=issn]": "ISSN",
          "input[type=text][name=eissn]": "E-ISSN",
          "input[type=text][name=isbn]": "ISBN",
          "input[type=text][name=eisbn]": "E-ISBN",
        };

        testFormAccessibility(form, "Publication type");
      });
    });

    it("should have clickable labels in the Issue editor form", () => {
      cy.login("librarian1");
      cy.setUpPublication("Issue editor");

      cy.visitPublication();

      cy.updateFields("Publication details", () => {
        const form = {
          "select[name=type]": "Publication type",
          "input[type=text][name=doi]": "DOI",
          "select[name=classification]": "Classification",
          "input[type=checkbox][name=legacy]": "Legacy publication",

          "input[type=text][name=title]": "Title",
          "input[type=text][name=alternative_title]": "Alternative title",
          "input[type=text][name=publication]": "Journal title",
          "input[type=text][name=publication_abbreviation]":
            "Short journal title",

          "select[name=language]": "Languages",
          "select[name=publication_status]": "Publishing status",
          ":checkbox[name=extern]":
            "Published while none of the authors and editors were employed at UGent",
          "input[type=text][name=year]": "Publication year",
          "input[type=text][name=place_of_publication]": "Place of publication",
          "input[type=text][name=publisher]": "Publisher",

          "input[type=text][name=volume]": "Volume",
          "input[type=text][name=issue]": "Issue",
          "input[type=text][name=issue_title]": "Special issue title",
          "input[type=text][name=edition]": "Edition",
          "input[type=text][name=page_count]": "Number of pages",

          "input[type=text][name=wos_type]": "Web of Science type",
          "input[type=text][name=wos_id]": "Web of Science ID",
          "input[type=text][name=issn]": "ISSN",
          "input[type=text][name=eissn]": "E-ISSN",
          "input[type=text][name=isbn]": "ISBN",
          "input[type=text][name=eisbn]": "E-ISBN",
        };

        testFormAccessibility(form, "Publication type");
      });
    });
  });

  describe("Projects", () => {
    beforeEach(() => {
      cy.setUpPublication();
      cy.visitPublication();
    });

    it("should be possible to add and delete projects", () => {
      cy.get("#projects-body").should("contain", "No projects");

      cy.contains(".card", "Project").contains(".btn", "Add project").click();

      cy.ensureModal("Select projects").within(() => {
        cy.intercept("/publication/*/projects/suggestions?*").as(
          "suggestProject",
        );

        cy.setFieldByLabel("Search project", "001D07903");
        cy.wait("@suggestProject");

        cy.contains(".list-group-item", "001D07903")
          .contains(".btn", "Add project")
          .click();
      });
      cy.ensureNoModal();

      cy.get("#projects-body")
        .contains(".list-group-item", "001D07903")
        .find(".if-more")
        .click();
      cy.contains(".dropdown-item", "Remove from publication").click();

      cy.ensureModal("Confirm deletion")
        .within(() => {
          cy.get(".modal-body").should(
            "contain",
            "Are you sure you want to remove this project from the publication?",
          );
        })
        .closeModal("Delete");
      cy.ensureNoModal();

      cy.get("#projects-body").should("contain", "No projects");
    });

    it("should have clickable labels in the project dialog", () => {
      cy.contains(".card", "Project").contains(".btn", "Add project").click();

      cy.ensureModal("Select projects").within(() => {
        cy.get("#project-q").should("be.focused");

        testFormAccessibility(
          {
            "#project-q": "Search project",
          },
          "Search project",
        );
      });
    });
  });

  describe("Abstracts", () => {
    beforeEach(() => {
      cy.setUpPublication();
      cy.visitPublication();
    });

    it("should be possible to add, edit and delete abstracts", () => {
      cy.get("#abstracts-body").should("contain", "No abstracts");

      cy.contains(".btn", "Add abstract").click();
      cy.ensureModal("Add abstract")
        .within(() => {
          cy.setFieldByLabel("Abstract", "");
          cy.setFieldByLabel("Language", "Danish");
        })
        .closeModal("Add abstract");

      cy.ensureModal("Add abstract")
        .within(() => {
          cy.contains(".alert-danger", "Abstract text can't be empty").should(
            "be.visible",
          );
          cy.get("textarea[name=text]")
            .should("have.class", "is-invalid")
            .next(".invalid-feedback")
            .should("have.text", "Abstract text can't be empty");

          cy.setFieldByLabel("Abstract", "The initial abstract");
        })
        .closeModal("Add abstract");
      cy.ensureNoModal();

      cy.get("#abstracts-body")
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
          cy.setFieldByLabel("Abstract", "");
          cy.setFieldByLabel("Language", "Northern Sami");
        })
        .closeModal("Update abstract");

      cy.ensureModal("Edit abstract")
        .within(() => {
          cy.contains(".alert-danger", "Abstract text can't be empty").should(
            "be.visible",
          );
          cy.get("textarea[name=text]")
            .should("have.class", "is-invalid")
            .next(".invalid-feedback")
            .should("have.text", "Abstract text can't be empty");

          cy.setFieldByLabel("Abstract", "The updated abstract");
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

      cy.ensureModal("Confirm deletion")
        .within(() => {
          cy.get(".modal-body").should(
            "contain",
            "Are you sure you want to remove this abstract?",
          );
        })
        .closeModal("Delete");
      cy.ensureNoModal();

      cy.get("#abstracts-body").should("contain", "No abstracts");
    });

    it("should have clickable labels in the Abstract dialog", () => {
      cy.updateFields("Abstract", () => {
        testFormAccessibility(
          {
            "textarea[name=text]": "Abstract",
            "select[name=lang]": "Language",
          },
          "Abstract",
        );
      });
    });
  });

  describe("Links", () => {
    beforeEach(() => {
      cy.setUpPublication();
      cy.visitPublication();
    });

    it("should be possible to add, edit and delete links", () => {
      cy.get("#links-body").should("contain", "No links");

      cy.contains(".btn", "Add link").click();
      cy.ensureModal("Add link")
        .within(() => {
          cy.setFieldByLabel("URL", "https://www.ugent.be");
          cy.setFieldByLabel("Relation", "Related information");
          cy.setFieldByLabel("Description", "The initial website");
        })
        .closeModal("Add link");
      cy.ensureNoModal();

      cy.get("#links-body")
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

      cy.ensureModal("Confirm deletion")
        .within(() => {
          cy.get(".modal-body").should(
            "contain",
            "Are you sure you want to remove this link?",
          );
        })
        .closeModal("Delete");
      cy.ensureNoModal();

      cy.get("#links-body").should("contain", "No links");
    });

    it("should have clickable labels in the Link dialog", () => {
      cy.updateFields("Link", () => {
        testFormAccessibility(
          {
            "input[type=text][name=url]": "URL",
            "select[name=relation]": "Relation",
            "input[type=text][name=description]": "Description",
          },
          "URL",
        );
      });
    });
  });

  describe("Lay summaries", () => {
    beforeEach(() => {
      cy.setUpPublication("Dissertation");
      cy.visitPublication();
    });

    it("should be possible to add, edit and delete lay summaries", () => {
      cy.get("#lay-summaries-body").should("contain", "No lay summaries");

      cy.contains(".btn", "Add lay summary").click();
      cy.ensureModal("Add lay summary")
        .within(() => {
          cy.setFieldByLabel("Lay summary", "");
          cy.setFieldByLabel("Language", "Italian");
        })
        .closeModal("Add lay summary");

      cy.ensureModal("Add lay summary")
        .within(() => {
          cy.contains(
            ".alert-danger",
            "Lay summary text can't be empty",
          ).should("be.visible");
          cy.get("textarea[name=text]")
            .should("have.class", "is-invalid")
            .next(".invalid-feedback")
            .should("have.text", "Lay summary text can't be empty");

          cy.setFieldByLabel("Lay summary", "The initial lay summary");
        })
        .closeModal("Add lay summary");
      cy.ensureNoModal();

      cy.get("#lay-summaries-body")
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
          cy.setFieldByLabel("Lay summary", "");
          cy.setFieldByLabel("Language", "Multiple languages");
        })
        .closeModal("Update lay summary");

      cy.ensureModal("Edit lay summary")
        .within(() => {
          cy.contains(
            ".alert-danger",
            "Lay summary text can't be empty",
          ).should("be.visible");
          cy.get("textarea[name=text]")
            .should("have.class", "is-invalid")
            .next(".invalid-feedback")
            .should("have.text", "Lay summary text can't be empty");

          cy.setFieldByLabel("Lay summary", "The updated lay summary");
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

      cy.ensureModal("Confirm deletion")
        .within(() => {
          cy.get(".modal-body").should(
            "contain",
            "Are you sure you want to remove this lay summary?",
          );
        })
        .closeModal("Delete");
      cy.ensureNoModal();

      cy.get("#lay-summaries-body").should("contain", "No lay summaries");
    });

    it("should have clickable labels in the Lay summary dialog", () => {
      cy.updateFields("Lay summary", () => {
        testFormAccessibility(
          {
            "textarea[name=text]": "Lay summary",
            "select[name=lang]": "Language",
          },
          "Lay summary",
        );
      });
    });
  });

  describe("Conference details", () => {
    beforeEach(() => {
      cy.setUpPublication("Conference contribution");
      cy.visitPublication();
    });

    it("should be possible to add and edit conference details", () => {
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

    it("should have clickable labels in the Conference details dialog", () => {
      cy.updateFields("Conference details", () => {
        testFormAccessibility(
          {
            "input[type=text][name=name]": "Conference",
            "input[type=text][name=location]": "Conference location",
            "input[type=text][name=organizer]": "Conference organiser",
            "input[type=text][name=start_date]": "Conference start date",
            "input[type=text][name=end_date]": "Conference end date",
          },
          "Conference",
        );
      });
    });
  });

  describe("Additional info", () => {
    beforeEach(() => {
      cy.setUpPublication();
      cy.visitPublication();
    });

    it("should be possible to add and edit additional info", () => {
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

          cy.setFieldByLabel(
            "Keywords",
            "these{enter}are{enter}the{enter}keywords",
          );

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
            .next()
            .find(".form-values")
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
          cy.get("tags span[contenteditable]").type("updated");

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

    it("should have clickable labels in the Additional info dialog", () => {
      cy.updateFields("Additional information", () => {
        cy.focused().should("have.attr", "id", "research_field-0");

        testFormAccessibility(
          {
            "select[name=research_field]": "Research field",
            "textarea[name=additional_info]": "Additional information",
            ".tags:has(textarea#keyword) tags span.tagify__input[contenteditable]":
              "Keywords",
          },
          "Research field",
          ["textarea[data-input-name=keyword]"],
        );

        cy.setFieldByLabel("Research field", "General Works");
      });
    });
  });
});
