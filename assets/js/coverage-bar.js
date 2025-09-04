"use strict";

(() => {
  /**
   * @param {CustomEvent} e
   */
  function handleGofutzInit(e) {
    /** @type {Message} */
    const details = e.detail;

    renderCoverage(details.params.coverage);
  }

  /**
   * @param {CustomEvent} e
   */
  function handleGofutzUpdate(e) {
    /** @type {Message} */
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

  window.addEventListener("gofutz:init", handleGofutzInit);
  window.addEventListener("gofutz:update", handleGofutzUpdate);
})();
