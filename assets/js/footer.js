"use strict";

(() => {
  /**
   * @param {CustomEvent} e
   */
  function handleGofutzInit(e) {
    /** @type {InitMessage} */
    const details = e.detail;

    renderFooterOutput(details.params.output);
  }

  /**
   * @param {CustomEvent} e
   */
  function handleGofutzOutput(e) {
    /** @type {OutputMessage} */
    const details = e.detail;

    renderFooterOutput(details.params.output);
  }

  /**
   * @param {string[]} output
   */
  function renderFooterOutput(output) {
    /** @type {HTMLPreElement | null} */
    const footerOutput = document.querySelector(".footer__output");

    if (!footerOutput) {
      console.error("Could not find footer output");

      return;
    }

    /** @type {HTMLElement | null} */
    const footer = footerOutput.parentElement;

    if (!footer) {
      console.error("Could not find footer");

      return;
    }

    const scrollThreshold = 10;
    const scrollPos = footer.scrollTop + footer.clientHeight;
    const shouldScroll = scrollPos >= footer.scrollHeight - scrollThreshold;

    footerOutput.innerHTML = output.join("\n");

    if (shouldScroll) {
      footer.scrollTop = footer.scrollHeight;
    }
  }

  window.addEventListener("gofutz:init", handleGofutzInit);
  window.addEventListener("gofutz:output", handleGofutzOutput);
})();
