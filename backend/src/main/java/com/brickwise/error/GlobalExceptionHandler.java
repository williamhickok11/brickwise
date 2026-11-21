package com.brickwise.error;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.MethodArgumentNotValidException;
import org.springframework.web.bind.annotation.ExceptionHandler;
import org.springframework.web.bind.annotation.RestControllerAdvice;

/**
 * Maps exceptions to JSON responses for every controller.
 *
 * <p><b>C# analogy:</b> Like a global exception filter or {@code UseExceptionHandler}
 * middleware that writes ProblemDetails. {@code @RestControllerAdvice} is the Spring
 * equivalent of registering one cross-cutting handler for all API controllers.
 */
@RestControllerAdvice
public class GlobalExceptionHandler {
    private static final Logger log = LoggerFactory.getLogger(GlobalExceptionHandler.class);

    @ExceptionHandler(ApiException.class)
    public ResponseEntity<ApiErrorBody> handleApi(ApiException ex) {
        if (ex.getStatus() >= 500) {
            log.error("API error: {}", ex.getMessage());
        } else {
            log.warn("API error: {} ({})", ex.getMessage(), ex.getCode());
        }
        return ResponseEntity.status(ex.getStatus()).body(new ApiErrorBody(ex.getCode(), ex.getMessage()));
    }

    @ExceptionHandler(MethodArgumentNotValidException.class)
    public ResponseEntity<ApiErrorBody> handleValidation(MethodArgumentNotValidException ex) {
        String msg = ex.getBindingResult().getFieldErrors().stream()
                .map(fe -> fe.getField() + ": " + fe.getDefaultMessage())
                .reduce((a, b) -> a + ", " + b)
                .orElse("Invalid input");
        return ResponseEntity.badRequest().body(new ApiErrorBody("VALIDATION_ERROR", msg));
    }

    @ExceptionHandler(Exception.class)
    public ResponseEntity<ApiErrorBody> handleOther(Exception ex) {
        log.error("Unexpected error", ex);
        return ResponseEntity.internalServerError()
                .body(new ApiErrorBody("INTERNAL_ERROR", "Internal server error"));
    }
}
