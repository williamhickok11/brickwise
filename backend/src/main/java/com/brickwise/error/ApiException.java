package com.brickwise.error;

/**
 * Structured API failure used across services and controllers.
 *
 * <p><b>C# analogy:</b> Similar to throwing an exception that an ASP.NET
 * {@code IExceptionFilter} / {@code ProblemDetails} middleware maps to an HTTP status.
 * In Go this was {@code APIError} with Code/Message/Status fields.
 */
public class ApiException extends RuntimeException {
    private final String code;
    private final int status;

    public ApiException(String code, String message, int status) {
        super(message);
        this.code = code;
        this.status = status;
    }

    public String getCode() {
        return code;
    }

    public int getStatus() {
        return status;
    }

    public static ApiException notFound() {
        return new ApiException("NOT_FOUND", "Resource not found", 404);
    }

    public static ApiException invalidInput(String message) {
        return new ApiException("INVALID_INPUT", message, 400);
    }

    public static ApiException validation(String message) {
        return new ApiException("VALIDATION_ERROR", message, 400);
    }

    public static ApiException internal(String message) {
        return new ApiException("INTERNAL_ERROR", message, 500);
    }
}
