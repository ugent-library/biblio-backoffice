// https://github.com/ugent-library/biblio-backoffice/issues/1402

describe("Issue #1402: Gohtml conversion to Templ", () => {
  describe("as researcher", () => {
    beforeEach(() => {
      cy.loginAsResearcher();
    });

    describe("for publications", () => {
      it("should be possible to delete a publication", () => {
        cy.setUpPublication();
        cy.visitPublication();

        cy.get(".btn .if-more").click();
        cy.contains(".dropdown-item", "Delete").click();

        cy.ensureModal("Are you sure?")
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

        cy.ensureModal(
          "Changing the publication type might result in data loss",
        )
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

        cy.ensureModal(
          "Changing the publication type might result in data loss",
        )
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

      it("should be possible to add, edit and delete abstracts", () => {
        cy.setUpPublication();
        cy.visitPublication();

        cy.get("#abstracts").find("table tbody tr").should("have.length", 0);

        cy.contains(".btn", "Add abstract").click();
        cy.ensureModal("Add abstract")
          .within(() => {
            cy.setFieldByLabel("Abstract", " ");
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

        cy.ensureModal("Are you sure")
          .within(() => {
            cy.get(".modal-body").should(
              "contain",
              "Are you sure you want to remove this abstract?",
            );
          })
          .closeModal("Delete");
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

        cy.ensureModal("Are you sure")
          .within(() => {
            cy.get(".modal-body").should(
              "contain",
              "Are you sure you want to remove this link?",
            );
          })
          .closeModal("Delete");
        cy.ensureNoModal();

        cy.get("#links").find("table tbody tr").should("have.length", 0);
      });

      it("should be possible to add, edit and delete lay summaries", () => {
        cy.setUpPublication("Dissertation");
        cy.visitPublication();

        cy.get("#lay-summaries")
          .find("table tbody tr")
          .should("have.length", 0);

        cy.contains(".btn", "Add lay summary").click();
        cy.ensureModal("Add lay summary")
          .within(() => {
            cy.setFieldByLabel("Lay summary", " ");
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

        cy.ensureModal("Are you sure")
          .within(() => {
            cy.get(".modal-body").should(
              "contain",
              "Are you sure you want to remove this lay summary?",
            );
          })
          .closeModal("Delete");
        cy.ensureNoModal();

        cy.get("#lay-summaries")
          .find("table tbody tr")
          .should("have.length", 0);
      });

      it("should be possible to add and edit conference details", () => {
        cy.setUpPublication("Conference contribution");
        cy.visitPublication();

        cy.get("#conference-details").contains(".btn", "Edit").click();
        cy.ensureModal("Edit conference details")
          .within(() => {
            cy.setFieldByLabel("Conference", "The conference name");
            cy.setFieldByLabel(
              "Conference location",
              "The conference location",
            );
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

      it("should be possible to delete authors", () => {
        cy.setUpPublication();
        cy.visitPublication();

        cy.updateFields(
          "Authors",
          () => {
            cy.intercept("/publication/*/contributors/author/suggestions?*").as(
              "suggestContributor",
            );

            cy.setFieldByLabel("First name", "Jane");
            cy.wait("@suggestContributor");
            cy.setFieldByLabel("Last name", "Doe");
            cy.wait("@suggestContributor");

            cy.contains(".btn", "Add external author").click();
          },
          true,
        );

        cy.contains("#authors tr", "Jane Doe").find(".btn .if-delete").click();

        cy.ensureModal("Are you sure?")
          .within(() => {
            cy.contains("Are you sure you want to remove this author?").should(
              "be.visible",
            );
          })
          .closeModal("Delete");
        cy.ensureNoModal();

        cy.contains("#authors", "Jane Doe").should("not.exist");
      });

      it("should not be possible to delete the last UGent author of a published publication", () => {
        cy.setUpPublication(undefined, { prepareForPublishing: true });
        cy.visitPublication();

        cy.contains(".btn", "Publish to Biblio").click();
        cy.ensureModal("Are you sure?").closeModal("Publish");
        cy.ensureToast("Publication was successfully published.").closeToast();

        cy.updateFields(
          "Authors",
          () => {
            cy.intercept("/publication/*/contributors/author/suggestions?*").as(
              "suggestContributor",
            );

            cy.setFieldByLabel("First name", "Jane");
            cy.wait("@suggestContributor");
            cy.setFieldByLabel("Last name", "Doe");
            cy.wait("@suggestContributor");
            cy.contains(".btn", "Add external author").click();
          },
          true,
        );

        cy.contains("#authors tr", "John Doe").find(".btn .if-delete").click();
        cy.ensureModal("Are you sure?").closeModal("Delete");

        cy.ensureModal(
          "Can't delete this contributor due to the following errors",
        ).within(() => {
          cy.contains(
            ".alert-danger",
            "At least one UGent author is required",
          ).should("be.visible");
        });
      });

      it("should be possible to delete editors", () => {
        cy.setUpPublication("Book");
        cy.visitPublication();

        cy.updateFields(
          "Editors",
          () => {
            cy.intercept("/publication/*/contributors/editor/suggestions?*").as(
              "suggestContributor",
            );

            cy.setFieldByLabel("First name", "Jane");
            cy.wait("@suggestContributor");
            cy.setFieldByLabel("Last name", "Doe");
            cy.wait("@suggestContributor");

            cy.contains(".btn", "Add external editor").click();
          },
          true,
        );

        cy.contains("#editors tr", "Jane Doe").find(".btn .if-delete").click();

        cy.ensureModal("Are you sure?")
          .within(() => {
            cy.contains("Are you sure you want to remove this editor?").should(
              "be.visible",
            );
          })
          .closeModal("Delete");
        cy.ensureNoModal();

        cy.contains("#editors", "Jane Doe").should("not.exist");
      });

      it("should be possible to delete supervisors", () => {
        cy.setUpPublication("Dissertation");
        cy.visitPublication();

        cy.updateFields(
          "Supervisors",
          () => {
            cy.intercept(
              "/publication/*/contributors/supervisor/suggestions?*",
            ).as("suggestContributor");

            cy.setFieldByLabel("First name", "Jane");
            cy.wait("@suggestContributor");
            cy.setFieldByLabel("Last name", "Doe");
            cy.wait("@suggestContributor");

            cy.contains(".btn", "Add external supervisor").click();
          },
          true,
        );

        cy.contains("#supervisors tr", "Jane Doe")
          .find(".btn .if-delete")
          .click();

        cy.ensureModal("Are you sure?")
          .within(() => {
            cy.contains(
              "Are you sure you want to remove this supervisor?",
            ).should("be.visible");
          })
          .closeModal("Delete");
        cy.ensureNoModal();

        cy.contains("#supervisors", "Jane Doe").should("not.exist");
      });

      it("should be possible to add and delete departments", () => {
        cy.setUpPublication();
        cy.visitPublication();

        cy.contains(".nav .nav-item", "People & Affiliations").click();

        cy.get("#departments").contains(".btn", "Add department").click();

        cy.ensureModal("Select departments").within(() => {
          cy.intercept("/publication/*/departments/suggestions?*").as(
            "suggestDepartment",
          );

          cy.getLabel("Search").next("input").type("LW17");
          cy.wait("@suggestDepartment");

          cy.contains(".btn", "Add department").click();
        });
        cy.ensureNoModal();

        cy.get("#departments-body .list-group-item-text h4")
          .map("textContent")
          .should("have.ordered.members", [
            "Department of Art, music and theatre sciences",
          ]);

        cy.get("#departments").contains(".btn", "Add department").click();

        cy.ensureModal("Select departments").within(() => {
          cy.getLabel("Search").next("input").type("DI62");
          cy.wait("@suggestDepartment");

          cy.contains(".btn", "Add department").click();
        });
        cy.ensureNoModal();

        cy.get("#departments-body .list-group-item-text h4")
          .map("textContent")
          .should("have.ordered.members", [
            "Department of Art, music and theatre sciences",
            "Biocenter AgriVet",
          ]);

        cy.contains(
          "#departments-body tr",
          "Department of Art, music and theatre sciences",
        )
          .find(".if-more")
          .click();
        cy.contains(".dropdown-item", "Remove from publication").click();

        cy.ensureModal("Are you sure")
          .within(() => {
            cy.get(".modal-body").should(
              "contain",
              "Are you sure you want to remove this department from the publication?",
            );
          })
          .closeModal("Delete");
        cy.ensureNoModal();

        cy.get("#departments-body .list-group-item-text h4")
          .map("textContent")
          .should("have.ordered.members", ["Biocenter AgriVet"]);
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

      it("should be possible to publish, withdraw and republish a publication", () => {
        cy.setUpPublication(undefined, { prepareForPublishing: true });
        cy.visitPublication();

        cy.contains(".btn-success", "Publish to Biblio").click();
        cy.ensureModal("Are you sure?")
          .within(() => {
            cy.get(".modal-body").should(
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
            cy.get(".modal-body").should(
              "contain",
              "Are you sure you want to withdraw this publication from Biblio?",
            );
          })
          .closeModal("Withdraw");
        cy.ensureNoModal();
        cy.ensureToast("Publication was successfully withdrawn.").closeToast();

        cy.contains(".btn-success", "Republish to Biblio").click();
        cy.ensureModal("Are you sure?")
          .within(() => {
            cy.get(".modal-body").should(
              "contain",
              "Are you sure you want to republish this publication to Biblio?",
            );
          })
          .closeModal("Republish");
        cy.ensureNoModal();
        cy.ensureToast(
          "Publication was successfully republished.",
        ).closeToast();
      });

      it("should error when publication is not ready for publication", () => {
        cy.setUpPublication("Miscellaneous", { prepareForPublishing: false });
        cy.visitPublication();

        cy.contains(".btn-success", "Publish to Biblio").click();
        cy.ensureModal("Are you sure?")
          .within(() => {
            cy.get(".modal-body").should(
              "contain",
              "Are you sure you want to publish this publication to Biblio?",
            );
          })
          .closeModal("Publish");

        cy.ensureModal(
          "Unable to publish this publication due to the following errors",
        )
          .within(() => {
            cy.get("ul > li")
              .map("textContent")
              .should("have.members", [
                "Publication year is required",
                "One or more authors are required",
                "At least one UGent author is required",
              ]);
          })
          .closeModal("Close");
        cy.ensureNoModal();
        cy.ensureNoToast();

        cy.reload();

        cy.contains(".btn-success", "Publish to Biblio").should("be.visible");
      });

      it("should error when publication is not ready for republication", () => {
        cy.setUpPublication("Miscellaneous", { prepareForPublishing: true });
        cy.visitPublication();

        cy.contains(".btn-success", "Publish to Biblio").click();
        cy.ensureModal("Are you sure?").closeModal("Publish");
        cy.ensureNoModal();
        cy.ensureToast("Publication was successfully published.").closeToast();

        cy.contains(".btn-outline-danger", "Withdraw").click();
        cy.ensureModal("Are you sure?").closeModal("Withdraw");
        cy.ensureNoModal();
        cy.ensureToast("Publication was successfully withdrawn.").closeToast();

        cy.updateFields(
          "Publication details",
          () => {
            cy.setFieldByLabel("Publication year", " ");
          },
          true,
        );

        cy.contains(".nav-link", "People & Affiliations").click();

        // Delete interal author
        cy.get("#authors button:has(.if-delete)").click();
        cy.ensureModal("Are you sure?").closeModal("Delete");

        // Add external author
        cy.updateFields(
          "Authors",
          () => {
            cy.intercept("/publication/*/contributors/author/suggestions?*").as(
              "suggestContributor",
            );
            cy.setFieldByLabel("First name", "John");
            cy.wait("@suggestContributor");
            cy.setFieldByLabel("Last name", "Doe");
            cy.wait("@suggestContributor");
            cy.contains(".btn", "Add external author").click();
          },
          true,
        );

        cy.contains(".btn-success", "Republish to Biblio").click();
        cy.ensureModal("Are you sure?")
          .within(() => {
            cy.get(".modal-body").should(
              "contain",
              "Are you sure you want to republish this publication to Biblio?",
            );
          })
          .closeModal("Republish");

        cy.ensureModal(
          "Unable to republish this publication due to the following errors",
        )
          .within(() => {
            cy.get("ul > li")
              .map("textContent")
              .should("have.members", [
                "Publication year is required",
                "At least one UGent author is required",
              ]);
          })
          .closeModal("Close");
        cy.ensureNoModal();
        cy.ensureNoToast();

        cy.reload();

        cy.contains(".btn-success", "Republish to Biblio").should("be.visible");
      });
    });

    describe("for datasets", () => {
      beforeEach(() => {
        cy.setUpDataset();
        cy.visitDataset();
      });

      it("should be possible to delete a dataset", () => {
        cy.setUpDataset();
        cy.visitDataset();

        cy.get(".btn .if-more").click();
        cy.contains(".dropdown-item", "Delete").click();

        cy.ensureModal("Are you sure?")
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

      it("should be possible to add, edit and delete abstracts", () => {
        cy.get("#abstracts").find("table tbody tr").should("have.length", 0);

        cy.contains(".btn", "Add abstract").click();
        cy.ensureModal("Add abstract")
          .within(() => {
            cy.setFieldByLabel("Abstract", " ");
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

        cy.ensureModal("Are you sure")
          .within(() => {
            cy.get(".modal-body").should(
              "contain",
              "Are you sure you want to remove this abstract?",
            );
          })
          .closeModal("Delete");
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

        cy.ensureModal("Are you sure")
          .within(() => {
            cy.get(".modal-body").should(
              "contain",
              "Are you sure you want to remove this link?",
            );
          })
          .closeModal("Delete");
        cy.ensureNoModal();

        cy.get("#links").find("table tbody tr").should("have.length", 0);
      });

      it("should be possible to delete creators", () => {
        cy.setUpDataset();
        cy.visitDataset();

        cy.updateFields(
          "Creators",
          () => {
            cy.intercept("/dataset/*/contributors/author/suggestions?*").as(
              "suggestContributor",
            );

            cy.setFieldByLabel("First name", "Jane");
            cy.wait("@suggestContributor");
            cy.setFieldByLabel("Last name", "Doe");
            cy.wait("@suggestContributor");

            cy.contains(".btn", "Add external creator").click();
          },
          true,
        );

        cy.contains("#authors tr", "Jane Doe").find(".btn .if-delete").click();

        cy.ensureModal("Are you sure?")
          .within(() => {
            cy.contains("Are you sure you want to remove this creator?").should(
              "be.visible",
            );
          })
          .closeModal("Delete");
        cy.ensureNoModal();

        cy.contains("#authors", "Jane Doe").should("not.exist");
      });

      it("should not be possible to delete the last UGent creator of a published dataset", () => {
        cy.setUpDataset({ prepareForPublishing: true });
        cy.visitDataset();

        cy.contains(".btn", "Publish to Biblio").click();
        cy.ensureModal("Are you sure?").closeModal("Publish");
        cy.ensureToast("Dataset was successfully published.").closeToast();

        cy.updateFields(
          "Creators",
          () => {
            cy.intercept("/dataset/*/contributors/author/suggestions?*").as(
              "suggestContributor",
            );

            cy.setFieldByLabel("First name", "Jane");
            cy.wait("@suggestContributor");
            cy.setFieldByLabel("Last name", "Doe");
            cy.wait("@suggestContributor");
            cy.contains(".btn", "Add external creator").click();
          },
          true,
        );

        cy.contains("#authors tr", "John Doe").find(".btn .if-delete").click();
        cy.ensureModal("Are you sure?").closeModal("Delete");

        cy.ensureModal(
          "Can't delete this contributor due to the following errors",
        ).within(() => {
          cy.contains(
            ".alert-danger",
            "At least one UGent author is required",
          ).should("be.visible");
        });
      });

      it("should be possible to add and delete departments", () => {
        cy.contains(".nav .nav-item", "People & Affiliations").click();

        cy.get("#departments").contains(".btn", "Add department").click();

        cy.ensureModal("Select departments").within(() => {
          cy.intercept("/dataset/*/departments/suggestions?*").as(
            "suggestDepartment",
          );

          cy.getLabel("Search").next("input").type("LW17");
          cy.wait("@suggestDepartment");

          cy.contains(".btn", "Add department").click();
        });
        cy.ensureNoModal();

        cy.get("#departments-body .list-group-item-text h4")
          .map("textContent")
          .should("have.ordered.members", [
            "Department of Art, music and theatre sciences",
          ]);

        cy.get("#departments").contains(".btn", "Add department").click();

        cy.ensureModal("Select departments").within(() => {
          cy.getLabel("Search").next("input").type("DI62");
          cy.wait("@suggestDepartment");

          cy.contains(".btn", "Add department").click();
        });
        cy.ensureNoModal();

        cy.get("#departments-body .list-group-item-text h4")
          .map("textContent")
          .should("have.ordered.members", [
            "Department of Art, music and theatre sciences",
            "Biocenter AgriVet",
          ]);

        cy.contains(
          "#departments-body tr",
          "Department of Art, music and theatre sciences",
        )
          .find(".if-more")
          .click();
        cy.contains(".dropdown-item", "Remove from dataset").click();

        cy.ensureModal("Are you sure")
          .within(() => {
            cy.get(".modal-body").should(
              "contain",
              "Are you sure you want to remove this department from the dataset?",
            );
          })
          .closeModal("Delete");
        cy.ensureNoModal();

        cy.get("#departments-body .list-group-item-text h4")
          .map("textContent")
          .should("have.ordered.members", ["Biocenter AgriVet"]);
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

      it("should be possible to publish, withdraw and republish a dataset", () => {
        cy.setUpDataset({ prepareForPublishing: true });
        cy.visitDataset();

        cy.contains(".btn-success", "Publish to Biblio").click();
        cy.ensureModal("Are you sure?")
          .within(() => {
            cy.get(".modal-body").should(
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
            cy.get(".modal-body").should(
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
            cy.get(".modal-body").should(
              "contain",
              "Are you sure you want to republish this dataset to Biblio?",
            );
          })
          .closeModal("Republish");
        cy.ensureNoModal();
        cy.ensureToast("Dataset was successfully republished.").closeToast();
      });

      it("should error when dataset is not ready for publication", () => {
        cy.setUpDataset({ prepareForPublishing: false });
        cy.visitDataset();

        cy.contains(".btn-success", "Publish to Biblio").click();
        cy.ensureModal("Are you sure?")
          .within(() => {
            cy.get(".modal-body").should(
              "contain",
              "Are you sure you want to publish this dataset to Biblio?",
            );
          })
          .closeModal("Publish");

        cy.ensureModal(
          "Unable to publish this dataset due to the following errors",
        )
          .within(() => {
            cy.get("ul > li")
              .map("textContent")
              .should("have.members", [
                "Access level is required",
                "Format is required",
                "Publisher is required",
                "Publication year is required",
                "One or more authors are required",
                "At least one UGent author is required",
                "License is required",
              ]);
          })
          .closeModal("Close");
        cy.ensureNoModal();
        cy.ensureNoToast();

        cy.reload();

        cy.contains(".btn-success", "Publish to Biblio").should("be.visible");
      });

      it("should error when dataset is not ready for republication", () => {
        cy.setUpDataset({ prepareForPublishing: true });
        cy.visitDataset();

        cy.contains(".btn-success", "Publish to Biblio").click();
        cy.ensureModal("Are you sure?").closeModal("Publish");
        cy.ensureNoModal();
        cy.ensureToast("Dataset was successfully published.").closeToast();

        cy.contains(".btn-outline-danger", "Withdraw").click();
        cy.ensureModal("Are you sure?").closeModal("Withdraw");
        cy.ensureNoModal();
        cy.ensureToast("Dataset was successfully withdrawn.").closeToast();

        cy.updateFields(
          "Dataset details",
          () => {
            cy.setFieldByLabel("Publisher", " ");
            cy.setFieldByLabel("Publication year", " ");
          },
          true,
        );

        cy.contains(".nav-link", "People & Affiliations").click();

        // Delete interal author
        cy.get("#authors button:has(.if-delete)").click();
        cy.ensureModal("Are you sure?").closeModal("Delete");

        // Add external author
        cy.updateFields(
          "Creators",
          () => {
            cy.intercept("/dataset/*/contributors/author/suggestions?*").as(
              "suggestContributor",
            );
            cy.setFieldByLabel("First name", "John");
            cy.wait("@suggestContributor");
            cy.setFieldByLabel("Last name", "Doe");
            cy.wait("@suggestContributor");
            cy.contains(".btn", "Add external creator").click();
          },
          true,
        );

        cy.contains(".btn-success", "Republish to Biblio").click();
        cy.ensureModal("Are you sure?")
          .within(() => {
            cy.get(".modal-body").should(
              "contain",
              "Are you sure you want to republish this dataset to Biblio?",
            );
          })
          .closeModal("Republish");

        cy.ensureModal(
          "Unable to republish this dataset due to the following errors",
        )
          .within(() => {
            cy.get("ul > li")
              .map("textContent")
              .should("have.members", [
                "Publisher is required",
                "Publication year is required",
                "At least one UGent author is required",
              ]);
          })
          .closeModal("Close");
        cy.ensureNoModal();
        cy.ensureNoToast();

        cy.reload();

        cy.contains(".btn-success", "Republish to Biblio").should("be.visible");
      });
    });

    describe("media type suggestions", () => {
      it("should provide format type suggestions", () => {
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
          .should(
            "have.attr",
            "data-value",
            "application/vnd.google-earth.kml+xml",
          )
          .find(".badge")
          .should("contains.text", "application/vnd.google-earth.kml+xml")
          .parent()
          .prop("innerText")
          .should("contain", "application/vnd.google-earth.kml+xml (.kml)");
      });
    });
  });

  describe("as librarian", () => {
    beforeEach(() => {
      cy.loginAsLibrarian();
    });

    describe("for publications", () => {
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

      it("should be possible to lock and unlock a publication", () => {
        cy.setUpPublication(undefined, { prepareForPublishing: true });
        cy.visitPublication();

        cy.contains(".btn-success", "Publish to Biblio").click();
        cy.ensureModal("Are you sure?").closeModal("Publish");
        cy.ensureToast("Publication was successfully published.").closeToast();

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

      it("should be possible to lock and unlock a dataset", () => {
        cy.setUpDataset({ prepareForPublishing: true });
        cy.visitDataset();

        cy.contains(".btn-success", "Publish to Biblio").click();
        cy.ensureModal("Are you sure?").closeModal("Publish");
        cy.ensureToast("Dataset was successfully published.").closeToast();

        cy.contains(".btn-outline-secondary", "Lock").click();
        cy.ensureNoModal();
        cy.ensureToast("Dataset was successfully locked.").closeToast();

        cy.contains(".btn-outline-secondary", "Unlock").click();
        cy.ensureNoModal();
        cy.ensureToast("Dataset was successfully unlocked.").closeToast();
      });
    });
  });

  describe("authorization", () => {
    it("should not be possible to edit or delete publications from another user", () => {
      cy.loginAsLibrarian();
      cy.setUpPublication();

      cy.loginAsResearcher();

      testForbiddenPublicationRoute("/add/description");
      testForbiddenPublicationRoute("/add/confirm");
      testForbiddenPublicationRoute("/add/publish", "POST");
      testForbiddenPublicationRoute("/add/finish");

      testForbiddenPublicationRoute(""); // The regular view publication route
      testForbiddenPublicationRoute("/description");
      testForbiddenPublicationRoute("/files");
      testForbiddenPublicationRoute("/contributors");
      testForbiddenPublicationRoute("/datasets");
      testForbiddenPublicationRoute("/activity");
      testForbiddenPublicationRoute("/files/file-123");

      testForbiddenPublicationRoute("/confirm-delete");
      testForbiddenPublicationRoute("", "DELETE");
      testForbiddenPublicationRoute("/lock", "POST");
      testForbiddenPublicationRoute("/unlock", "POST");
      testForbiddenPublicationRoute("/publish/confirm");
      testForbiddenPublicationRoute("/publish", "POST");
      testForbiddenPublicationRoute("/withdraw/confirm");
      testForbiddenPublicationRoute("/withdraw", "POST");
      testForbiddenPublicationRoute("/republish/confirm");
      testForbiddenPublicationRoute("/republish", "POST");

      testForbiddenPublicationRoute("/message/edit", "GET", "PUT");
      testForbiddenPublicationRoute("/message", "PUT");
      testForbiddenPublicationRoute("/reviewer-tags/edit");
      testForbiddenPublicationRoute("/reviewer-tags", "PUT");
      testForbiddenPublicationRoute("/reviewer-note/edit");
      testForbiddenPublicationRoute("/reviewer-note", "PUT");

      testForbiddenPublicationRoute("/details/edit", "GET", "PUT");
      testForbiddenPublicationRoute("/type/confirm");
      testForbiddenPublicationRoute("/type", "PUT");

      testForbiddenPublicationRoute("/conference/edit");
      testForbiddenPublicationRoute("/conference", "PUT");
      testForbiddenPublicationRoute("/additional-info/edit");
      testForbiddenPublicationRoute("/additional-info", "PUT");

      testForbiddenPublicationRoute("/projects/add");
      testForbiddenPublicationRoute("/projects/suggestions");
      testForbiddenPublicationRoute("/projects", "POST");
      testForbiddenPublicationRoute(
        "/snapshot-123/projects/confirm-delete/project-123",
      );
      testForbiddenPublicationRoute("/projects/project-123", "DELETE");

      testForbiddenPublicationRoute("/links/add");
      testForbiddenPublicationRoute("/links", "POST");
      testForbiddenPublicationRoute("/links/link-123/edit");
      testForbiddenPublicationRoute("/links/link-123", "PUT", "DELETE");
      testForbiddenPublicationRoute(
        "/snapshot-123/links/link-123/confirm-delete",
      );

      testForbiddenPublicationRoute("/departments/add");
      testForbiddenPublicationRoute("/departments/suggestions");
      testForbiddenPublicationRoute("/departments", "POST");
      testForbiddenPublicationRoute(
        "/snapshot-123/departments/department-123/confirm-delete",
      );
      testForbiddenPublicationRoute("/departments/department-123", "DELETE");

      testForbiddenPublicationRoute("/abstracts/add");
      testForbiddenPublicationRoute("/abstracts", "POST");
      testForbiddenPublicationRoute("/abstracts/abstract-123/edit");
      testForbiddenPublicationRoute("/abstracts/abstract-123", "PUT", "DELETE");
      testForbiddenPublicationRoute(
        "/snapshot-123/abstracts/abstract-123/confirm-delete",
      );

      testForbiddenPublicationRoute("/lay_summaries/add");
      testForbiddenPublicationRoute("/lay_summaries", "POST");
      testForbiddenPublicationRoute("/lay_summaries/lay-summary-123/edit");
      testForbiddenPublicationRoute(
        "/lay_summaries/lay-summary-123",
        "PUT",
        "DELETE",
      );
      testForbiddenPublicationRoute(
        "/snapshot-123/lay_summaries/lay-summary-123/confirm-delete",
      );

      testForbiddenPublicationRoute("/datasets/add");
      testForbiddenPublicationRoute("/datasets/suggestions");
      testForbiddenPublicationRoute("/datasets", "POST");
      testForbiddenPublicationRoute(
        "/snapshot-123/datasets/dataset-123/confirm-delete",
      );
      testForbiddenPublicationRoute("/datasets/dataset-123", "DELETE");

      testForbiddenPublicationRoute("/contributors/role-123/order", "POST");
      testForbiddenPublicationRoute("/contributors/role-123/add");
      testForbiddenPublicationRoute("/contributors/role-123/suggestions");
      testForbiddenPublicationRoute("/contributors/role-123/confirm-create");
      testForbiddenPublicationRoute("/contributors/role-123", "POST");
      testForbiddenPublicationRoute("/contributors/role-123/position-123/edit");
      testForbiddenPublicationRoute(
        "/contributors/role-123/position-123/suggestions",
      );
      testForbiddenPublicationRoute(
        "/contributors/role-123/position-123/confirm-update",
      );
      testForbiddenPublicationRoute(
        "/contributors/role-123/position-123/confirm-delete",
      );
      testForbiddenPublicationRoute(
        "/contributors/role-123/position-123",
        "PUT",
        "DELETE",
      );

      testForbiddenPublicationRoute("/files", "POST");
      testForbiddenPublicationRoute("/files/file-123/edit");
      testForbiddenPublicationRoute("/refresh-files");
      testForbiddenPublicationRoute("/files/file-123/refresh-form");
      testForbiddenPublicationRoute("/files/file-123", "PUT", "DELETE");
      testForbiddenPublicationRoute(
        "/snapshot123/files/file-123/confirm-delete",
      );
    });

    it("should not be possible to edit or delete datasets from another user", () => {
      cy.loginAsLibrarian();
      cy.setUpDataset();

      cy.loginAsResearcher();

      testForbiddenDatasetRoute("/add/description");
      testForbiddenDatasetRoute("/add/confirm");
      testForbiddenDatasetRoute("/save", "POST");
      testForbiddenDatasetRoute("/add/publish", "POST");
      testForbiddenDatasetRoute("/add/finish");

      testForbiddenDatasetRoute(""); // The regular view dataset route
      testForbiddenDatasetRoute("/description");
      testForbiddenDatasetRoute("/contributors");
      testForbiddenDatasetRoute("/publications");
      testForbiddenDatasetRoute("/activity");

      testForbiddenDatasetRoute("/confirm-delete");
      testForbiddenDatasetRoute("", "DELETE");
      testForbiddenDatasetRoute("/lock", "POST");
      testForbiddenDatasetRoute("/unlock", "POST");
      testForbiddenDatasetRoute("/publish/confirm");
      testForbiddenDatasetRoute("/publish", "POST");
      testForbiddenDatasetRoute("/withdraw/confirm");
      testForbiddenDatasetRoute("/withdraw", "POST");
      testForbiddenDatasetRoute("/republish/confirm");
      testForbiddenDatasetRoute("/republish", "POST");

      testForbiddenDatasetRoute("/message/edit", "GET", "PUT");
      testForbiddenDatasetRoute("/message", "PUT");
      testForbiddenDatasetRoute("/reviewer-tags/edit");
      testForbiddenDatasetRoute("/reviewer-tags", "PUT");
      testForbiddenDatasetRoute("/reviewer-note/edit");
      testForbiddenDatasetRoute("/reviewer-note", "PUT");

      testForbiddenDatasetRoute("/details/edit", "GET", "PUT");
      testForbiddenDatasetRoute("/details/edit/refresh-form", "PUT");

      testForbiddenDatasetRoute("/projects/add");
      testForbiddenDatasetRoute("/projects/suggestions");
      testForbiddenDatasetRoute("/projects", "POST");
      testForbiddenDatasetRoute(
        "/snapshot-123/projects/confirm-delete/project-123",
      );
      testForbiddenDatasetRoute("/projects/project-123", "DELETE");

      testForbiddenDatasetRoute("/links/add");
      testForbiddenDatasetRoute("/links", "POST");
      testForbiddenDatasetRoute("/links/link-123/edit");
      testForbiddenDatasetRoute("/links/link-123", "PUT", "DELETE");
      testForbiddenDatasetRoute("/snapshot-123/links/link-123/confirm-delete");

      testForbiddenDatasetRoute("/departments/add");
      testForbiddenDatasetRoute("/departments/suggestions");
      testForbiddenDatasetRoute("/departments", "POST");
      testForbiddenDatasetRoute(
        "/snapshot-123/departments/department-123/confirm-delete",
      );
      testForbiddenDatasetRoute("/departments/department-123", "DELETE");

      testForbiddenDatasetRoute("/abstracts/add");
      testForbiddenDatasetRoute("/abstracts", "POST");
      testForbiddenDatasetRoute("/abstracts/abstract-123/edit");
      testForbiddenDatasetRoute("/abstracts/abstract-123", "PUT", "DELETE");
      testForbiddenDatasetRoute(
        "/snapshot-123/abstracts/abstract-123/confirm-delete",
      );

      testForbiddenDatasetRoute("/publications/add");
      testForbiddenDatasetRoute("/publications/suggestions");
      testForbiddenDatasetRoute("/publications", "POST");
      testForbiddenDatasetRoute(
        "/snapshot-123/publications/publication-123/confirm-delete",
      );
      testForbiddenDatasetRoute("/publications/publication-123", "DELETE");

      testForbiddenDatasetRoute("/contributors/role-123/order", "POST");
      testForbiddenDatasetRoute("/contributors/role-123/add");
      testForbiddenDatasetRoute("/contributors/role-123/suggestions");
      testForbiddenDatasetRoute("/contributors/role-123/confirm-create");
      testForbiddenDatasetRoute("/contributors/role-123", "POST");
      testForbiddenDatasetRoute("/contributors/role-123/position-123/edit");
      testForbiddenDatasetRoute(
        "/contributors/role-123/position-123/suggestions",
      );
      testForbiddenDatasetRoute(
        "/contributors/role-123/position-123/confirm-update",
      );
      testForbiddenDatasetRoute(
        "/contributors/role-123/position-123/confirm-delete",
      );
      testForbiddenDatasetRoute(
        "/contributors/role-123/position-123",
        "PUT",
        "DELETE",
      );
    });

    function testForbiddenDatasetRoute(
      route: string,
      ...methods: ("GET" | "PUT" | "POST" | "DELETE")[]
    ) {
      cy.then(function () {
        testForbiddenRoute(`/dataset/${this.biblioId}${route}`, ...methods);
      });
    }

    function testForbiddenPublicationRoute(
      route: string,
      ...methods: ("GET" | "PUT" | "POST" | "DELETE")[]
    ) {
      cy.then(function () {
        testForbiddenRoute(`/publication/${this.biblioId}${route}`, ...methods);
      });
    }

    function testForbiddenRoute(
      route: string,
      ...methods: ("GET" | "PUT" | "POST" | "DELETE")[]
    ) {
      if (methods.length === 0) {
        methods.push("GET");
      }

      for (const method of methods) {
        cy.request({ url: route, method, failOnStatusCode: false }).should(
          "have.property",
          "status",
          403,
        );
      }
    }
  });
});
