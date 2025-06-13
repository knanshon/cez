/**
 * Adds custom headers to an event based on the data attributes of the selected option in a <select> element.
 *
 * For each `data-*` attribute on the selected option, a header is added to `$event.detail.headers`
 * with the name transformed to "Hx-" followed by the attribute name in PascalCase.
 *
 * @param {Event} $event - The event object whose `detail.headers` will be modified.
 * @param {HTMLSelectElement} $el - The <select> element containing the selected option.
 */
function hxAddSelectDataHeaders($event, $el) {
    const opt = $el.selectedOptions[0];
    if (!opt) return;
    for (const attr of opt.attributes) {
        if (attr.name.startsWith("data-")) {
            const headerName = "Hx-" + attr.name.replace(/(^|-)([a-z])/g, (_, sep, c) => sep + c.toUpperCase());
            $event.detail.headers[headerName]= attr.value;
        }
    }
}