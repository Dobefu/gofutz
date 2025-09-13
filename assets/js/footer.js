"use strict";

(() => {
  function handleGofutzUpdates() {
    renderFooterOutput();
  }

  function renderFooterOutput() {
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

    globalThis.testData.output = globalThis.testData.output
      .map((line) => {
        if (line.startsWith("{") && line.endsWith("}")) {
          try {
            const json = JSON.parse(line);

            if (!("Output" in json)) {
              return "";
            }

            return json.Output.replace(/\n$/, "");
          } catch (error) {
            console.warn("Could not parse line:", line);

            return "";
          }
        }

        return line;
      })
      .filter((line) => line !== "");

    const scrollThreshold = 10;
    const scrollPos = footer.scrollTop + footer.clientHeight;
    const shouldScroll = scrollPos >= footer.scrollHeight - scrollThreshold;

    footerOutput.innerHTML = globalThis.testData.output.join("\n");

    if (shouldScroll) {
      footer.scrollTop = footer.scrollHeight;
    }
  }

  window.addEventListener("gofutz:init", handleGofutzUpdates);
  window.addEventListener("gofutz:output", handleGofutzUpdates);
})();
