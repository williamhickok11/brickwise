package com.brickwise.web;

import com.brickwise.model.Document;
import com.brickwise.service.DocumentService;
import org.springframework.core.io.Resource;
import org.springframework.http.HttpHeaders;
import org.springframework.http.HttpStatus;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;
import org.springframework.web.multipart.MultipartFile;

@RestController
@RequestMapping("/api/v1/residences/{residenceId}/document")
public class DocumentController {
    private final DocumentService service;

    public DocumentController(DocumentService service) {
        this.service = service;
    }

    @GetMapping
    public Document get(@PathVariable int residenceId) {
        return service.getByResidenceId(residenceId);
    }

    @PostMapping
    @ResponseStatus(HttpStatus.CREATED)
    public Document upload(@PathVariable int residenceId, @RequestParam("file") MultipartFile file) {
        return service.upload(residenceId, file);
    }

    @GetMapping("/download")
    public ResponseEntity<Resource> download(@PathVariable int residenceId) {
        Document meta = service.getByResidenceId(residenceId);
        Resource body = service.download(residenceId);
        return ResponseEntity.ok()
                .header(HttpHeaders.CONTENT_DISPOSITION, "attachment; filename=\"" + meta.getFileName() + "\"")
                .contentType(MediaType.APPLICATION_PDF)
                .body(body);
    }
}
