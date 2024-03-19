import getLabel from "./get-label";
import getParams from "./get-params";

Cypress.Commands.addQuery("getLabel", getLabel);
Cypress.Commands.addQuery("getParams", getParams);
