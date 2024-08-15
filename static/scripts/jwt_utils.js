const SECOND = 1000;
const MINUTE = 60 * SECOND;
const HOUR = 60 * MINUTE;
const DAY = 24 * HOUR;
const jwt_header_name = "jwt_token";

function addCookie(name, value, expiresIn) {
    try {
        checkCookieProperties(name, value, expiresIn);
    } catch (error) {
        throw error;
    }
    name = name.trim();
    value = value.trim();
    document.cookie = name + "=" + value + "; expires=" + new Date(Date.now() + expiresIn).toUTCString();
}

function checkCookieProperties(name, value, expiresIn) {
    if (typeof name !== "string") {
        throw new TypeError("name must be a string");
    } else if (typeof value !== "string") {
        throw new TypeError("value must be a string");
    } else if (typeof expiresIn !== "number") {
        throw new TypeError("expiresIn must be a number");
    }

    if (name.length === 0) {
        throw new RangeError("name must not be empty");
    } else if (value.length === 0) {
        throw new RangeError("value must not be empty");
    } else if (expiresIn < 0) {
        throw new RangeError("expiresIn must be a positive number or zero in milliseconds");
    }
}

function getCookie(name) {
    for (var cookie of document.cookie.split(";")) {
        cookie = cookie.trim();
        if (cookie.startsWith(name + "=")) {
            return cookie.substring(name.length + 1);
        }
    }
    return null;
}

function addJWT(jwtTokenn, expiresIn) {
    addCookie(jwt_header_name, jwtTokenn, expiresIn);
}

function getJWT() {
    return getCookie(jwt_header_name);
}