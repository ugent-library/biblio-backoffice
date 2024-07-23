// Parent commands
import login from "./login";
import switchMode from "./switch-mode";
import ensureModal from "./ensure-modal";
import ensureNoModal from "./ensure-no-modal";
import visitPublication from "./visit-publication";
import visitDataset from "./visit-dataset";
import ensureToast from "./ensure-toast";
import ensureNoToast from "./ensure-no-toast";
import setFieldByLabel from "./set-field-by-label";
import search from "./search";
import updateFields from "./update-fields";
import setUpPublication from "./set-up-publication";
import setUpDataset from "./set-up-dataset";
import clocked from "./clocked";
import {
  addAuthor,
  addCreator,
  addEditor,
  addSupervisor,
} from "./add-contributor";
import htmxRequest from "./htmx-request";

// Child commands
import finishLog from "./finish-log";
import setField from "./set-field";
import triggerHtmx from "./trigger-htmx";

// Dual commands
import extractBiblioId from "./extract-biblio-id";
import closeModal from "./close-modal";
import closeToast from "./close-toast";

// Parent commands
Cypress.Commands.addAll({
  login,

  switchMode,

  ensureModal,

  ensureNoModal,

  visitPublication,

  visitDataset,

  ensureToast,

  ensureNoToast,

  setFieldByLabel,

  search,

  updateFields,

  setUpPublication,

  setUpDataset,

  clocked,

  addAuthor,
  addEditor,
  addSupervisor,
  addCreator,

  htmxRequest,
});

// Child commands
Cypress.Commands.addAll(
  { prevSubject: true },
  {
    finishLog,

    setField,

    triggerHtmx,
  },
);

// Dual commands
Cypress.Commands.addAll(
  {
    prevSubject: "optional",
  },
  {
    extractBiblioId,

    closeModal,

    closeToast,
  },
);
