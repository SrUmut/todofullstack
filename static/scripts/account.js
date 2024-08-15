function resetForm(target) {
    target.reset();
}


function displayError(event) {
    console.log(event);
    document.querySelector("#error").textContent = event.detail.xhr.response;
}


function responseHeaders(event) {
    const jwt_token = event.detail.xhr.getResponseHeader("jwt-token");
    let expDur = event.detail.xhr.getResponseHeader("jwt-exp");
    expDur = parseInt(expDur);
    addJWT(jwt_token, expDur);

    //window.location.href = "/";
}