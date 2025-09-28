"use strict";

(() => {
  function handleGofutzUpdates() {
    renderCoverage();
  }

  function renderCoverage() {
    /** @type {HTMLDivElement | null} */
    const coveredContainer = document.querySelector(".covered");

    if (!coveredContainer) {
      console.error("Could not find covered container");

      return;
    }

    if (globalThis.testData.coverage >= 0) {
      coveredContainer.style.width = `${globalThis.testData.coverage}%`;
    }
  }

  globalThis.addEventListener("gofutz:init", handleGofutzUpdates);
  globalThis.addEventListener("gofutz:update", handleGofutzUpdates);
})();
