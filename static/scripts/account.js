function resetForm(target) {
    target.reset();
}


function displayError(event) {
    console.log(event);
    document.querySelector("#error").textContent = event.detail.xhr.response;
}
