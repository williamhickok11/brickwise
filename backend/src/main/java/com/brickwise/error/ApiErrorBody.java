package com.brickwise.error;

/** JSON error body matching the Go API: {@code { "code", "message" }}. */
public record ApiErrorBody(String code, String message) {}
