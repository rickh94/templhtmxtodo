/**
 * Picks the right icon for an alert
 * @param {AlertVariant} variant - The variant
 */
function getIcon(variant) {
  switch (variant) {
    case "success":
      return `
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor" class="w-6 h-6 text-emerald-200">
          <path fill-rule="evenodd" d="M2.25 12c0-5.385 4.365-9.75 9.75-9.75s9.75 4.365 9.75 9.75-4.365 9.75-9.75 9.75S2.25 17.385 2.25 12zm13.36-1.814a.75.75 0 10-1.22-.872l-3.236 4.53L9.53 12.22a.75.75 0 00-1.06 1.06l2.25 2.25a.75.75 0 001.14-.094l3.75-5.25z" clip-rule="evenodd" />
        </svg>
        `;
    case "error":
      return `
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor" class="w-6 h-6 text-rose-200">
          <path fill-rule="evenodd" d="M2.25 12c0-5.385 4.365-9.75 9.75-9.75s9.75 4.365 9.75 9.75-4.365 9.75-9.75 9.75S2.25 17.385 2.25 12zM12 8.25a.75.75 0 01.75.75v3.75a.75.75 0 01-1.5 0V9a.75.75 0 01.75-.75zm0 8.25a.75.75 0 100-1.5.75.75 0 000 1.5z" clip-rule="evenodd" />
        </svg>
        `;
    case "warning":
      return `
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor" class="w-6 h-6 text-amber-200">
          <path fill-rule="evenodd" d="M9.401 3.003c1.155-2 4.043-2 5.197 0l7.355 12.748c1.154 2-.29 4.5-2.599 4.5H4.645c-2.309 0-3.752-2.5-2.598-4.5L9.4 3.003zM12 8.25a.75.75 0 01.75.75v3.75a.75.75 0 01-1.5 0V9a.75.75 0 01.75-.75zm0 8.25a.75.75 0 100-1.5.75.75 0 000 1.5z" clip-rule="evenodd" />
        </svg>
        `;
    default:
      return `
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor" class="w-6 h-6 text-sky-200">
          <path fill-rule="evenodd" d="M2.25 12c0-5.385 4.365-9.75 9.75-9.75s9.75 4.365 9.75 9.75-4.365 9.75-9.75 9.75S2.25 17.385 2.25 12zm8.706-1.442c1.146-.573 2.437.463 2.126 1.706l-.709 2.836.042-.02a.75.75 0 01.67 1.34l-.04.022c-1.147.573-2.438-.463-2.127-1.706l.71-2.836-.042.02a.75.75 0 11-.671-1.34l.041-.022zM12 9a.75.75 0 100-1.5.75.75 0 000 1.5z" clip-rule="evenodd" />
        </svg>
        `;
  }
}

/**
 * Displays an alert
 * @param {string} message - The message to display
 * @param {string} title - The title of the alert
 * @param {AlertVariant} variant - The variant of the alert
 * @param {number} duration - A duration to auto dismiss the alert
 */
function showAlert(message, title, variant, duration) {
  const toastId = "toast-" + Math.random().toString(36).substring(2, 15);
  const icon = getIcon(variant);

  const toastHTML = `
      <div class="p-4">
        <div class="flex items-start">
          <div class="flex-shrink-0">
            ${icon}
          </div>
          <div class="flex-1 pt-0.5 ml-3 w-0">
            <p class="text-sm font-medium">${title}</p>
            <p class="mt-1 text-sm">${message}</p>
          </div>
          <div class="flex flex-shrink-0 ml-4">
            <button type="button" class="inline-flex focus:ring-2 focus:ring-white focus:ring-offset-2 focus:outline-none"
              onclick="closeAlert('${toastId}')"
              >
              <span class="sr-only">Close</span>
              <svg class="w-5 h-5" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
                <path d="M6.28 5.22a.75.75 0 00-1.06 1.06L8.94 10l-3.72 3.72a.75.75 0 101.06 1.06L10 11.06l3.72 3.72a.75.75 0 101.06-1.06L11.06 10l3.72-3.72a.75.75 0 00-1.06-1.06L10 8.94 6.28 5.22z" />
              </svg>
            </button>
          </div>
        </div>
      </div>
  `;
  const toastDiv = document.createElement("div");
  toastDiv.className =
    "overflow-hidden w-full max-w-sm text-white bg-black border-2 border-white ring-1 ring-black ring-opacity-5 shadow-lg transition-transform pointer-events-auto shadow-white/20";
  toastDiv.id = toastId;
  toastDiv.classList.add("transform", "ease-out", "duration-300", "transition");
  toastDiv.classList.add(
    "translate-y-2",
    "opacity-0",
    "sm:translate-y-0",
    "sm:translate-x-2",
  );
  requestAnimationFrame(() => {
    requestAnimationFrame(() => {
      toastDiv.classList.remove(
        "translate-y-2",
        "opacity-0",
        "sm:translate-y-0",
        "sm:translate-x-2",
      );
      toastDiv.classList.add(
        "translate-y-0",
        "opacity-100",
        "sm:translate-x-0",
        "sm:translate-y-0",
      );
      setTimeout(() => {
        toastDiv.classList.remove(
          "translate-y-0",
          "opacity-100",
          "sm:translate-x-0",
          "sm:translate-y-0",
        );
        toastDiv.classList.remove(
          "transform",
          "ease-out",
          "duration-300",
          "transition",
        );
      }, 300);
    });
  });
  toastDiv.innerHTML = toastHTML;
  const toastContainer = document.getElementById("toast-container");
  toastContainer.append(toastDiv);
  setTimeout(() => {
    closeAlert(toastId);
  }, duration);
}

function closeAlert(id) {
  const toastDiv = document.getElementById(id);
  if (!toastDiv) return;
  toastDiv.classList.add("transform", "ease-in", "duration-300", "transition");
  toastDiv.classList.add("opacity-100");

  requestAnimationFrame(() => {
    requestAnimationFrame(() => {
      toastDiv.classList.remove("opacity-100");
      toastDiv.classList.add("opacity-0");
      setTimeout(() => {
        toastDiv.remove();
      }, 300);
    });
  });
}

/**
 * A custom event for showing alerts from the server
 * @typedef ShowAlertEvent
 * @type {object}
 * @property {ShowAlertDetail} detail
 *
 */
/**
 * @typedef {Object} ShowAlertDetail
 * @property {string} message - The message to display
 * @property {string} title - The title of the alert
 * @property {AlertVariant} variant - The variant of the alert
 * @property {number} duration - A duration to auto dismiss the alert
 */

/**
 * @typedef {"success" | "info" | "warning" | "error"} AlertVariant
 */

/**
 * Handles an alert event
 * @param {ShowAlertEvent} evt - The event
 */
function handleAlertEvent(evt) {
  const { message, title, variant, duration } = evt.detail;
  if (!message || !title || !variant || !duration) {
    throw new Error("Invalid event received from server");
  }
  showAlert(message, title, variant, duration);
}

/**
 * A custom event for showing alerts from the server
 * @typedef FocusInputEvent
 * @type {object}
 * @property {FocusInputDetail} detail
 *
 */
/**
 * @typedef {Object} FocusInputDetail
 * @property {string} id - The message to display
 */

/**
 * @typedef {"success" | "info" | "warning" | "error"} AlertVariant
 */

/**
 * Handles a focus input event
 * @param {FocusEvent} evt - The event
 */
function handleFocusInputEvent(evt) {
  const { id } = evt.detail;
  if (!id) {
    throw new Error("Invalid event received from server");
  }
  document.getElementById(id).focus();
  document.getElementById(id).select();
}

document.addEventListener("ShowAlert", handleAlertEvent);
document.addEventListener("FocusInput", handleFocusInputEvent);
