function resetForm(target) {
    target.reset();
}

function displayError(event) {
    const errorDiv = document.querySelector("#error");
    errorDiv.textContent = event.detail.xhr.response;
    errorDiv.classList.remove("d-none");
}
