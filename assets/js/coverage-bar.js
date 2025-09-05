"use strict";

(() => {
  /**
   * @param {CustomEvent} e
   */
  function handleGofutzUpdates(e) {
    /** @type {UpdateMessage} */
    const details = e.detail;

    renderCoverage(details.params.coverage);
  }

  /**
   * @param {number} coverage
   */
  function renderCoverage(coverage) {
    /** @type {HTMLDivElement | null} */
    const coveredContainer = document.querySelector(".covered");

    if (!coveredContainer) {
      console.error("Could not find covered container");

      return;
    }

    if (coverage >= 0) {
      coveredContainer.style.width = `${coverage}%`;
    }
  }

  window.addEventListener("gofutz:init", handleGofutzUpdates);
  window.addEventListener("gofutz:update", handleGofutzUpdates);
})();
