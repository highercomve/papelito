if (!Element.prototype.matches) {
  Element.prototype.matches =
    Element.prototype.msMatchesSelector ||
    Element.prototype.webkitMatchesSelector;
}

if (!Element.prototype.closest) {
  Element.prototype.closest = function (s) {
    var el = this;

    do {
      if (Element.prototype.matches.call(el, s)) return el;
      el = el.parentElement || el.parentNode;
    } while (el !== null && el.nodeType === 1);
    return null;
  };
}

function navbarToggle(event) {
  let childSelector = event.target.tagName;
  if (Array.from(event.target.classList).length > 0) {
    childSelector += "." + Array.from(event.target.classList).join(".");
  }
  const selector = ".navbar-toggler " + childSelector;
  if (
    !event.target.classList.contains("navbar-toggler") &&
    !document.querySelector(selector)
  ) {
    return;
  }

  if (!event.target.closest) {
    return;
  }
  const navbar = event.target
    .closest(".navbar")
    .querySelector(".navbar-collapse");
  navbar.classList.toggle("show");
}

function preffier() {
  const elements = document.getElementsByClassName("prettyprint-json");

  Array.from(elements).forEach((elem) => {
    if (elem.classList.contains("preffier-done")) { return }
    elem.classList.add("preffier-done")
    const jsonString = JSON.stringify(JSON.parse(elem.innerHTML), null, 2)
    elem.innerHTML = jsonString;
  });
}

function confirmableSubmit(event) {
  const formData = new FormData(event.target);
  const confirmable = formData.get("confirmable");
  const confirmation = formData.get("confirmation");
  if (
    confirmable === "" ||
    confirmation === "" ||
    confirmation !== confirmable
  ) {
    window.alert("The confirmation is not correct");
    event.preventDefault();
    return;
  }

  const message = event.target.dataset.confirmation || "Are you sure?";
  if (!window.confirm(message)) {
    event.preventDefault();
    return;
  }
}

function enableConfirmableSubmit() {
  Array.from(document.querySelectorAll("form.confirmable")).forEach((e) => {
    e.addEventListener("submit", confirmableSubmit);
  });
}

function collapseToggle(event) {
  if (event.target.dataset["toggle"] === "collapse") {
    event.preventDefault();
    const target = document.querySelector(event.target.dataset["target"]);
    target.classList.toggle("show");
  }
}

function modalToggle(event) {
  if (event.target.dataset["toggle"] === "modal") {
    event.preventDefault();
    const target = document.querySelector(event.target.dataset["target"]);
    target.classList.toggle("show");
    document.body.classList.toggle("modal-open");
  }
}

function dismissModal(event) {
  if (
    event.target &&
    event.target.closest &&
    event.target.dataset["dismiss"] === "modal"
  ) {
    event.preventDefault();
    const target = event.target.closest(".modal");
    target.classList.remove("show");
    document.body.classList.remove("modal-open");
  }

  if (
    event.target.dataset["toggle"] !== "modal" &&
    !event.target.closest(".modal-dialog") &&
    document.body.classList.contains("modal-open")
  ) {
    event.preventDefault();
    const target = document.querySelector(".modal");
    target.classList.remove("show");
    document.body.classList.remove("modal-open");
  }
}

function dropDownToggle(event) {
  if (!event.target || !event.target.closest) {
    return;
  }

  const dropdown = event.target.closest(".dropdown-toggle");
  if (!dropdown) {
    const targetElements = document.querySelectorAll(".dropdown-menu");
    targetElements.forEach((e) => e.classList.remove("show"));
    return;
  }


  const target = dropdown.dataset["toggleTarget"];
  if (!target || target === "") {
    return;
  }

  const targetElement = document.getElementById(target);
  if (!targetElement) {
    return;
  }

  if (targetElement === event.target && targetElement.contains(event.target)) {
    return
  }

  event.preventDefault();
  targetElement.classList.toggle("show");
}

export function Load() {
  document.addEventListener("click", navbarToggle);
  document.addEventListener("click", collapseToggle);
  document.addEventListener("click", modalToggle);
  document.addEventListener("click", dismissModal);
  document.addEventListener("click", dropDownToggle);

  preffier();
  enableConfirmableSubmit();
}
